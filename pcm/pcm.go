/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2015 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pcm

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/intelsdi-x/snap-plugin-utilities/ns"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"

	"github.com/intelsdi-x/snap-plugin-collector-pcm/parser"
)

const (
	// Name of plugin
	name = "pcm"
	// Version of plugin
	version = 11
	// Type of plugin
	pluginType = plugin.CollectorPluginType
)

func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(name, version, pluginType, []string{plugin.SnapGOBContentType}, []string{plugin.SnapGOBContentType},
		plugin.Exclusive(true)) // this should prevent PCM racing for MSRs from multiple instances
}

// PCM
type PCM struct {
	parser           parser.Parser
	initializedMutex *sync.Mutex
	initialized      bool

	RawToNs, NsToRaw map[parser.Key]parser.Key
}

func (p *PCM) initialize() error {
	p.initializedMutex.Lock()
	defer p.initializedMutex.Unlock()

	if !p.initialized {
		err := p.run()
		if err != nil {
			log.WithFields(log.Fields{
				"block":    "initialize",
				"function": "run",
			}).Error(err)
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		keys, err := p.parser.GetKeys(ctx)
		if err != nil {
			log.WithFields(log.Fields{
				"block":    "initialize",
				"function": "GetKeys",
			}).Error(err)
			return err
		}
		p.RawToNs = map[parser.Key]parser.Key{}
		p.NsToRaw = map[parser.Key]parser.Key{}
		for _, key := range keys {
			nsComp := ns.ReplaceNotAllowedCharsInNamespacePart(key.Component)
			nsMt := ns.ReplaceNotAllowedCharsInNamespacePart(key.MetricName)
			nsKey := parser.Key{Component: nsComp, MetricName: nsMt}
			p.RawToNs[key] = nsKey
			p.NsToRaw[nsKey] = key

		}

		p.initialized = true
	}

	return nil

}

// // CollectMetrics returns metrics from pcm
func (p *PCM) CollectMetrics(mts []plugin.MetricType) ([]plugin.MetricType, error) {
	err := p.initialize()
	if err != nil {
		return nil, err
	}

	vals, err := p.parser.GetValues(context.TODO())
	if err != nil {
		log.WithFields(log.Fields{
			"block":    "CollectMetrics",
			"function": "GetValues",
		}).Error(err)
		return nil, err
	}

	for i := range mts {
		ns := mts[i].Namespace()
		key := p.NsToRaw[parser.Key{Component: ns[2].Value, MetricName: ns[3].Value}]
		if v, ok := vals[key]; ok {
			mts[i].Data_ = v
			mts[i].Timestamp_ = time.Now()
		}
	}
	// fmt.Fprintf(os.Stderr, "collected >>> %+v\n", metrics)
	return mts, nil
}

// GetMetricTypes returns the metric types exposed by pcm
func (p *PCM) GetMetricTypes(_ plugin.ConfigType) ([]plugin.MetricType, error) {
	err := p.initialize()
	if err != nil {
		return nil, err
	}
	mts := []plugin.MetricType{}
	for nsKey := range p.NsToRaw {
		mt := plugin.MetricType{Namespace_: core.NewNamespace("intel", "pcm", nsKey.Component, nsKey.MetricName)}
		mts = append(mts, mt)
	}
	return mts, nil
}

//GetConfigPolicy
func (p *PCM) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	return c, nil
}

func NewPCMCollector() *PCM {
	return &PCM{initializedMutex: &sync.Mutex{}}
}

var dirtyMock io.Reader

func (pcm *PCM) run() error {
	if dirtyMock != nil {
		pcm.parser = parser.RunParser(dirtyMock)
		return nil
	}
	var cmd *exec.Cmd
	if path := os.Getenv("SNAP_PCM_PATH"); path != "" {
		cmd = exec.Command(filepath.Join(path, "pcm.x"), "/csv", "-nc", "-r", "1")
	} else {
		c, err := exec.LookPath("pcm.x")
		if err != nil {
			fmt.Fprint(os.Stderr, "Unable to find PCM.  Ensure it's in your path or set SNAP_PCM_PATH.")
			return err
		}
		cmd = exec.Command(c, "/csv", "-nc", "-r", "1")
	}

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Error creating StdoutPipe %v", err)
	}

	pcm.parser = parser.RunParser(cmdReader)

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("Error starting pcm %v", err)
	}

	// LEAVE the following block for debugging
	// err = cmd.Wait()
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, "Error waiting for pcm", err)
	// 	return nil, err
	// }

	return nil
}

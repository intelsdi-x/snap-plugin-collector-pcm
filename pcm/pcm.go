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
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/intelsdi-x/snap-plugin-utilities/ns"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	log "github.com/sirupsen/logrus"
)

const (
	// Name of plugin
	name = "pcm"
	// Version of plugin
	version = 11
	// Type of plugin
	pluginType = plugin.CollectorPluginType
)

var fieldsToSkip = 2

func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(name, version, pluginType, []string{plugin.SnapGOBContentType}, []string{plugin.SnapGOBContentType})
}

// PCM
type PCM struct {
	keys        []string
	data        map[string]float64
	mutex       *sync.RWMutex
	initialized bool
}

func (p *PCM) Keys() []string {
	return p.keys
}

func (p *PCM) Data() map[string]float64 {
	return p.data
}

// // CollectMetrics returns metrics from pcm
func (p *PCM) CollectMetrics(mts []plugin.MetricType) ([]plugin.MetricType, error) {
	if p.initialized == false {
		err := p.run()
		if err != nil {
			log.WithFields(log.Fields{
				"block":    "CollectMetrics",
				"function": "run",
			}).Error(err)
			return nil, err
		}
		p.initialized = true
	}
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	for i := range mts {
		if v, ok := p.data[mts[i].Namespace().String()]; ok {
			mts[i].Data_ = v
			mts[i].Timestamp_ = time.Now()
		}
	}
	// fmt.Fprintf(os.Stderr, "collected >>> %+v\n", metrics)
	return mts, nil
}

// GetMetricTypes returns the metric types exposed by pcm
func (p *PCM) GetMetricTypes(_ plugin.ConfigType) ([]plugin.MetricType, error) {
	mts := []plugin.MetricType{}
	if p.initialized == false {
		err := p.run()
		if err != nil {
			log.WithFields(log.Fields{
				"block":    "GetMetricTypes",
				"function": "run",
			}).Error(err)
			return nil, err
		}
		p.initialized = true
	}
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	for _, k := range p.keys {
		mt := plugin.MetricType{Namespace_: core.NewNamespace(strings.Split(strings.TrimPrefix(k, "/"), "/")...)}
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
	return &PCM{mutex: &sync.RWMutex{}, data: map[string]float64{}, initialized: false}
}
func (pcm *PCM) run() error {
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

	go func() {
		pcm.parse(cmdReader)
	}()

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("Error starting pcm %v", err)
	}

	// we need to wait until we have our metric types
	st := time.Now()
	for {
		pcm.mutex.RLock()
		c := len(pcm.keys)
		pcm.mutex.RUnlock()
		if c > 0 {
			break
		}
		if time.Since(st) > time.Second*2 {
			return fmt.Errorf("Timed out waiting for metrics from pcm")
		}
	}

	// LEAVE the following block for debugging
	// err = cmd.Wait()
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, "Error waiting for pcm", err)
	// 	return nil, err
	// }

	return nil
}

func (pcm *PCM) parse(reader io.Reader) {
	// read the data from stdout
	scanner := bufio.NewScanner(reader)
	first := true
	header := []string{}
	for scanner.Scan() {
		if first {
			first = false
			currentKey := ""
			keys := strings.Split(strings.TrimSuffix(scanner.Text(), ";"), ";")
			for _, key := range keys {
				if key != "" {
					currentKey = key
				}
				header = append(header, currentKey)
			}
			continue
		}
		if len(pcm.keys) == 0 {
			pcm.mutex.Lock()
			keys := strings.Split(strings.TrimSuffix(scanner.Text(), ";"), ";")
			//skip the date and time fields
			pcm.keys = make([]string, len(keys[fieldsToSkip:]))
			for i, k := range keys[fieldsToSkip:] {
				// removes all spaces from metric key
				metricKey := ns.ReplaceNotAllowedCharsInNamespacePart(k)
				metricComponent := ns.ReplaceNotAllowedCharsInNamespacePart(header[i+fieldsToSkip])
				pcm.keys[i] = fmt.Sprintf("/intel/pcm/%s/%s", metricComponent, metricKey)
			}
			pcm.mutex.Unlock()
			continue
		}

		pcm.mutex.Lock()
		datal := strings.Split(strings.TrimSuffix(scanner.Text(), ";"), ";")
		for i, d := range datal[fieldsToSkip:] {
			v, err := strconv.ParseFloat(strings.TrimSpace(d), 64)
			if err == nil {
				pcm.data[pcm.keys[i]] = v
			} else {
				fmt.Fprintln(os.Stderr, "Invalid metric value", err)
				pcm.data[pcm.keys[i]] = math.NaN()
			}
		}
		pcm.mutex.Unlock()
		// fmt.Fprintf(os.Stderr, "data >>> %+v\n", pcm.data)
		// fmt.Fprintf(os.Stdout, "data >>> %+v\n", pcm.data)
	}
}

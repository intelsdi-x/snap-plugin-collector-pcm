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
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"io"

	"math"

	"github.com/intelsdi-x/snap-plugin-utilities/ns"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
)

const (
	// Name of plugin
	name = "pcm"
	// Version of plugin
	version = 9
	// Type of plugin
	pluginType = plugin.CollectorPluginType
)

func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(name, version, pluginType, []string{plugin.SnapGOBContentType}, []string{plugin.SnapGOBContentType})
}

// PCM
type PCM struct {
	keys  []string
	data  map[string]float64
	mutex *sync.RWMutex
}

func (p *PCM) Keys() []string {
	return p.keys
}

func (p *PCM) Data() map[string]float64 {
	return p.data
}

// // CollectMetrics returns metrics from pcm
func (p *PCM) CollectMetrics(mts []plugin.MetricType) ([]plugin.MetricType, error) {
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
	mts := make([]plugin.MetricType, len(p.keys))
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	for i, k := range p.keys {
		mts[i] = plugin.MetricType{Namespace_: core.NewNamespace(strings.Split(strings.TrimPrefix(k, "/"), "/")...)}
	}
	return mts, nil
}

//GetConfigPolicy
func (p *PCM) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	return c, nil
}

func NewPCMCollector() (*PCM, error) {
	pcm := &PCM{mutex: &sync.RWMutex{}, data: map[string]float64{}}

	err := pcm.run()
	if err != nil {
		return nil, err
	}

	return pcm, nil
}
func (pcm *PCM) run() error {
	var cmd *exec.Cmd
	if path := os.Getenv("SNAP_PCM_PATH"); path != "" {
		cmd = exec.Command(filepath.Join(path, "pcm.x"), "/csv", "-nc", "-r", "1")
	} else {
		c, err := exec.LookPath("pcm.x")
		if err != nil {
			fmt.Fprint(os.Stderr, "Unable to find PCM.  Ensure it's in your path or set SNAP_PCM_PATH.")
			panic(err)
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
	for scanner.Scan() {
		if first {
			first = false
			continue
		}
		if len(pcm.keys) == 0 {
			pcm.mutex.Lock()
			keys := strings.Split(strings.TrimSuffix(scanner.Text(), ";"), ";")
			//skip the date and time fields
			pcm.keys = make([]string, len(keys[2:]))
			for i, k := range keys[2:] {
				// removes all spaces from metric key
				metricKey := ns.ReplaceNotAllowedCharsInNamespacePart(k)
				pcm.keys[i] = fmt.Sprintf("/intel/pcm/%s", metricKey)
			}
			pcm.mutex.Unlock()
			continue
		}

		pcm.mutex.Lock()
		datal := strings.Split(strings.TrimSuffix(scanner.Text(), ";"), ";")
		for i, d := range datal[2:] {
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

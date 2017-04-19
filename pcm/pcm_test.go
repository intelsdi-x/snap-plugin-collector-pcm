//
// +build small

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
	"strings"
	"sync"
	"testing"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	. "github.com/smartystreets/goconvey/convey"
)

var mockCmdOut = `System;;;;;;;;;;;;;;;;;;;;;System Core C-States;;;;;System Pack C-States;;;;;;;;Socket0;;;;;;;;;;;;;SKT0 Core C-State;;;;;SKT0 Package C-State;;;;;;;Proc Energy (Joules);
Date;Time;EXEC;IPC;FREQ;AFREQ;L3MISS;L2MISS;L3HIT;L2HIT;L3MPI;L2MPI;READ;WRITE;INST;ACYC;TIME(ticks);PhysIPC;PhysIPC%;INSTnom;INSTnom%;C0res%;C1res%;C3res%;C6res%;C7res%;C2res%;C3res%;C6res%;C7res%;C8res%;C9res%;C10res%;Proc Energy (Joules);EXEC;IPC;FREQ;AFREQ;L3MISS;L2MISS;L3HIT;L2HIT;L3MPI;L2MPI;READ;WRITE;TEMP;C0res%;C1res%;C3res%;C6res%;C7res%;C2res%;C3res%;C6res%;C7res%;C8res%;C9res%;C10res%;SKT0;
2017-04-14;13:52:31.359;0.00245;0.589;0.00415;0.262;0.249;0.666;0.613;0.439;0.00372;0.00997;0.761;0.134;66.8;113;3.41e+03;1.18;29.5;0.0049;0.122;1.59;12.6;0.0181;1.64;84.1;50.1;0;0;0;0;0;0;2.37;0.00245;0.589;0.00415;0.262;0.249;0.666;0.613;0.439;0.00372;0.00997;0.761;0.134;69;1.59;12.6;0.0181;1.64;84.1;50.1;0;0;0;0;0;0;2.37;
`

var refMap = map[string]float64{
	"/intel/pcm/System/L2MPI":                            0.00997,
	"/intel/pcm/SKT0_Core_C-State/C7res%":                84.1,
	"/intel/pcm/SKT0_Package_C-State/C3res%":             0,
	"/intel/pcm/System/FREQ":                             0.00415,
	"/intel/pcm/System_Pack_C-States/Proc_Energy_Joules": 2.37,
	"/intel/pcm/System/IPC":                              0.589,
	"/intel/pcm/System/INSTnom%":                         0.122,
	"/intel/pcm/Socket0/FREQ":                            0.00415,
	"/intel/pcm/Socket0/L3MISS":                          0.249,
	"/intel/pcm/System/AFREQ":                            0.262,
	"/intel/pcm/Socket0/L2HIT":                           0.439,
	"/intel/pcm/System/ACYC":                             113,
	"/intel/pcm/Socket0/IPC":                             0.589,
	"/intel/pcm/Proc_Energy_Joules/SKT0":                 2.37,
	"/intel/pcm/System_Core_C-States/C3res%":             0.0181,
	"/intel/pcm/SKT0_Core_C-State/C1res%":                12.6,
	"/intel/pcm/System/INSTnom":                          0.0049,
	"/intel/pcm/System_Core_C-States/C6res%":             1.64,
	"/intel/pcm/System_Pack_C-States/C2res%":             50.1,
	"/intel/pcm/Socket0/L3MPI":                           0.00372,
	"/intel/pcm/SKT0_Core_C-State/C6res%":                1.64,
	"/intel/pcm/System_Pack_C-States/C3res%":             0,
	"/intel/pcm/SKT0_Core_C-State/C3res%":                0.0181,
	"/intel/pcm/System_Pack_C-States/C6res%":             0,
	"/intel/pcm/System_Pack_C-States/C8res%":             0,
	"/intel/pcm/Socket0/WRITE":                           0.134,
	"/intel/pcm/SKT0_Package_C-State/C7res%":             0,
	"/intel/pcm/SKT0_Package_C-State/C10res%":            0,
	"/intel/pcm/System/READ":                             0.761,
	"/intel/pcm/System/PhysIPC":                          1.18,
	"/intel/pcm/Socket0/AFREQ":                           0.262,
	"/intel/pcm/Socket0/L3HIT":                           0.613,
	"/intel/pcm/Socket0/READ":                            0.761,
	"/intel/pcm/System/L3HIT":                            0.613,
	"/intel/pcm/System_Pack_C-States/C7res%":             0,
	"/intel/pcm/Socket0/L2MPI":                           0.00997,
	"/intel/pcm/Socket0/TEMP":                            69,
	"/intel/pcm/System/EXEC":                             0.00245,
	"/intel/pcm/System/L3MISS":                           0.249,
	"/intel/pcm/Socket0/L2MISS":                          0.666,
	"/intel/pcm/SKT0_Package_C-State/C8res%":             0,
	"/intel/pcm/SKT0_Package_C-State/C9res%":             0,
	"/intel/pcm/System/WRITE":                            0.134,
	"/intel/pcm/System/TIME_ticks":                       3410,
	"/intel/pcm/System_Core_C-States/C1res%":             12.6,
	"/intel/pcm/System_Core_C-States/C7res%":             84.1,
	"/intel/pcm/Socket0/EXEC":                            0.00245,
	"/intel/pcm/System/L3MPI":                            0.00372,
	"/intel/pcm/System/PhysIPC%":                         29.5,
	"/intel/pcm/System_Core_C-States/C0res%":             1.59,
	"/intel/pcm/SKT0_Core_C-State/C0res%":                1.59,
	"/intel/pcm/System/L2HIT":                            0.439,
	"/intel/pcm/System/INST":                             66.8,
	"/intel/pcm/System_Pack_C-States/C9res%":             0,
	"/intel/pcm/System_Pack_C-States/C10res%":            0,
	"/intel/pcm/SKT0_Package_C-State/C2res%":             50.1,
	"/intel/pcm/System/L2MISS":                           0.666,
	"/intel/pcm/SKT0_Package_C-State/C6res%":             0,
}

func TestPCMPlugin(t *testing.T) {
	Convey("Meta should return metadata for the plugin", t, func() {
		meta := Meta()
		So(meta.Name, ShouldResemble, name)
		So(meta.Version, ShouldResemble, version)
		So(meta.Type, ShouldResemble, plugin.CollectorPluginType)
	})

	Convey("Create PCM Collector", t, func() {
		pcm := &PCM{mutex: &sync.RWMutex{}, data: map[string]float64{}}
		reader := strings.NewReader(mockCmdOut)
		pcm.parse(reader)

		configPolicy, err := pcm.GetConfigPolicy()
		Convey("pcmCol.GetConfigPolicy() should return a config policy", func() {
			Convey("So config policy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)
			})
			Convey("So getting config policy should not return an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("So config policy should be a cpolicy.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, &cpolicy.ConfigPolicy{})
			})
		})

		Convey("Given valid static metric namespace collect metrics", func() {

			var mockMts []plugin.MetricType

			for key := range refMap {
				mockMts = append(mockMts, plugin.MetricType{Namespace_: core.NewNamespace(strings.Split(strings.TrimPrefix(key, "/"), "/")...)})
			}

			So(func() { pcm.CollectMetrics(mockMts) }, ShouldNotPanic)
			result, err := pcm.CollectMetrics(mockMts)
			So(err, ShouldBeNil)
			So(len(result), ShouldEqual, len(refMap))

			m := make(map[string]float64, len(result))

			for _, r := range result {
				m[r.Namespace().String()] = r.Data().(float64)
			}

			So(m, ShouldResemble, refMap)
		})

		Convey("Get metric types", func() {
			mts, err := pcm.GetMetricTypes(plugin.ConfigType{})
			So(err, ShouldBeNil)
			So(len(mts), ShouldEqual, len(refMap))

			namespaces := []string{}
			for _, m := range mts {
				namespaces = append(namespaces, m.Namespace().String())
			}

			for k := range refMap {
				So(namespaces, ShouldContain, k)
			}
		})
	})
}

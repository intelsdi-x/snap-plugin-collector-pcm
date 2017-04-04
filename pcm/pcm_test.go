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
	"testing"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	. "github.com/smartystreets/goconvey/convey"
)

var mockCmdOut = `System;;;;;;;;;;;;;;;;;;;;;System Core C-States;;;;System Pack C-States;;;Socket0;;;;;;;;;;;;;SKT0 Core C-State;;;;SKT0 Package C-State;;;
Date;Time;EXEC;IPC;FREQ;AFREQ;L3MISS;L2MISS;L3HIT;L2HIT;L3MPI;L2MPI;READ;WRITE;INST;ACYC;TIME(ticks);PhysIPC;PhysIPC%;INSTnom;INSTnom%;C0res%;C1res%;C3res%;C6res%;C3res%;C6res%;C7res%;EXEC;IPC;FREQ;AFREQ;L3MISS;L2MISS;L3HIT;L2HIT;L3MPI;L2MPI;READ;WRITE;TEMP;C0res%;C1res%;C3res%;C6res%;C3res%;C6res%;C7res%;
2016-11-02;14:10:27.600;0.0175;1.13;0.0155;0.84;0.0925;0.617;0.85;0.335;0.000237;0.00158;0.0361;0.0157;391;347;2.8e+03;2.26;56.4;0.0349;0.873;1.85;2.13;1.14;94.9;0.382;76.6;0;0.0175;1.13;0.0155;0.84;0.0925;0.617;0.85;0.335;0.000237;0.00158;0.0361;0.0157;N/A;1.85;2.13;1.14;94.9;0.382;76.6;0;
`

var refMap = map[string]float64{
	"/intel/pcm/Socket0/L3MISS":              0.0925,
	"/intel/pcm/Socket0/L2MISS":              0.617,
	"/intel/pcm/Socket0/READ":                0.0361,
	"/intel/pcm/SKT0_Core_C-State/C1res%":    2.13,
	"/intel/pcm/System/L2MISS":               0.617,
	"/intel/pcm/System/ACYC":                 347,
	"/intel/pcm/System/PhysIPC":              2.26,
	"/intel/pcm/System_Core_C-States/C1res%": 2.13,
	"/intel/pcm/System_Pack_C-States/C6res%": 76.6,
	"/intel/pcm/Socket0/IPC":                 1.13,
	"/intel/pcm/SKT0_Core_C-State/C0res%":    1.85,
	"/intel/pcm/System/IPC":                  1.13,
	"/intel/pcm/System/INSTnom":              0.0349,
	"/intel/pcm/Socket0/FREQ":                0.0155,
	"/intel/pcm/Socket0/L3MPI":               0.000237,
	"/intel/pcm/SKT0_Core_C-State/C6res%":    94.9,
	"/intel/pcm/System/TIME_ticks":           2800,
	"/intel/pcm/System_Pack_C-States/C3res%": 0.382,
	"/intel/pcm/System_Pack_C-States/C7res%": 0,
	"/intel/pcm/Socket0/L2MPI":               0.00158,
	"/intel/pcm/SKT0_Package_C-State/C3res%": 0.382,
	"/intel/pcm/System/L3HIT":                0.85,
	"/intel/pcm/System/WRITE":                0.0157,
	"/intel/pcm/System/PhysIPC%":             56.4,
	"/intel/pcm/System/INSTnom%":             0.873,
	"/intel/pcm/System/L3MISS":               0.0925,
	"/intel/pcm/System/FREQ":                 0.0155,
	"/intel/pcm/System/AFREQ":                0.84,
	"/intel/pcm/System/L3MPI":                0.000237,
	"/intel/pcm/System/L2MPI":                0.00158,
	"/intel/pcm/Socket0/EXEC":                0.0175,
	"/intel/pcm/System/EXEC":                 0.0175,
	"/intel/pcm/System/INST":                 391,
	"/intel/pcm/System_Core_C-States/C6res%": 94.9,
	"/intel/pcm/Socket0/AFREQ":               0.84,
	"/intel/pcm/Socket0/L2HIT":               0.335,
	"/intel/pcm/SKT0_Core_C-State/C3res%":    1.14,
	"/intel/pcm/SKT0_Package_C-State/C6res%": 76.6,
	"/intel/pcm/System/READ":                 0.0361,
	"/intel/pcm/System_Core_C-States/C0res%": 1.85,
	"/intel/pcm/System_Core_C-States/C3res%": 1.14,
	"/intel/pcm/Socket0/L3HIT":               0.85,
	"/intel/pcm/Socket0/WRITE":               0.0157,
	"/intel/pcm/SKT0_Package_C-State/C7res%": 0,
	"/intel/pcm/System/L2HIT":                0.335,
}

func TestPCMPlugin(t *testing.T) {
	Convey("Meta should return metadata for the plugin", t, func() {
		meta := Meta()
		So(meta.Name, ShouldResemble, name)
		So(meta.Version, ShouldResemble, version)
		So(meta.Type, ShouldResemble, plugin.CollectorPluginType)
	})

	Convey("Create PCM Collector", t, func() {
		dirtyMock = strings.NewReader(mockCmdOut)
		pcm := NewPCMCollector()

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
			So(len(result), ShouldEqual, 45)
			So(err, ShouldBeNil)

			m := make(map[string]float64, len(result))

			for _, r := range result {
				m[r.Namespace().String()] = r.Data().(float64)
			}

			So(m, ShouldResemble, refMap)
		})

		Convey("Get metric types", func() {
			mts, err := pcm.GetMetricTypes(plugin.ConfigType{})
			So(err, ShouldBeNil)
			So(len(mts), ShouldEqual, 46)

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

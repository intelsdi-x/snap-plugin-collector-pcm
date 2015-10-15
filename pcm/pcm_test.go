//
// +build unit

package pcm

import (
	"testing"

	"github.com/intelsdi-x/pulse/control/plugin"
	"github.com/intelsdi-x/pulse/control/plugin/cpolicy"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPCMPlugin(t *testing.T) {
	Convey("Meta should return metadata for the plugin", t, func() {
		meta := Meta()
		So(meta.Name, ShouldResemble, name)
		So(meta.Version, ShouldResemble, version)
		So(meta.Type, ShouldResemble, plugin.CollectorPluginType)
	})

	Convey("Create PCM Collector", t, func() {
		pcmCol, err := NewPCMCollector()
		Convey("So pcmCol should not be nil", func() {
			So(pcmCol, ShouldNotBeNil)
		})
		Convey("So err should be nil", func() {
			So(err, ShouldBeNil)
		})
		Convey("So pcmCol should be of Psutil type", func() {
			So(pcmCol, ShouldHaveSameTypeAs, &PCM{})
		})
		configPolicy, err := pcmCol.GetConfigPolicy()
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
	})
}

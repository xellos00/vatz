package prometheus

import (
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	"reflect"
	"testing"
)

func Test_prometheusManager_getPluginUp(t *testing.T) {
	type fields struct {
		Protocol string
	}
	type args struct {
		plugins  []config.Plugin
		hostName string
	}
	var tests []struct {
		name         string
		fields       fields
		args         args
		wantPluginUp map[int]*prometheusValue
	}
	var grpcClients []pluginpb.PluginClient
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &prometheusManager{
				Protocol: tt.fields.Protocol,
			}
			if gotPluginUp := c.getPluginUp(tt.args.plugins, tt.args.hostName, grpcClients); !reflect.DeepEqual(gotPluginUp, tt.wantPluginUp) {
				t.Errorf("getPluginUp() = %v, want %v", gotPluginUp, tt.wantPluginUp)
			}
		})
	}
}

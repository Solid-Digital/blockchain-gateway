package nodeauth

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		rawurl string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "http",
			args: args{
				rawurl: "https://mainnet.nodes.unchain.io/v0/30",
			},
			want: map[string]string{
				"networkUUID": "30",
				"protocol":    "http",
			},
			wantErr: false,
		},
		{
			name: "ws",
			args: args{
				rawurl: "https://mainnet.nodes.unchain.io/v0/ws/30",
			},
			want: map[string]string{
				"networkUUID": "30",
				"protocol":    "websocket",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.rawurl)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

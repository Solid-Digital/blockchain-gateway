package mario

import (
	"testing"

	"bitbucket.org/unchain/ares/pkg/testhelper"

	"github.com/unchainio/pkg/xlogger"

	"bitbucket.org/unchain/ares/pkg/ares"
)

func Test_marioClient_Build(t *testing.T) {
	if testhelper.InBitBucket() {
		t.Skip("Skipping in bitbucket")
	}

	type fields struct {
		config *Config
	}
	type args struct {
		manifest *ares.BuildManifest
	}

	manifest1 := &ares.BuildManifest{
		Tag:       "registry:5000/thatcher/pipe:3",
		BaseImage: "debian:stable-slim",
		Components: []*ares.Component{
			{
				FileName: "trigger-1.trigger.0.0.15.ares-v2-demo.so",
				FileId:   "ares-v2-demo/trigger/trigger-1/0.0.15/trigger-1.trigger.0.0.15.ares-v2-demo.so.tar.gz",
			},
		},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "build from manifest",
			args: args{
				manifest: manifest1,
			},
			fields: fields{
				config: &Config{
					URL: "http://localhost:8012/",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewClient(xlogger.NewSimpleLogger(), tt.fields.config)

			if err := m.BuildImage(tt.args.manifest); (err != nil) != tt.wantErr {
				t.Errorf("Client.BuildImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

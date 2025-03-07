package jamfprointegration

import (
	"net/http"
	"reflect"
	"testing"

	"go.uber.org/zap"
)

func TestIntegration_marshalRequest(t *testing.T) {
	type fields struct {
		jamfProFQDN          string
		AuthMethodDescriptor string
		Sugar                *zap.SugaredLogger
		auth                 authInterface
		http                 http.Client
	}
	type args struct {
		body     interface{}
		method   string
		endpoint string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Integration{
				JamfProFQDN:          tt.fields.jamfProFQDN,
				AuthMethodDescriptor: tt.fields.AuthMethodDescriptor,
				Sugar:                tt.fields.Sugar,
				auth:                 tt.fields.auth,
				http:                 tt.fields.http,
			}
			got, err := j.marshalRequest(tt.args.body, tt.args.method, tt.args.endpoint)
			if (err != nil) != tt.wantErr {
				t.Errorf("Integration.marshalRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Integration.marshalRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

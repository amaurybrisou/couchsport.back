package utils

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestParseBody(t *testing.T) {
	type args struct {
		i    interface{}
		body io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "should return true",
			args: args{
				i:    struct{}{},
				body: strings.NewReader(`{"Name":"Amaury", "Age": 32}`),
			},
			want:    map[string]interface{}{"Name": "Amaury", "Age": float64(32)},
			wantErr: false,
		},
		{
			name: "should return true",
			args: args{
				i:    struct{}{},
				body: strings.NewReader(`{"Name":"Amaury", "Age": "33"}`),
			},
			want:    map[string]interface{}{"Name": "Amaury", "Age": "33"},
			wantErr: false,
		},
		{
			name: "malformed JSON, should return nil",
			args: args{
				i:    struct{}{},
				body: strings.NewReader(`{"Name":"Amaury, "Age": "33"}`),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseBody(tt.args.i, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseBody() = %v, want %v", reflect.TypeOf(got), reflect.TypeOf(tt.want))
				t.Errorf("ParseBody() = %s, want %s", got, tt.want)
				t.Errorf("ParseBody() = %v, want %v", got, tt.want)
			}
		})
	}
}

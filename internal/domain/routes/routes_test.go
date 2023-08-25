package routes

import (
	"reflect"
	"testing"
)

func Test_transformPorts(t *testing.T) {
	type args struct {
		ports []string
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "normal input, normal output",
			args: args{ports: []string{"9000", "10000"}},
			want: []int{9000, 10000},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := transformPorts(tt.args.ports); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("transformPorts() = %v, want %v", got, tt.want)
			}
		})
	}
}

package service

import (
	"reflect"
	"testing"
)

func Test_resolveExcludedPorts(t *testing.T) {
	type args struct {
		ports    []int
		excluded []string
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "matched",
			args: args{
				ports:    []int{8080, 8081, 8082},
				excluded: []string{"8081"},
			},
			want: []string{"8080", "8082"},
		},
		{
			name: "not matched",
			args: args{
				ports:    []int{8080, 8081, 8082},
				excluded: []string{"8083"},
			},
			want: []string{"8080", "8081", "8082"},
		},
		{
			name: "nothing to match",
			args: args{
				ports:    []int{8080, 8081, 8082},
				excluded: []string{""},
			},
			want: []string{"8080", "8081", "8082"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveExcludedPorts(tt.args.ports, tt.args.excluded); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resolveExcludedPorts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parsePorts(t *testing.T) {
	type args struct {
		portsStr  string
		toExclude []string
	}

	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "right format. matching",
			args: args{
				portsStr:  "8080...8081",
				toExclude: []string{"8080"},
			},
			want:    []string{"8081"},
			wantErr: false,
		},
		{
			name: "right format. doesn't matching",
			args: args{
				portsStr:  "8080...8081",
				toExclude: []string{"8082"},
			},
			want:    []string{"8080", "8081"},
			wantErr: false,
		},
		{
			name: "right format. nothing to match",
			args: args{
				portsStr:  "8080...8081",
				toExclude: []string{""},
			},
			want:    []string{"8080", "8081"},
			wantErr: false,
		},
		{
			name: "right format. nothing to match. nil safety",
			args: args{
				portsStr:  "8080...8081",
				toExclude: []string{""},
			},
			want:    []string{"8080", "8081"},
			wantErr: false,
		},
		{
			name: "bad format",
			args: args{
				portsStr:  "8080..8081",
				toExclude: []string{""},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "right format. negative numbers",
			args: args{
				portsStr:  "-1...-2",
				toExclude: []string{""},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "right format. equal numbers",
			args: args{
				portsStr:  "8080...8080",
				toExclude: []string{""},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "right format. bigger than uint16 numbers",
			args: args{
				portsStr:  "65535...65537",
				toExclude: []string{""},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePorts(tt.args.portsStr, tt.args.toExclude)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePorts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePorts() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseHostPorts(t *testing.T) {
	type args struct {
		host  string
		ports []string
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "",
			args: args{
				host:  "localhost",
				ports: []string{"8080", "8081"},
			},
			want: []string{"localhost:8080", "localhost:8081"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseHostPorts(tt.args.host, tt.args.ports); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseHostPorts() = %v, want %v", got, tt.want)
			}
		})
	}
}

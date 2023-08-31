package config_test

import (
	"io"
	"loadder/config"
	lb "loadder/load_balancer"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	type args struct {
		file io.Reader
	}

	tests := []struct {
		name    string
		args    args
		want    *config.Config
		wantErr bool
	}{
		{
			name: "right format",
			args: args{file: strings.NewReader(`
algorithm: weighted-round-robin
load-balancer-address: :8080

services:
    first-service:
        name: hello-world-1
        address: http://localhost:9000
        weight: 2
        healthcheck:
           path: /healthcheck
           interval: 1s
           timeout: 2s
           unhealthy-threshold: 2
           timeout-threshold: 2`)},
			want: &config.Config{
				Services: map[string]*config.Service{
					"first-service": {
						Name:    "hello-world-1",
						Address: "http://localhost:9000",
						Weight:  2.0,
						Healthcheck: &config.Healthcheck{
							Path:               "/healthcheck",
							Interval:           "1s",
							Timeout:            "2s",
							UnhealthyThreshold: 2,
							TimeoutThreshold:   2,
						},
					},
				},
				Algorithm:           lb.WightedRoundRobinAlgorithm,
				LoadBalancerAddress: ":8080",
			},
		},
		{
			name: "right format. two services",
			args: args{file: strings.NewReader(`
algorithm: weighted-round-robin
load-balancer-address: :8080

services:
    first-service:
        name: hello-world-1
        address: http://localhost:9000
        weight: 2
        healthcheck:
           path: /healthcheck
           interval: 1s
           timeout: 2s
           unhealthy-threshold: 2
           timeout-threshold: 2
    second-service:
        name: hello-world-1
        address: http://localhost:9001
        weight: 2
        healthcheck:
           path: /healthcheck
           interval: 1s
           timeout: 2s
           unhealthy-threshold: 2
           timeout-threshold: 2`)},
			want: &config.Config{
				Services: map[string]*config.Service{
					"first-service": {
						Name:    "hello-world-1",
						Address: "http://localhost:9000",
						Weight:  2.0,
						Healthcheck: &config.Healthcheck{
							Path:               "/healthcheck",
							Interval:           "1s",
							Timeout:            "2s",
							UnhealthyThreshold: 2,
							TimeoutThreshold:   2,
						},
					},
					"second-service": {
						Name:    "hello-world-1",
						Address: "http://localhost:9001",
						Weight:  2.0,
						Healthcheck: &config.Healthcheck{
							Path:               "/healthcheck",
							Interval:           "1s",
							Timeout:            "2s",
							UnhealthyThreshold: 2,
							TimeoutThreshold:   2,
						},
					},
				},
				Algorithm:           lb.WightedRoundRobinAlgorithm,
				LoadBalancerAddress: ":8080",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := config.Parse(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, *tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHealthcheck_IntervalDuration(t *testing.T) {
	type fields struct {
		Interval string
	}

	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{
			name:   "1ns",
			fields: fields{Interval: "1ns"},
			want:   time.Nanosecond,
		},
		{
			name:   "1us",
			fields: fields{Interval: "1us"},
			want:   time.Microsecond,
		},
		{
			name:   "1ms",
			fields: fields{Interval: "1ms"},
			want:   time.Millisecond,
		},
		{
			name:   "1s",
			fields: fields{Interval: "1s"},
			want:   time.Second,
		},
		{
			name:   "1m",
			fields: fields{Interval: "1m"},
			want:   time.Minute,
		},
		{
			name:   "1h",
			fields: fields{Interval: "1h"},
			want:   time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &config.Healthcheck{
				Interval: tt.fields.Interval,
			}
			if got := h.IntervalDuration(); got != tt.want {
				t.Errorf("IntervalDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestHealthcheck_TimoutDuration(t *testing.T) {
	type fields struct {
		Timout string
	}

	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{
			name:   "1ns",
			fields: fields{Timout: "1ns"},
			want:   time.Nanosecond,
		},
		{
			name:   "1us",
			fields: fields{Timout: "1us"},
			want:   time.Microsecond,
		},
		{
			name:   "1ms",
			fields: fields{Timout: "1ms"},
			want:   time.Millisecond,
		},
		{
			name:   "1s",
			fields: fields{Timout: "1s"},
			want:   time.Second,
		},
		{
			name:   "1m",
			fields: fields{Timout: "1m"},
			want:   time.Minute,
		},
		{
			name:   "1h",
			fields: fields{Timout: "1h"},
			want:   time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &config.Healthcheck{
				Timeout: tt.fields.Timout,
			}
			if got := h.TimoutDuration(); got != tt.want {
				t.Errorf("IntervalDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

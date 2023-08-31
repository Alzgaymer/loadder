package lb

import (
	"reflect"
	"testing"
)

func TestDefineAlgorithm(t *testing.T) {
	type args struct {
		algorithm string
	}

	tests := []struct {
		name    string
		args    args
		want    Algorithm
		wantErr bool
	}{
		{
			name:    "round-robin",
			args:    args{algorithm: "round-robin"},
			want:    NewRoundRobin(),
			wantErr: false,
		},
		{
			name:    "weighted-round-robin",
			args:    args{algorithm: "weighted-round-robin"},
			want:    NewWeightedRoundRobin(),
			wantErr: false,
		},
		{
			name:    "un-weighted-round-robin",
			args:    args{algorithm: "un-weighted-round-robin"},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DefineAlgorithm(tt.args.algorithm)
			if (err != nil) != tt.wantErr {
				t.Errorf("DefineAlgorithm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefineAlgorithm() got = %v, want %v", got, tt.want)
			}
		})
	}
}

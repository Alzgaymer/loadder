package lb

import (
	"reflect"
	"testing"
)

func TestRoundRobin_NextAliveAllAlive(t *testing.T) {
	services := []*Service{
		{
			alive: true,
		},
		{
			alive: true,
		},
	}

	robin := NewRoundRobin()
	robin.Add(services...)

	if !reflect.DeepEqual(robin.NextAlive(), services[0]) {
		t.Errorf("NextAlive() = %v, want %v", robin.NextAlive(), services[0])
	}

	if !reflect.DeepEqual(robin.NextAlive(), services[1]) {
		t.Errorf("NextAlive() = %v, want %v", robin.NextAlive(), services[1])
	}

	if !reflect.DeepEqual(robin.NextAlive(), services[0]) {
		t.Errorf("NextAlive() = %v, want %v", robin.NextAlive(), services[0])
	}
}

func TestRoundRobin_NextAliveNotAllAlive(t *testing.T) {
	services := []*Service{
		{
			alive: true,
		},
		{
			alive: false,
		},
		{
			alive: true,
		},
	}

	robin := NewRoundRobin()
	robin.Add(services...)

	if !reflect.DeepEqual(robin.NextAlive(), services[0]) {
		t.Errorf("NextAlive() = %v, want %v", robin.NextAlive(), services[0])
	}

	if !reflect.DeepEqual(robin.NextAlive(), services[2]) {
		t.Errorf("NextAlive() = %v, want %v", robin.NextAlive(), services[1])
	}

	if !reflect.DeepEqual(robin.NextAlive(), services[0]) {
		t.Errorf("NextAlive() = %v, want %v", robin.NextAlive(), services[0])
	}
}

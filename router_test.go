package servergo

import (
	"testing"
)

func TestNew(t *testing.T) {
	result := New()
  t.Logf("%v", result)
	if result == nil {
		t.Error(`"New()" failed expected a value and returned nil`)
		return
	}

	if result.routes == nil {
		t.Errorf(`"New()" failed router's route list is not initialised`)
	}

	if result.MethodNotAllowedHandler == nil {
		t.Errorf(`"New()" failed router's MethodNotAllowedHandler is not initialised`)
	}

	if result.NotFoundHandler == nil {
		t.Errorf(`"New()" failed router's NotFoundHandler is not initialised`)
	}

	if result.UnauthorizedHandler == nil {
		t.Errorf(`"New()" failed router's UnauthorizedHandler is not initialised`)
	}
}
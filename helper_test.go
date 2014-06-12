package irc

import (
	"reflect"
	"testing"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	defer func() {
		if err := recover(); err != nil {
			if !reflect.DeepEqual(a, b) {
				t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
			}
		}
	}()
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

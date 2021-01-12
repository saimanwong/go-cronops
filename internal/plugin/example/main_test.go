package example

import (
	"testing"

	"github.com/mitchellh/mapstructure"
)

func TestBasic(t *testing.T) {
	param := map[string]interface{}{
		"logprefix": "example",
		"msg":       "hello world",
	}
	var e Example
	err := mapstructure.Decode(param, &e)
	if err != nil {
		t.Errorf("could not decode param")
	}
}

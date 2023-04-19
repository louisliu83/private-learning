package config

import (
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	fileName := "../conf.json"
	file, err := filepath.Abs(fileName)
	if err != nil {
		t.Fail()
	}
	t.Logf("%v", file)
	c, err := Load(file)
	if err != nil {
		t.Fail()
	}
	t.Logf("%v", c)
}

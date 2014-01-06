package config

import "testing"
import "io/ioutil"

func TestOperatingSystem_Default(t *testing.T) {
	operatingSystems := DefaultOperatingSystems()
	v := len(operatingSystems)
	if v < 1 {
		t.Error(
			"For", operatingSystems,
			"expected", 1,
			"got", v,
		)
	}
}

func TestOperatingSystem_WriteDefault(t *testing.T) {
	destination, _ := ioutil.TempFile("", "operatingsystems.json")
	WriteOperatingSystems(destination)
}

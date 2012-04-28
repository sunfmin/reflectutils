package reflectutils

import (
	"testing"
)

func TestPopulate(t *testing.T) {
	var v1 *Person
	Populate(&v1)
	t.Errorf("%+v", v1)

	// var v2 Person
	// Populate(&v2)
}

package reflectutils

import (
	"testing"
)

func TestPopulate(t *testing.T) {
	var v1 *Person
	Populate(&v1)
	if v1.Company == nil {
		t.Errorf("not populated %+v", v1)
	}

}

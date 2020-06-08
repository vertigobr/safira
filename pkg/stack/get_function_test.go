package stack

import "testing"

func TestGetAllFunctions(t *testing.T) {
	functions, err := GetAllFunctions()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(functions)
}

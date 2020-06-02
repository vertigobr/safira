package deploy

import "testing"

func TestCreateYamlFunction(t *testing.T) {
	if err := CreateYamlFunction("./deploy/function.yaml", "func-teste"); err != nil {
		t.Fatal(err)
	}
}

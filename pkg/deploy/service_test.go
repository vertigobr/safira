package deploy

import "testing"

func TestCreateYamlService(t *testing.T) {
	if err := CreateYamlService("./deploy/service.yaml", "func-teste"); err != nil {
		t.Fatal(err)
	}
}

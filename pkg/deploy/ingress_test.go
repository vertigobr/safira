package deploy

import "testing"

func TestCreateYamlIngress(t *testing.T) {
	if err := CreateYamlIngress("./deploy/ingress.yaml", "func-teste"); err != nil {
		t.Fatal(err)
	}
}

package deploy

import "testing"

func TestCreateYamlKongPlugin(t *testing.T) {
	if err := CreateYamlKongPlugin("./deploy/kong_plugin.yaml"); err != nil {
		t.Fatal(err)
	}
}

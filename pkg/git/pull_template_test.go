package git

import "testing"

func TestPullTemplate(t *testing.T) {
	if err := PullTemplate("https://github.com/vertigobr/openfaas-templates.git", true); err != nil {
		t.Fatal(err)
	}
}

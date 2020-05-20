package cmd

import "testing"

func TestCheckRepositoryTemplate(t *testing.T) {
	repo1, official1 := checkRepositoryTemplate("https://github.com/vertigobr/openfaas-templates")
	repo2, official2 := checkRepositoryTemplate("vertigobr/openfaas-templates")
	repo3, official3 := checkRepositoryTemplate("openfaas-templates")

	t.Log(repo1, official1)
	t.Log(repo2, official2)
	t.Log(repo3, official3)
}

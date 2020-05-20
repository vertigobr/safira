package repository

import "testing"

func TestGetRepositoryURL(t *testing.T)  {
	ListTemplates()
	GetRepositoryURL("https://google.com")
	GetRepositoryURL("google.com")
	GetRepositoryURL("google")
	GetRepositoryURL("/foo/bar")
}

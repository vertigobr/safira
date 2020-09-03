// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package git

import (
	"testing"
)

func TestPullTemplate(t *testing.T) {
	if err := PullTemplate("https://github.com/vertigobr/openfaas-templates.git", true, true); err != nil {
		t.Fatal(err)
	}
}

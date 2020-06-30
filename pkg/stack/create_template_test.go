// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package stack

import (
	"testing"
)

func TestCreateTemplate(t *testing.T) {
	f := Function{
		Name:     "func-teste",
		Template: "teste",
		Handler:  "teste.js",
		Image:    "teste:latest",
	}

	if err := CreateTemplate(f); err != nil {
		t.Fatal(err)
	}
}

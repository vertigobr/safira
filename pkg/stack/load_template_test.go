// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the  Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package stack

import (
	"testing"
)

func TestLoadStackFile(t *testing.T) {
	stack, err := LoadStackFile("stack.yml")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(stack)
}

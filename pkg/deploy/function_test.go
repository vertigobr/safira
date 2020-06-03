// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package deploy

import "testing"

func TestCreateYamlFunction(t *testing.T) {
	if err := CreateYamlFunction("./deploy/function.yaml", "func-teste"); err != nil {
		t.Fatal(err)
	}
}

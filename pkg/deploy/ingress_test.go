// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package deploy

import "testing"

func TestCreateYamlIngress(t *testing.T) {
	if err := CreateYamlIngress("./deploy/ingress.yaml", "func-teste"); err != nil {
		t.Fatal(err)
	}
}

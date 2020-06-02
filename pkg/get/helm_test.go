// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the  Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package get

import "testing"

func TestDownloadHelm(t *testing.T) {
	if err := DownloadHelm(); err != nil {
		t.Log(err)
	}
}

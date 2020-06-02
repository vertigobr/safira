// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the  Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package execute

import (
	"testing"
)

func TestExecuteWithStreamStdio(t *testing.T) {
	task := Task{Command: "ls", StreamStdio: true}
	_, err := task.Execute()
	if err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}
}

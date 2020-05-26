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

package get

import "testing"

func TestDownloadKubectl(t *testing.T) {
	if err := DownloadKubectl(); err != nil {
		t.Log(err)
	}
}

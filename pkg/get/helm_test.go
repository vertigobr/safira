package get

import "testing"

func TestDownloadHelm(t *testing.T) {
	if err := DownloadHelm(); err != nil {
		t.Log(err)
	}
}

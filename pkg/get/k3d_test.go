package get

import "testing"

func TestDownloadK3d(t *testing.T) {
	if err := DownloadK3d(); err != nil {
		t.Log(err)
	}
}

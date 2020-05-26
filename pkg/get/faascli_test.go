package get

import "testing"

func TestDownloadFaasCli(t *testing.T) {
	if err := DownloadFaasCli(); err != nil {
		t.Log(err)
	}
}

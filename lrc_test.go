package lyrics_lrc

import (
	"encoding/json"
	"testing"
)

func TestOpenLRCFile(t *testing.T) {
	fl, err := OpenLRCFile("test.lrc")
	if err != nil {
		t.Logf("Error: %v", err)
		return
	}

	jsn, err := json.Marshal(fl.fragments)
	if err != nil {
		t.Logf("Error: %v", err)
		return
	}
	t.Log(string(jsn))
}

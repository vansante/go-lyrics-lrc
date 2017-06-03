package lyrics_lrc

import (
	"fmt"
	"testing"
	"time"
)

func TestLRCTimer(t *testing.T) {
	fl, err := OpenLRCFile("test.lrc")
	if err != nil {
		t.Logf("Error: %v", err)
		return
	}

	tm := NewLRCTimer(fl)

	tm.AddListener(func(start int64, content string, last bool) {
		fmt.Printf("[%v] %s [%v]\n", time.Duration(start)*time.Millisecond, content, last)
	})

	tm.Start()
	//time.Sleep(40 * time.Second)
}

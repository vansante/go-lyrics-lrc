package lyrics_lrc

import "time"

type LRCTimer struct {
	file      *LRCFile
	timer     *time.Timer
	stop      bool
	listeners []func(startTimeMs int64, content string, last bool)
}

func NewLRCTimer(file *LRCFile) (timer *LRCTimer) {
	timer = &LRCTimer{
		file: file,

	}
	return
}

func (t *LRCTimer) AddListener(l func(startTimeMs int64, content string, last bool)) {
	t.listeners = append(t.listeners, l)
}

func (t *LRCTimer) Start() {
	fragments := t.file.fragments

	if len(fragments) < 1 {
		return
	}

	currentIdx := 0
	current := fragments[0]
	startTime := time.Now()
	t.timer = time.NewTimer(time.Millisecond * time.Duration(current.StartTimeMs))
	t.stop = false
	for {
		<-t.timer.C
		if t.stop {
			break
		}

		currentIdx++
		last := currentIdx >= len(fragments)

		for _, l := range t.listeners {
			go l(current.StartTimeMs, current.Content, last)
		}

		if last {
			break
		}

		current = fragments[currentIdx]
		elapsedTime := time.Now().Sub(startTime)
		t.timer.Reset((time.Millisecond * time.Duration(current.StartTimeMs)) - elapsedTime)
	}

	t.timer.Stop()
	t.timer = nil
}

func (t LRCTimer) IsStarted() bool {
	return t.timer != nil
}

func (t LRCTimer) Stop() {
	t.stop = true
	if t.timer != nil {
		t.timer.Stop()
	}
}
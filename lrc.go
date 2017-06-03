package lyrics_lrc

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Parses an LRC file
// https://en.wikipedia.org/wiki/LRC_(file_format)

const LRC_TIME_FORMAT = "04:05.00"

var zeroTime, _ = time.Parse(LRC_TIME_FORMAT, "00:00.00")

type LRCFile struct {
	fragments []LRCFragment
}

type LRCFragment struct {
	StartTimeMs int64
	Content     string
}

func OpenLRCFile(filePath string) (lrcFile *LRCFile, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	lrcFile, err = ReadLRC(file)
	return
}

func ReadLRC(reader io.Reader) (lrcFile *LRCFile, err error) {
	var fragments []LRCFragment

	lineNo := 1
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		var lineFragments []LRCFragment
		lineFragments, err = readLRCLine(scanner.Text(), lineNo)
		if err != nil {
			continue
		}
		for _, fragment := range lineFragments {
			fragments = append(fragments, fragment)
		}
		lineNo++
	}

	sort.Slice(fragments, func(i, j int) bool {
		return fragments[i].StartTimeMs < fragments[j].StartTimeMs
	})

	lrcFile = &LRCFile{
		fragments: fragments,
	}
	return
}

func readLRCLine(line string, lineNo int) (fragments []LRCFragment, err error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	tm, err := parseLRCTime(line, "[", "]")
	if err != nil {
		err = fmt.Errorf("Error on line %d: %v", lineNo, err)
		return
	}

	line = line[len(LRC_TIME_FORMAT)+2:]

	var extraTms []time.Time
	for {
		extraTm, tmErr := parseLRCTime(line, "[", "]")
		if tmErr != nil {
			break
		}
		extraTms = append(extraTms, extraTm)
		line = line[len(LRC_TIME_FORMAT)+2:]
	}

	line = strings.TrimSpace(line)
	if len(extraTms) > 0 {
		fragments = append(fragments, LRCFragment{
			StartTimeMs: getMillisecondsFromTime(tm),
			Content:     line,
		})

		for _, extraTm := range extraTms {
			fragments = append(fragments, LRCFragment{
				StartTimeMs: getMillisecondsFromTime(extraTm),
				Content:     line,
			})
		}
		return
	}

	lineFragments, err := parseContentLine(line, tm)
	for _, fragment := range lineFragments {
		fragments = append(fragments, fragment)
	}
	return
}

func parseLRCTime(line, openChar, closeChar string) (tm time.Time, err error) {
	if len(line) < len(LRC_TIME_FORMAT)+2 {
		err = errors.New("line too short")
		return
	}
	if line[0:1] != openChar || line[len(LRC_TIME_FORMAT)+1:len(LRC_TIME_FORMAT)+2] != closeChar {
		err = errors.New("brackets missing")
		return
	}

	_, err = strconv.Atoi(line[1:3])
	if err != nil {
		// A tag line
		return
	}

	tm, err = time.Parse(LRC_TIME_FORMAT, line[1:len(LRC_TIME_FORMAT)+1])
	return
}

func parseContentLine(line string, tm time.Time) (fragments []LRCFragment, err error) {
	if !strings.Contains(line, "<") {
		fragments = append(fragments, LRCFragment{
			StartTimeMs: getMillisecondsFromTime(tm),
			Content:     line,
		})
		return
	}

	previousTm := tm
	startIndex := 0
	lastIndex := 0
	for {
		idx := strings.Index(line[lastIndex:], "<")
		if idx < 0 {
			break
		}
		idx += lastIndex

		splitTm, tmErr := parseLRCTime(line[idx:], "<", ">")
		if tmErr == nil {
			fragments = append(fragments, LRCFragment{
				StartTimeMs: getMillisecondsFromTime(previousTm),
				Content:     strings.TrimSpace(line[startIndex:idx]),
			})
			startIndex = idx + len(LRC_TIME_FORMAT) + 2
			previousTm = splitTm
		}
		lastIndex = idx + 1
	}

	fragments = append(fragments, LRCFragment{
		StartTimeMs: getMillisecondsFromTime(previousTm),
		Content:     strings.TrimSpace(line[startIndex:]),
	})
	return
}

func getMillisecondsFromTime(tm time.Time) (ms int64) {
	ms = tm.Sub(zeroTime).Nanoseconds() / int64(time.Millisecond)
	return
}

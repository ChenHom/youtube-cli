package chapterizer

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Chapter 代表影片章節的資訊
type Chapter struct {
	Title string
	Start string
	End   string
}

// configuration driven by spec: mode=paragraph, time_gap_ms, transitions, use_embedding, title_generation
var (
	timeGap     = 1500 * time.Millisecond
	transitions = []string{"接下來", "那麼", "然後", "最後", "另外", "另一方面"}
)

// parseTimestamp converts VTT timestamp to time.Duration
func parseTimestamp(ts string) (time.Duration, error) {
	// 清理 ts 字串
	ts = strings.TrimSpace(ts)
	ts = strings.ReplaceAll(ts, "\n", "")
	ts = strings.Trim(ts, "\x00") // Remove null/BOM if any
	parts := strings.Split(ts, ":")
	var hours, mins int
	var secPart string
	switch len(parts) {
	case 3:
		// hh:mm:ss.mmm
		h, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}
		m, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, err
		}
		hours, mins = h, m
		secPart = parts[2]
	case 2:
		// mm:ss.mmm
		m, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}
		hours, mins = 0, m
		secPart = parts[1]
	default:
		return 0, fmt.Errorf("invalid timestamp: %s", ts)
	}
	secParts := strings.Split(secPart, ".")
	if len(secParts) != 2 {
		return 0, fmt.Errorf("invalid seconds format: %s", secPart)
	}
	secs, err := strconv.Atoi(secParts[0])
	if err != nil {
		return 0, err
	}
	millis, err := strconv.Atoi(secParts[1])
	if err != nil {
		return 0, err
	}
	return time.Duration(hours)*time.Hour + time.Duration(mins)*time.Minute + time.Duration(secs)*time.Second + time.Duration(millis)*time.Millisecond, nil
}

// formatTimestamp converts time.Duration back to VTT timestamp string
func formatTimestamp(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	ms := int(d.Milliseconds()) % 1000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}

// DetectParagraphChapters implements paragraph-based chapter detection
func DetectParagraphChapters(transcriptPath string) ([]Chapter, error) {
	f, err := os.Open(transcriptPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	type cue struct {
		start, end time.Duration
		text       string
	}
	var cues []cue
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "-->") {
			parts := strings.Split(line, "-->")
			startDur, err := parseTimestamp(strings.TrimSpace(parts[0]))
			if err != nil {
				return nil, err
			}
			endDur, err := parseTimestamp(strings.TrimSpace(parts[1]))
			if err != nil {
				return nil, err
			}
			if scanner.Scan() {
				text := strings.TrimSpace(scanner.Text())
				cues = append(cues, cue{start: startDur, end: endDur, text: text})
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	var chapters []Chapter
	if len(cues) == 0 {
		return chapters, nil
	}
	var currText strings.Builder
	prev := cues[0]
	currText.WriteString(cues[0].text)
	paragraphStart := prev.start
	for i := 1; i < len(cues); i++ {
		c := cues[i]
		isBoundary := false
		// 條件 A：時間間隔大於
		if c.start-prev.end > timeGap { /* new paragraph boundary */
			isBoundary = true
		}
		// 條件 B：句子開頭為轉折詞
		if !isBoundary {
			for _, t := range transitions {
				if strings.HasPrefix(c.text, t) {
					isBoundary = true
					break
				}
			}
		}
		if isBoundary {
			full := currText.String()
			title := full
			if idx := strings.Index(full, "。"); idx != -1 {
				title = full[:idx+len("。")]
			}
			chapters = append(chapters, Chapter{Title: title, Start: formatTimestamp(paragraphStart), End: formatTimestamp(prev.end)})
			currText.Reset()
			currText.WriteString(c.text)
			paragraphStart = c.start
		} else {
			currText.WriteString(" ")
			currText.WriteString(c.text)
		}
		prev = c
	}
	full := currText.String()
	title := full
	if idx := strings.Index(full, "。"); idx != -1 {
		title = full[:idx+len("。")]
	}
	chapters = append(chapters, Chapter{Title: title, Start: formatTimestamp(paragraphStart), End: formatTimestamp(prev.end)})
	return chapters, nil
}

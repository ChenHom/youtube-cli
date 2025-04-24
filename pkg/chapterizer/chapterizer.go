package chapterizer

import (
	"bufio"
	"os"
	"strings"
)

// Chapter 代表影片章節的資訊
type Chapter struct {
	Title string
	Start string
	End   string
}

// DetectChapters 解析 transcriptPath 中的字幕，回傳章節清單
func DetectChapters(transcriptPath string) ([]Chapter, error) {
	f, err := os.Open(transcriptPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var chapters []Chapter
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "-->") {
			parts := strings.Split(line, "-->")
			start := strings.TrimSpace(parts[0])
			end := strings.TrimSpace(parts[1])
			if scanner.Scan() {
				title := strings.TrimSpace(scanner.Text())
				chapters = append(chapters, Chapter{Title: title, Start: start, End: end})
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return chapters, nil
}

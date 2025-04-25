package pipeline

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/ChenHom/ytcli/pkg/audio"
	"github.com/ChenHom/ytcli/pkg/chapterizer"
	"github.com/ChenHom/ytcli/pkg/downloader"
	"github.com/ChenHom/ytcli/pkg/knowledge"
	"github.com/ChenHom/ytcli/pkg/summarizer"
	"github.com/ChenHom/ytcli/pkg/transcript"
	"github.com/ChenHom/ytcli/pkg/translator"
)

// extractVideoID 從 YouTube URL 解析出影片 ID
func extractVideoID(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return strings.TrimSuffix(filepath.Base(raw), filepath.Ext(raw))
	}
	host := u.Host
	if strings.Contains(host, "youtu.be") {
		return strings.Trim(u.Path, "/")
	}
	if strings.Contains(host, "youtube.com") {
		if v := u.Query().Get("v"); v != "" {
			return v
		}
	}
	return strings.TrimSuffix(filepath.Base(raw), filepath.Ext(raw))
}

// Process 執行完整處理流程
func Process(apiKey, urlStr, outputDir, model string, doChapters, doOverview, doSummary bool, deepDiveChapter, relatedChapter int) error {
	// 建立每次處理專屬子資料夾
	videoID := extractVideoID(urlStr)
	runDir := filepath.Join(outputDir, videoID)
	if err := os.MkdirAll(runDir, os.ModePerm); err != nil {
		return fmt.Errorf("create runDir failed: %w", err)
	}
	// 下載影片與字幕到子資料夾
	videoPath, subtitlePath, err := downloader.Download(urlStr, runDir)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	// 提取音訊
	audioPath := filepath.Join(runDir, "audio.wav")
	if err := audio.Extract(videoPath, audioPath); err != nil {
		return fmt.Errorf("audio extract failed: %w", err)
	}
	// 轉錄或使用內嵌字幕
	transcriptPath := subtitlePath
	if transcriptPath == "" {
		transcriptPath = filepath.Join(runDir, "transcript.vtt")
		if err := transcript.WhisperTranscribe(apiKey, audioPath, model, transcriptPath); err != nil {
			return fmt.Errorf("transcribe failed: %w", err)
		}
	}
	// 翻譯 transcript 為繁體中文
	if err := translator.TranslateFile(apiKey, transcriptPath, "Traditional Chinese"); err != nil {
		return fmt.Errorf("translate transcript failed: %w", err)
	}
	fmt.Println("Transcript:", transcriptPath)
	// 偵測章節
	if doChapters {
		chapters, err := chapterizer.DetectParagraphChapters(transcriptPath)
		if err != nil {
			return fmt.Errorf("detect chapters failed: %w", err)
		}
		data, _ := json.MarshalIndent(chapters, "", "  ")
		path := filepath.Join(runDir, "chapters.json")
		if err := os.WriteFile(path, data, 0644); err != nil {
			return fmt.Errorf("write chapters.json failed: %w", err)
		}
		// 翻譯 chapters.json
		if err := translator.TranslateFile(apiKey, path, "Traditional Chinese"); err != nil {
			return fmt.Errorf("translate chapters failed: %w", err)
		}
		fmt.Println("Chapters saved to", path)
	}
	// Overview
	if doOverview {
		overviewPath := filepath.Join(runDir, "overview.txt")
		if err := summarizer.Overview(apiKey, transcriptPath, overviewPath, model); err != nil {
			return fmt.Errorf("overview failed: %w", err)
		}
		// 翻譯 overview
		if err := translator.TranslateFile(apiKey, overviewPath, "Traditional Chinese"); err != nil {
			return fmt.Errorf("translate overview failed: %w", err)
		}
	}
	// Summary
	if doSummary {
		summaryPath := filepath.Join(runDir, "summary.txt")
		if err := summarizer.Summarize(apiKey, transcriptPath, summaryPath, model); err != nil {
			return fmt.Errorf("summary failed: %w", err)
		}
		// 翻譯 summary
		if err := translator.TranslateFile(apiKey, summaryPath, "Traditional Chinese"); err != nil {
			return fmt.Errorf("translate summary failed: %w", err)
		}
	}
	// Deep Dive
	if deepDiveChapter >= 0 {
		deepdivePath := filepath.Join(runDir, fmt.Sprintf("deep_dive_%d.txt", deepDiveChapter))
		if err := knowledge.DeepDive(apiKey, transcriptPath, deepDiveChapter, deepdivePath); err != nil {
			return fmt.Errorf("deep-dive failed: %w", err)
		}
	}
	// Related Topics
	if relatedChapter >= 0 {
		relatedPath := filepath.Join(runDir, fmt.Sprintf("related_%d.txt", relatedChapter))
		if err := knowledge.Related(apiKey, transcriptPath, relatedChapter, relatedPath); err != nil {
			return fmt.Errorf("related failed: %w", err)
		}
	}

	// 將重點檔案複製到 output 目錄
	patterns := []string{"*.mp4", "*.vtt", "chapters.json", "overview.txt", "summary.txt"}
	for _, pattern := range patterns {
		matches, _ := filepath.Glob(filepath.Join(runDir, pattern))
		for _, src := range matches {
			dst := filepath.Join(outputDir, filepath.Base(src))
			input, err := os.ReadFile(src)
			if err != nil {
				return fmt.Errorf("copy %s failed: %w", src, err)
			}
			if err := os.WriteFile(dst, input, 0644); err != nil {
				return fmt.Errorf("write %s failed: %w", dst, err)
			}
		}
	}
	return nil
}

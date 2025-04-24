package pipeline_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/ChenHom/ytcli/pkg/pipeline"
)

func TestFullPipelineWithAutoTranscript(t *testing.T) {
	// 檢查外部工具
	if _, err := exec.LookPath("yt-dlp"); err != nil {
		t.Skip("yt-dlp not installed, skipping integration test")
	}
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not installed, skipping integration test")
	}
	// 執行 pipeline
	output := t.TempDir()
	url := "https://youtu.be/dQw4w9WgXcQ"
	err := pipeline.Process("dummy-key", url, output, "whisper", true, true, true, -1, -1)
	if err != nil {
		t.Fatalf("Full pipeline failed: %v", err)
	}
	// 驗證輸出
	// 影片檔
	mp4s, _ := filepath.Glob(filepath.Join(output, "*.mp4"))
	if len(mp4s) == 0 {
		t.Errorf("expected video file, found none in %s", output)
	}
	// 字幕檔
	vtts, _ := filepath.Glob(filepath.Join(output, "*.vtt"))
	if len(vtts) == 0 {
		t.Errorf("expected transcript .vtt, found none in %s", output)
	}
	// chapters.json
	if _, err := os.Stat(filepath.Join(output, "chapters.json")); err != nil {
		t.Errorf("chapters.json not found: %v", err)
	}
	// overview.txt
	if _, err := os.Stat(filepath.Join(output, "overview.txt")); err != nil {
		t.Errorf("overview.txt not found: %v", err)
	}
	// summary.txt
	if _, err := os.Stat(filepath.Join(output, "summary.txt")); err != nil {
		t.Errorf("summary.txt not found: %v", err)
	}
}

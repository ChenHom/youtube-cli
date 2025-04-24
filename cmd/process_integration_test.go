package cmd_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestProcessCommandGeneratesAllOutputs(t *testing.T) {
	// 環境需有 yt-dlp 與 ffmpeg
	if _, err := exec.LookPath("yt-dlp"); err != nil {
		t.Skip("yt-dlp not installed, skip process CLI test")
	}
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not installed, skip process CLI test")
	}
	// 設定輸出與 env
	output := t.TempDir()
	cmd := exec.Command("go", "run", filepath.Join("..", "main.go"),
		"process",
		"--url", "https://youtu.be/dQw4w9WgXcQ",
		"--output", output,
		"--chapters", "--overview", "--summary",
	)
	cmd.Env = append(os.Environ(), "OPENAI_API_KEY=dummy-key")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("process command failed: %v", err)
	}
	// 驗證輸出檔案
	files := []string{"*.mp4", "*.vtt", "chapters.json", "overview.txt", "summary.txt"}
	for _, pattern := range files {
		matches, _ := filepath.Glob(filepath.Join(output, pattern))
		if len(matches) == 0 {
			t.Errorf("expected output %s in %s", pattern, output)
		}
	}
}

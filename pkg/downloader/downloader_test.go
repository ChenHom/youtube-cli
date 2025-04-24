package downloader_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/ChenHom/ytcli/pkg/downloader"
)

func TestDownloadVideoAndSubtitle(t *testing.T) {
	if _, err := exec.LookPath("yt-dlp"); err != nil {
		t.Skip("yt-dlp not installed, skipping integration test")
	}
	output := t.TempDir()
	url := "https://youtu.be/dQw4w9WgXcQ"
	video, subtitle, err := downloader.Download(url, output)
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}
	if _, err := os.Stat(video); err != nil {
		t.Errorf("expected video file %s to exist", video)
	}
	// subtitle may not exist if no embedded subs
	if subtitle != "" {
		if _, err := os.Stat(subtitle); err != nil {
			t.Errorf("expected subtitle file %s to exist", subtitle)
		}
	}
}

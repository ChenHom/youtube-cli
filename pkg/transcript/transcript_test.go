package transcript_test

import (
	"os"
	"strings"
	"testing"

	"github.com/ChenHom/ytcli/pkg/transcript"
)

func TestTranscribeGeneratesVTT(t *testing.T) {
	input := "testdata/sample.mp4"
	output := t.TempDir() + "/out.vtt"
	// 若未實作，可能回傳錯誤
	err := transcript.WhisperTranscribe("dummy-key", input, "whisper", output)
	if err != nil {
		t.Fatalf("WhisperTranscribe failed: %v", err)
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("read output vtt failed: %v", err)
	}
	content := string(data)
	if !strings.HasPrefix(content, "WEBVTT") {
		t.Errorf("expected VTT header, got: %s", content)
	}
}

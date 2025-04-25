package chapterizer_test

import (
	"testing"

	"github.com/ChenHom/ytcli/pkg/chapterizer"
)

func TestDetectChaptersFromTranscript(t *testing.T) {
	transcript := "../../testdata/sample.vtt"
	chapters, err := chapterizer.DetectParagraphChapters(transcript)
	if err != nil {
		t.Fatalf("DetectParagraphChapters failed: %v", err)
	}
	if len(chapters) == 0 {
		t.Errorf("expected at least one chapter, got 0")
	}
	for _, ch := range chapters {
		if ch.Start == "" || ch.End == "" || ch.Title == "" {
			t.Errorf("invalid chapter data: %+v", ch)
		}
	}
}

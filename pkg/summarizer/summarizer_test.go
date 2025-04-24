package summarizer_test

import (
	"os"
	"strings"
	"testing"

	"github.com/ChenHom/ytcli/pkg/summarizer"
)

func TestOverviewListsKeySentences(t *testing.T) {
	input := "../../testdata/sample.vtt"
	output := t.TempDir() + "/overview.txt"
	err := summarizer.Overview("dummy-key", input, output)
	if err != nil {
		t.Fatalf("Overview failed: %v", err)
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("read overview failed: %v", err)
	}
	content := string(data)
	if !strings.HasPrefix(content, "Overview:") {
		t.Errorf("expected overview prefix, got: %s", content)
	}
}

func TestSummaryProducesConciseText(t *testing.T) {
	input := "../../testdata/sample.vtt"
	output := t.TempDir() + "/summary.txt"
	err := summarizer.Summarize("dummy-key", input, output)
	if err != nil {
		t.Fatalf("Summarize failed: %v", err)
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("read summary failed: %v", err)
	}
	content := string(data)
	// placeholder content is short
	words := strings.Fields(content)
	if len(words) >= 500 {
		t.Errorf("expected summary < 500 words, got %d", len(words))
	}
	if !strings.Contains(content, "Summary placeholder") {
		t.Errorf("expected placeholder summary, got: %s", content)
	}
}

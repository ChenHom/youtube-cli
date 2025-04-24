package knowledge_test

import (
	"os"
	"strings"
	"testing"

	"github.com/ChenHom/ytcli/pkg/knowledge"
)

func TestDeepDiveExpandsChapter(t *testing.T) {
	out := t.TempDir() + "/deep.txt"
	err := knowledge.DeepDive("dummy-key", "testdata/sample.vtt", 1, out)
	if err != nil {
		t.Fatalf("DeepDive failed: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read deep-dive output failed: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "Deep dive") {
		t.Errorf("expected deep-dive placeholder, got: %s", content)
	}
}

func TestRelatedTopicsRetrieved(t *testing.T) {
	out := t.TempDir() + "/related.txt"
	err := knowledge.Related("dummy-key", "testdata/sample.vtt", 1, out)
	if err != nil {
		t.Fatalf("Related failed: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read related output failed: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "Related topics") {
		t.Errorf("expected related placeholder, got: %s", content)
	}
}

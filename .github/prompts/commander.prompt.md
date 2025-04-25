參照以下說明格式來實作

```json
{
  "commands": {
    "download": {
      "description": "下載指定 YouTube 影片檔與內嵌字幕（若有），並儲存到目錄。",
      "inputs": {
        "flags": {
          "url": { "type": "string", "required": true, "description": "YouTube 影片網址" },
          "output": { "type": "string", "required": false, "default": "./videos", "description": "存放影片與字幕的資料夾" }
        }
      },
      "outputs": [
        { "path": "output/<videoID>.mp4", "description": "Downloaded video file" },
        { "path": "output/<videoID>.vtt", "description": "Downloaded or auto-generated subtitle file" }
      ],
      "dependencies": {
        "external": ["yt-dlp"],
        "go_packages": ["github.com/spf13/cobra"]
      },
      "algorithm_pseudocode": [
        "parse flags",
        "if yt-dlp not in PATH → error",
        "run shell: yt-dlp -o \"{output}/%(id)s.%(ext)s\" --write-auto-sub --sub-lang en {url}",
        "if exit code != 0 → retry up to 2 times, then return error"
      ],
      "function_signature": "func DownloadVideo(url string, outputDir string) error",
      "error_handling": [
        { "condition": "yt-dlp not found", "action": "return error 'yt-dlp not found'" },
        { "condition": "download failure", "action": "retry twice then return error" }
      ],
      "example_code": "err := downloader.DownloadVideo(\"https://youtu.be/XXX\", \"./videos\")",
      "tests": [
        {
          "name": "TestDownloadVideo_Success",
          "type": "integration",
          "input": { "url": "https://youtu.be/dQw4w9WgXcQ", "outputDir": "test_videos" },
          "validate": ["file exists test_videos/dQw4w9WgXcQ.mp4", "subtitle file exists if embedded"]
        }
      ],
      "performance": {
        "notes": "可考慮非同步呼叫 yt-dlp，並行多影片下載"
      },
      "documentation": {
        "godoc": "Yes",
        "readme_section": "Included usage example in README.md"
      }
    },
    "transcribe": {
      "description": "抽取音訊並用 Whisper 轉文字，產生字幕檔。",
      "inputs": {
        "flags": {
          "input": { "type": "string", "required": true, "description": "影片檔案路徑" },
          "model": { "type": "string", "required": false, "default": "whisper", "description": "Whisper 模型名稱" }
        }
      },
      "outputs": [
        { "path": "<input_basename>.vtt", "description": "Generated subtitle file with WEBVTT header" }
      ],
      "dependencies": {
        "external": ["ffmpeg"],
        "go_packages": ["github.com/openai/whisper-go"]
      },
      "algorithm_pseudocode": [
        "extract audio via ffmpeg",
        "call Whisper API with audio stream",
        "write returned text segments into .vtt format"
      ],
      "function_signature": "func TranscribeVideo(inputFile string, model string) (string, error)",
      "error_handling": [
        { "condition": "ffmpeg not found", "action": "return error 'ffmpeg not found'" },
        { "condition": "Whisper API error", "action": "return API error message" }
      ],
      "example_code": "vttFile, err := transcript.TranscribeVideo(\"sample.mp4\", \"whisper\")",
      "tests": [
        {
          "name": "TestTranscribeGeneratesVTT",
          "type": "unit",
          "input": { "videoFile": "testdata/sample.mp4", "model": "whisper" },
          "validate": ["output file sample.vtt contains 'WEBVTT' header"]
        }
      ],
      "performance": {
        "notes": "轉錄大影片可能耗時，可加入進度顯示與分段並行處理"
      },
      "documentation": {
        "godoc": "Yes",
        "readme_section": "Described in README under 'transcribe' command"
      }
    },
    "chapters": {
      "description": "偵測章節並輸出時間區間。",
      "inputs": {
        "flags": {
          "transcript": { "type": "string", "required": true, "description": "字幕檔案 .vtt 路徑" }
        }
      },
      "outputs": [
        { "path": "stdout or chapters.json", "description": "List of chapters with start/end times and titles" }
      ],
      "dependencies": {
        "go_packages": ["github.com/yourusername/ytcli/pkg/chapterizer"]
      },
      "algorithm_pseudocode": [
        "load transcript segments",
        "slide window over segments to detect semantic breaks",
        "mark chapter boundaries at high-break points",
        "assign titles based on segment text clustering"
      ],
      "function_signature": "func DetectChapters(transcriptFile string) ([]Chapter, error)",
      "error_handling": [
        { "condition": "file not found", "action": "return error" },
        { "condition": "parse error", "action": "return error" }
      ],
      "example_code": "chapters, err := chapterizer.DetectChapters(\"sample.vtt\")",
      "tests": [
        {
          "name": "TestDetectChaptersFromTranscript",
          "type": "unit",
          "input": { "transcriptFile": "testdata/sample.vtt" },
          "validate": ["returned slice length > 0", "each Chapter has valid StartTime < EndTime"]
        }
      ],
      "performance": {
        "notes": "可設定最小章節長度與最大章節數以優化速度"
      },
      "documentation": {
        "godoc": "Yes",
        "readme_section": "Chapters command details"
      }
    },
    "overview": {
      "description": "列出章節標題、開始時間及重點句。",
      "inputs": {
        "flags": {
          "transcript": { "type": "string", "required": true, "description": "字幕檔案 .vtt 路徑" }
        }
      },
      "outputs": [
        { "path": "stdout or overview.txt", "description": "Each chapter with timestamp and key sentence" }
      ],
      "dependencies": {
        "go_packages": ["github.com/yourusername/ytcli/pkg/summarizer"]
      },
      "algorithm_pseudocode": [
        "load chapters via DetectChapters",
        "for each chapter: send transcript segment to GPT",
        "receive key sentence, print with timestamp"
      ],
      "function_signature": "func GenerateOverview(transcriptFile string) (string, error)",
      "error_handling": [
        { "condition": "API error", "action": "retry or return error" }
      ],
      "example_code": "overview, err := summarizer.GenerateOverview(\"sample.vtt\")",
      "tests": [
        {
          "name": "TestOverviewListsKeySentences",
          "type": "unit",
          "input": { "transcriptFile": "testdata/sample.vtt" },
          "validate": ["output contains timestamps", "contains key sentences"]
        }
      ],
      "performance": {
        "notes": "可批次合併多章節請求以減少 API 呼叫次數"
      },
      "documentation": {
        "godoc": "Yes",
        "readme_section": "Overview command usage"
      }
    },
    "summary": {
      "description": "整支影片或指定章節文字摘要。",
      "inputs": {
        "flags": {
          "transcript": { "type": "string", "required": true, "description": "字幕檔案 .vtt 路徑" },
          "chapter": { "type": "int", "required": false, "description": "指定章節編號" }
        }
      },
      "outputs": [
        { "path": "stdout or summary.txt", "description": "Generated summary text" }
      ],
      "dependencies": {
        "go_packages": ["github.com/yourusername/ytcli/pkg/summarizer"]
      },
      "algorithm_pseudocode": [
        "if chapter flag set: extract that segment text else entire transcript",
        "send to GPT summarization endpoint",
        "return summary text"
      ],
      "function_signature": "func GenerateSummary(transcriptFile string, chapterIndex int) (string, error)",
      "error_handling": [
        { "condition": "empty transcript", "action": "return error 'no transcript data'" }
      ],
      "example_code": "summary, err := summarizer.GenerateSummary(\"sample.vtt\", 0)",
      "tests": [
        {
          "name": "TestSummaryProducesConciseText",
          "type": "unit",
          "input": { "transcriptFile": "testdata/sample.vtt" },
          "validate": ["output length < 500 words", "contains main topic keywords"]
        }
      ],
      "performance": {
        "notes": "限制最大文字長度以符合 API 請求限制"
      },
      "documentation": {
        "godoc": "Yes",
        "readme_section": "Summary command details"
      }
    },
    "deep-dive": {
      "description": "對指定章節做詳盡解析與背景延伸。",
      "inputs": {
        "flags": {
          "chapter": { "type": "int", "required": true, "description": "章節編號" },
          "transcript": { "type": "string", "required": true, "description": "字幕檔案 .vtt 路徑" }
        }
      },
      "outputs": [
        { "path": "stdout or deepdive.txt", "description": "Detailed analysis text" }
      ],
      "dependencies": {
        "go_packages": ["github.com/yourusername/ytcli/pkg/summarizer", "github.com/yourusername/ytcli/pkg/knowledge"]
      },
      "algorithm_pseudocode": [
        "extract chapter transcript",
        "send to GPT with prompt for deep explanation",
        "retrieve related knowledge via knowledge module",
        "combine and return"
      ],
      "function_signature": "func DeepDiveChapter(transcriptFile string, chapterIndex int) (string, error)",
      "error_handling": [
        { "condition": "invalid chapter index", "action": "return error 'chapter not found'" }
      ],
      "example_code": "detail, err := summarizer.DeepDiveChapter(\"sample.vtt\", 2)",
      "tests": [
        {
          "name": "TestDeepDiveExpandsChapter",
          "type": "unit",
          "input": { "chapterIndex": 1, "transcriptFile": "testdata/sample.vtt" },
          "validate": ["output contains background context", "output contains examples"]
        }
      ],
      "performance": {
        "notes": "可分段呼叫 API 以避免超時"
      },
      "documentation": {
        "godoc": "Yes",
        "readme_section": "Deep-dive command usage"
      }
    },
    "related": {
      "description": "列出與該章節相關的其他知識主題。",
      "inputs": {
        "flags": {
          "chapter": { "type": "int", "required": true, "description": "章節編號" },
          "transcript": { "type": "string", "required": true, "description": "字幕檔案 .vtt 路徑" }
        }
      },
      "outputs": [
        { "path": "stdout or related.txt", "description": "List of related topic titles" }
      ],
      "dependencies": {
        "go_packages": ["github.com/yourusername/ytcli/pkg/knowledge"]
      },
      "algorithm_pseudocode": [
        "extract chapter keywords via GPT",
        "query Wikipedia or vector DB with keywords",
        "return top 3–5 topic titles and summaries"
      ],
      "function_signature": "func GetRelatedTopics(transcriptFile string, chapterIndex int) ([]Topic, error)",
      "error_handling": [
        { "condition": "no related topics found", "action": "return empty list" }
      ],
      "example_code": "topics, err := knowledge.GetRelatedTopics(\"sample.vtt\", 2)",
      "tests": [
        {
          "name": "TestRelatedTopicsRetrieved",
          "type": "unit",
          "input": { "chapterIndex": 1, "transcriptFile": "testdata/sample.vtt" },
          "validate": ["len(topics) >= 3"]
        }
      ],
      "performance": {
        "notes": "將查詢結果快取以減少重複 API 呼叫"
      },
      "documentation": {
        "godoc": "Yes",
        "readme_section": "Related command details"
      }
    },
    "process": {
      "description": "一鍵執行 download →（如需）transcribe → chapters/overview/summary/deep-dive/related",
      "inputs": {
        "flags": {
          "url": { "type": "string", "required": true },
          "output": { "type": "string", "required": true },
          "model": { "type": "string", "required": false },
          "chapters": { "type": "boolean", "required": false },
          "overview": { "type": "boolean", "required": false },
          "summary": { "type": "boolean", "required": false },
          "deep-dive": { "type": "int", "required": false },
          "related": { "type": "int", "required": false }
        }
      },
      "outputs": [
        "video file",
        "transcript file",
        "chapters.json",
        "overview.txt",
        "summary.txt",
        "deepdive.txt",
        "related.txt"
      ],
      "dependencies": {
        "all_commands": true
      },
      "algorithm_pseudocode": [
        "call DownloadVideo",
        "if no subtitle → call TranscribeVideo",
        "if chapters flag → call DetectChapters",
        "if overview flag → call GenerateOverview",
        "if summary flag → call GenerateSummary",
        "if deep-dive flag → call DeepDiveChapter",
        "if related flag → call GetRelatedTopics"
      ],
      "function_signature": "func RunFullPipeline(opts ProcessOptions) error",
      "error_handling": [
        { "condition": "any step error", "action": "abort and report step name and error" }
      ],
      "example_code": "err := pipeline.RunFullPipeline(opts)",
      "tests": [
        {
          "name": "TestFullPipelineWithAutoTranscript",
          "type": "integration",
          "input": { "url": "https://youtu.be/dQw4w9WgXcQ", "output": "test_videos", "flags": ["--chapters","--overview","--summary"] },
          "validate": ["video file exists", "transcript exists", "chapters.json exists", "overview.txt exists", "summary.txt exists"]
        }
      ],
      "performance": {
        "notes": "可針對各步驟並行或分段執行以提高效率"
      },
      "documentation": {
        "godoc": "Yes",
        "readme_section": "Process command usage"
      }
    }
  }
}
```
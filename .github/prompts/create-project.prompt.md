```json
{
  "project": {
    "name": "ytcli",
    "description": "YouTube AI CLI：下載影片、轉錄、章節偵測、摘要、深入與相關知識檢索",
    "language": "go",
    "module": "github.com/yourusername/ytcli"
  },
  "structure": {
    "cmd": [
      "process.go",
      "download.go",
      "transcribe.go",
      "chapters.go",
      "overview.go",
      "summary.go",
      "deepdive.go",
      "related.go"
    ],
    "pkg": {
      "pipeline": ["pipeline.go"],
      "downloader": ["downloader.go"],
      "audio": ["audio.go"],
      "transcript": ["whisper.go"],
      "chapterizer": ["chapterizer.go"],
      "summarizer": ["summarizer.go"],
      "knowledge": ["knowledge.go"],
      "config": ["config.go"]
    },
    "root": ["go.mod", "README.md"]
  },
  "commands": {
    "download": {
      "flags": ["--url <YouTube_URL>", "--output <path>"],
      "description": "下載影片檔與內嵌字幕（若有）"
    },
    "transcribe": {
      "flags": ["--input <video.mp4>", "--model whisper"],
      "description": "抽取音訊並用 Whisper 轉文字，產生字幕檔"
    },
    "chapters": {
      "flags": ["--transcript <file.vtt>"],
      "description": "偵測章節並輸出時間區間"
    },
    "overview": {
      "flags": ["--transcript <file.vtt>"],
      "description": "列出章節標題、開始時間及重點句"
    },
    "summary": {
      "flags": ["--transcript <file.vtt>"],
      "description": "整支影片或指定章節文字摘要"
    },
    "deep-dive": {
      "flags": ["--chapter <n>", "--transcript <file.vtt>"],
      "description": "對指定章節做詳盡解析與背景延伸"
    },
    "related": {
      "flags": ["--chapter <n>", "--transcript <file.vtt>"],
      "description": "列出與該章節相關的其他知識主題"
    },
    "process": {
      "flags": [
        "--url <YouTube_URL>",
        "--output <path>",
        "[--model whisper]",
        "[--chapters]",
        "[--overview]",
        "[--summary]",
        "[--deep-dive <章節>]",
        "[--related <章節>]"
      ],
      "description": "一鍵執行 download →（如需）transcribe → chapters/overview/summary/deep-dive/related"
    }
  },
  "dependencies": {
    "cobra": "github.com/spf13/cobra",
    "viper": "github.com/spf13/viper",
    "yt-dlp": "external CLI",
    "ffmpeg": "external CLI",
    "openai": "OpenAI Go SDK (for GPT & Whisper)"
  }
}

```
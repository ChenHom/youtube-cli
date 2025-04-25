# ytcli

YouTube AI CLI：下載影片、轉錄、章節偵測、摘要、深入與相關知識檢索

## 安裝

```bash
# clone
git clone https://github.com/ChenHom/ytcli.git
cd ytcli

# 安裝相依套件
go install github.com/ChenHom/ytcli/cmd/ytcli@latest
```

## 使用範例

### 下載影片

```bash
ytcli download --url <YouTube_URL> --output <path>
```

### 轉錄字幕

```bash
ytcli transcribe --input <video.mp4> [--model whisper]
```

### 偵測章節

```bash
ytcli chapters --transcript <file.vtt>
```

### 列出章節概要

```bash
ytcli overview --transcript <file.vtt>
```

### 生成摘要

```bash
ytcli summary --transcript <file.vtt>
```

### 深入探討章節

```bash
ytcli deep-dive --transcript <file.vtt> --chapter <n>
```

### 檢索相關知識

```bash
ytcli related --query <topic>
```

### 一鍵全流程

```bash
ytcli process --url <YouTube_URL> --output <path> [--model whisper] [--chapters] [--overview] [--summary] [--deep-dive <n>] [--related <topic>]
```

## 全流程處理範例

```bash
ytcli process --url https://youtu.be/dQw4w9WgXcQ --output ./output --chapters --overview --summary
```

執行後將依序下載影片、轉錄字幕、偵測章節並生成 chapters.json、列出章節重點（overview.txt）、生成摘要（summary.txt）。

## 相依套件

- github.com/spf13/cobra
- github.com/spf13/viper
- OpenAI Go SDK (`github.com/sashabaranov/go-openai`)
- 需要外部命令: `yt-dlp`, `ffmpeg`

## 專案結構

```
cmd/
  download.go
  transcribe.go
  chapters.go
  overview.go
  summary.go
  deepdive.go
  related.go
  process.go
  root.go
pkg/
  pipeline/pipeline.go
  downloader/downloader.go
  audio/audio.go
  transcript/whisper.go
  chapterizer/chapterizer.go
  summarizer/summarizer.go
  knowledge/knowledge.go
  config/config.go
go.mod
README.md
```

package downloader

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// 將 youtu.be 或 youtube.com URL 正規化為只含 video ID 的 watch URL
func normalizeYouTubeURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	host := u.Host
	if strings.Contains(host, "youtu.be") {
		id := strings.Trim(u.Path, "/")
		return "https://www.youtube.com/watch?v=" + id
	}
	if strings.Contains(host, "youtube.com") {
		q := u.Query()
		if v := q.Get("v"); v != "" {
			return "https://www.youtube.com/watch?v=" + v
		}
	}
	return raw
}

// Download 使用 yt-dlp 下載影片與字幕，回傳影片檔與字幕檔路徑
func Download(urlStr, outputDir string) (videoPath, subtitlePath string, err error) {
	// 建立輸出資料夾
	if err = os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", "", err
	}
	tpl := filepath.Join(outputDir, "%(id)s.%(ext)s")
	// 正規化 URL，去掉多餘參數
	urlStr = normalizeYouTubeURL(urlStr)
	cmd := exec.Command("yt-dlp",
		"-f", "bestvideo+bestaudio",
		"--merge-output-format", "mp4",
		"--write-subs",
		"--sub-format", "vtt/best",
		"-o", tpl,
		urlStr,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return
	}
	// 搜尋下載結果
	mp4s, _ := filepath.Glob(filepath.Join(outputDir, "*.mp4"))
	if len(mp4s) > 0 {
		videoPath = mp4s[0]
	}
	vtts, _ := filepath.Glob(filepath.Join(outputDir, "*.vtt"))
	if len(vtts) > 0 {
		subtitlePath = vtts[0]
	}
	fmt.Printf("Download complete: video=%s, subtitle=%s\n", videoPath, subtitlePath)
	return
}

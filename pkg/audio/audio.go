package audio

import (
	"io"
	"os/exec"
)

// Extract 將影片檔 input 轉為 WAV 音訊，並存至 output
func Extract(input, output string) error {
	cmd := exec.Command("ffmpeg", "-y", "-i", input, "-vn", "-acodec", "pcm_s16le", "-ar", "16000", "-ac", "1", output)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	return cmd.Run()
}

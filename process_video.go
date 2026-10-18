package main

import (
	"fmt"
	"os"
	"os/exec"
)

func processVideoForFastStart(filePath string) (string, error) {
	outputPath := filePath + ".processing"

	cmd := exec.Command("ffmpeg",
		"-nostdin",
		"-i", filePath,
		"-c", "copy",
		"-movflags", "+faststart",
		"-f", "mp4",
		outputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		os.Remove(outputPath)
		return "", fmt.Errorf("ffmpeg failed: %w: %s", err, output)
	}

	return outputPath, nil
}

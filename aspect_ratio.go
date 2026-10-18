package main

import (
	"bytes"
	"encoding/json"
	"math"
	"os/exec"
)

type ffprobeResult struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"streams"`
}

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-print_format", "json",
		"-show_streams", filePath,
	)

	buf := bytes.Buffer{}
	cmd.Stdout = &buf

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	ffprobeResult := ffprobeResult{}

	if err := json.Unmarshal(buf.Bytes(), &ffprobeResult); err != nil {
		return "", err
	}

	width := ffprobeResult.Streams[0].Width
	height := ffprobeResult.Streams[0].Height

	aspectRatio := calcAspecRatio(width, height)

	return aspectRatio, nil
}

func calcAspecRatio(width int, height int) string {
	const tolerance = 0.1
	const landscape = 16.0 / 9.0
	const portrait = 9.0 / 16.0
	rasio := (float64(width) / float64(height))

	if math.Abs(landscape-rasio) < tolerance {
		return "16:9"
	}
	if math.Abs(portrait-rasio) < tolerance {
		return "9:16"
	}
	return "other"
}

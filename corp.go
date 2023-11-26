package main

import (
	"fmt"
	"os"
	"os/exec"
)

type Corp struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Format string  `json:"format"`
}

func (c *Corp) isEmpty() bool {
	return c.X == 0 &&
		c.Y == 0 &&
		c.Width == 0 &&
		c.Height == 0
}

func (c *Corp) Crop(tempFile *os.File, outFilename string) error {
	args := []string{"-i", tempFile.Name()}
	if c.isEmpty() == false {
		cropArgs := fmt.Sprintf("crop=%d:%d:%d:%d", int(c.Width), int(c.Height), int(c.X), int(c.Y))
		args = append(args, "-filter:v", cropArgs)
	}

	args = append(args, outFilename)

	return exec.Command("ffmpeg", args...).Run()
}

func (c *Corp) GetFormat() string {
	if c.Format == "" {
		return "mp4"
	}
	return c.Format
}

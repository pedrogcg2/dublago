package media

import (
	"bytes"
	"log/slog"
	"os/exec"
	"path/filepath"
)

type FfmpegWrapper struct{}

func (w *FfmpegWrapper) Merge(inputVideoPath string, inputAudioPath string, outputVideoPath string) error {
	slog.Info("[MEDIA] Start merge video and audio files")

	args := []string{"-i", inputAudioPath, "-i", inputVideoPath, "-c", "copy", outputVideoPath, "-y"}
	cmd := exec.Command("ffmpeg", args...)

	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		slog.Error("[MEDIA] Error on merge files: " + err.Error())
		return err
	}

	// TODO:the ffmpeg output log are too long.
	// retrieve just usable info.
	// slog.Info(stdErr.String())

	absOutPath, outPathErr := filepath.Abs(outputVideoPath)

	outValue := absOutPath
	if outPathErr != nil {
		outValue = outputVideoPath
	}

	slog.Info("[MEDIA] Successfully merged files into path " + outValue)

	return err
}

func (w *FfmpegWrapper) Unmerge(inputVideoPath string, outputVideoPath string, outputAudioPath string) error {
	inputValue := inputVideoPath

	if inputAbsValue, error := filepath.Abs(inputVideoPath); error == nil {
		inputValue = inputAbsValue
	}

	slog.Info("[MERGE] Start unmerge audio and video from " + inputValue)

	args := []string{"-i", inputVideoPath, "-an", "-c:v", "copy", outputVideoPath, "-vn", "-c:a", "copy", outputAudioPath, "-y"}
	cmd := exec.Command("ffmpeg", args...)

	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		slog.Error("[MEDIA] Error on unmerge files: " + err.Error())
		return err
	}

	// TODO:the ffmpeg output logs are too long.
	// retrieve just usable info.
	// slog.Info(stdErr.String())

	slog.Info("[MEDIA] Successfully unmerge audio and video")

	return err
}

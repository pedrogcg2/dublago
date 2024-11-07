package tts

import (
	"bytes"
	"fmt"
	"log/slog"
	"os/exec"
	"time"
	"tradutor-dos-crias/shared"

	"github.com/google/uuid"
)

type CoquiTTS struct{}

func (c CoquiTTS) Speech(message string) (string, error) {
	slog.Info("[TTS] Generating speech now")
	modelName := "tts_models/multilingual/multi-dataset/xtts_v2"

	speakerWavPath := shared.PWD + "pipe/primegean-speaker.wav"

	t := time.Now()
	outFile := fmt.Sprintf("%d-%d-%d-%s-tts.wav", t.Day(), t.Month(), t.Year(), uuid.New())
	outPath := shared.PWD + "pipe/" + outFile

	languageIdx := "pt"

	cmd := exec.Command("tts", "--model_name", modelName,
		"--text", message,
		"--speaker_wav", speakerWavPath,
		"--out_path", outPath,
		"--language_idx", languageIdx)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		slog.Error(`[TTS] Error in running 'cmd.Run': ` + err.Error() +
			"\n\t\t\tCLI Error Output: " + stderr.String())
		return "", err
	}

	slog.Info(`[TTS] Speech generated with success
	Output path: ` + outPath)

	return outPath, nil
}

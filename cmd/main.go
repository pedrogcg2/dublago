package main

import (
	"io"
	"log/slog"
	"os"
	"time"
	"tradutor-dos-crias/caption"
	"tradutor-dos-crias/input"
	"tradutor-dos-crias/media"
	"tradutor-dos-crias/pipeline"
	"tradutor-dos-crias/transcript"
	"tradutor-dos-crias/translator"
	"tradutor-dos-crias/tts"

	"github.com/google/uuid"
)

func main() {
	mediaHandler := &media.FfmpegWrapper{}
	transcripter := transcript.NewWhisper()
	translator := &translator.GoogleTranslator{}
	//TODO: Change to kokoro, coqui is a very problematic lib
	tts := tts.NewCoquiTTS()
	subtitler := caption.NewStablets()
	downloader := new(input.YouTube)

	pipeline := pipeline.NewPipeline(transcripter, mediaHandler, translator, tts, subtitler, downloader)
	start := time.Now()

	file, _ := os.Open("pipe/story.md")
	defer func() {
		if err := file.Close(); err != nil {
			slog.Error(err.Error())
		}
	}()

	tBuff, _ := io.ReadAll(file)
	output := "output/" + uuid.NewString() + ".mp4"
	text := string(tBuff[:])
	pipeline.RunFromText(text, "pipe/parkour.webm", output)
	slog.Info("[TIME SPEND] " + time.Since(start).String())
}

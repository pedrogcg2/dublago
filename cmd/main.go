package main

import (
	"log/slog"
	"time"
	"tradutor-dos-crias/caption"
	"tradutor-dos-crias/media"
	"tradutor-dos-crias/pipeline"
	"tradutor-dos-crias/transcript"
	"tradutor-dos-crias/translator"
	"tradutor-dos-crias/tts"
)

func main() {
	mediaHandler := &media.FfmpegWrapper{}
	transcripter := transcript.NewWhisper()
	translator := &translator.GoogleTranslator{}
	tts := tts.NewCoquiTTS()
	subtitler := caption.NewStablets()

	pipeline := pipeline.NewPipeline(transcripter, mediaHandler, translator, tts, subtitler)
	start := time.Now()
	pipeline.Run("pipe/videoDubbed.mp4")
	slog.Info("[TIME SPEND] " + time.Since(start).String())
}

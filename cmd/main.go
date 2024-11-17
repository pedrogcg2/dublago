package main

import (
	"tradutor-dos-crias/media"
	"tradutor-dos-crias/pipeline"
	"tradutor-dos-crias/transcript"
	"tradutor-dos-crias/translator"
	"tradutor-dos-crias/tts"
)

func main() {
	transcripter := transcript.NewWhisper()
	mediaHandler := &media.FfmpegWrapper{}
	translator := translator.GoogleTranslator{}
	tts := tts.NewCoquiTTS()

	pipeline := pipeline.NewPipeline(transcripter, mediaHandler, translator, tts)

	pipeline.Run("pipe/video.mp4", "pipe/videoDubbed.mp4")
}

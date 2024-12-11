package pipeline

import (
	"log/slog"
	"os"
	"path/filepath"
	"tradutor-dos-crias/caption"
	"tradutor-dos-crias/input"
	"tradutor-dos-crias/media"
	"tradutor-dos-crias/transcript"
	"tradutor-dos-crias/translator"
	"tradutor-dos-crias/tts"

	"github.com/google/uuid"
)

type Pipeline struct {
	transcripter transcript.Transcripter
	mediaHandler media.MediaHandler
	translator   translator.Translator
	speaker      tts.Speecher
	subtitler    caption.Subtitler
}

func NewPipeline(transcripter transcript.Transcripter,
	mediaHandler media.MediaHandler,
	translator translator.Translator,
	speaker tts.Speecher,
	subtitler caption.Subtitler,
) *Pipeline {
	i := &Pipeline{
		transcripter,
		mediaHandler,
		translator,
		speaker,
		subtitler,
	}

	return i
}

func (i *Pipeline) Run(outputVideoPath string) error {
	filesToRemove := make([]string, 0)
	defer removeTmpFiles(&filesToRemove)

	url, err := input.Parse()
	if err != nil {
		slog.Error("[INPUT] Error in getting input from cli. Err: " + err.Error())
		return err
	}

	inputVideoPath, inputAudioPath := input.Download(url)
	filesToRemove = append(filesToRemove, inputVideoPath, inputAudioPath)

	outputFolderDefault, err := filepath.Abs("pipe/")
	if err != nil {
		return err
	}

	text, err := i.transcripter.Transcript(inputAudioPath)
	filesToRemove = append(filesToRemove, "audio.txt")
	if err != nil {
		return err
	}

	translatedText, err := i.translator.Translate(text)
	if err != nil {
		return err
	}

	dubbedAudio, err := i.speaker.Speech(translatedText)
	if err != nil {
		return err
	}

	tmpDubbedFileName := outputFolderDefault + "/" + uuid.NewString() + ".mp4"
	filesToRemove = append(filesToRemove, tmpDubbedFileName, dubbedAudio)

	err = i.mediaHandler.Merge(inputVideoPath, dubbedAudio, tmpDubbedFileName)
	if err != nil {
		return err
	}

	subtitlesPath, err := i.subtitler.GenerateSubtitle(tmpDubbedFileName)
	if err != nil {
		return err
	}
	filesToRemove = append(filesToRemove, subtitlesPath)
	_, err = i.mediaHandler.MergeSubtitle(tmpDubbedFileName, subtitlesPath, outputVideoPath)

	return nil
}

func removeTmpFiles(filesPath *[]string) {
	for _, file := range *filesPath {
		os.Remove(file)
	}
}

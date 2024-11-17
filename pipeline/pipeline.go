package pipeline

import (
	"path/filepath"
	"tradutor-dos-crias/caption"
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

func (i *Pipeline) Run(inputVideoPath, outputVideoPath string) error {
	outputFolderDefault, err := filepath.Abs("pipe/")
	if err != nil {
		return err
	}

	tmpVideoName := outputFolderDefault + "/" + uuid.New().String() + ".mp4"
	tmpAudioName := outputFolderDefault + "/" + uuid.New().String() + ".mp4"

	err = i.mediaHandler.Unmerge(inputVideoPath, tmpVideoName, tmpAudioName)
	if err != nil {
		return err
	}

	text, err := i.transcripter.Transcript(tmpAudioName)
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

	err = i.mediaHandler.Merge(tmpVideoName, dubbedAudio, tmpDubbedFileName)
	if err != nil {
		return err
	}

	subtitlesPath, err := i.subtitler.GenerateSubtitle(tmpDubbedFileName)
	if err != nil {
		return err
	}

	_, err = i.mediaHandler.MergeSubtitle(tmpDubbedFileName, subtitlesPath, outputVideoPath)

	// filesToRemove := []string{tmpAudioName, tmpVideoName, dubbedAudio, subtitlesPath}
	// for _, name := range filesToRemove {
	// 	os.Remove(name)
	// }
	return nil
}

package pipeline

import (
	"os"
	"path/filepath"
	"strings"
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
	filesToRemove := &[]string{tmpAudioName, tmpVideoName}
	defer removeTmpFiles(filesToRemove)
	err = i.mediaHandler.Unmerge(inputVideoPath, tmpVideoName, tmpAudioName)
	if err != nil {
		return err
	}

	text, err := i.transcripter.Transcript(tmpAudioName)
	if err != nil {
		return err
	}

	exts := []string{".srt", ".json", ".txt", ".tsv", ".vtt"}
	for _, ext := range exts {
		name := strings.ReplaceAll(tmpAudioName, ".mp4", ext)
		*filesToRemove = append(*filesToRemove, name)
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
	*filesToRemove = append(*filesToRemove, tmpDubbedFileName, dubbedAudio)

	err = i.mediaHandler.Merge(tmpVideoName, dubbedAudio, tmpDubbedFileName)
	if err != nil {
		return err
	}

	subtitlesPath, err := i.subtitler.GenerateSubtitle(tmpDubbedFileName)
	if err != nil {
		return err
	}
	*filesToRemove = append(*filesToRemove, subtitlesPath)
	_, err = i.mediaHandler.MergeSubtitle(tmpDubbedFileName, subtitlesPath, outputVideoPath)

	return nil
}

func removeTmpFiles(filesPath *[]string) {
	for _, file := range *filesPath {
		os.Remove(file)
	}
}

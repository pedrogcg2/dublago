package media

type MediaHandler interface {
	Unmerge(inputVideoPath string, outputVideoPath string, outputAudioPath string) error
	Merge(inputVideoPath string, inputAudioPath string, outputVideoPath string, strategy MergeStrategy) error
	MergeSubtitle(inputVideoPath, inputSubtitlePath, outputVideoPath string) (string, error)
	ConvertToWav(inputAudioPath string, wavAudioPath string) error
}

type MergeStrategy int

const (
	SpeedUpAudio MergeStrategy = iota
	CutStream
)



package caption

type Subtitler interface {
	GenerateSubtitle(path string) (string, error)
}

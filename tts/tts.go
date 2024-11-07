package tts

type Speeacher interface {
	Speech(message string) (string, error)
}

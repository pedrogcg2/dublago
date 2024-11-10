package translator

type Translator interface {
	Translate(text string) (string, error)
}

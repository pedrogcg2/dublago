package translator

import (
	"context"
	"fmt"
	"log/slog"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

type GoogleTranslator struct{}

func (gt GoogleTranslator) Translate(text string) (string, error) {
	slog.Info("[TRANSLATE] Generating translation now")
	ctx := context.Background()

	client, err := translate.NewClient(ctx)
	if err != nil {
		slog.Error(`[TRANSLATE] Error occured when client has been created: ` + err.Error())
		return "", err
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{text}, language.BrazilianPortuguese, nil)
	if err != nil {
		slog.Error(`[TRANSLATE] Error occured when client called to "Translate": ` + err.Error())
		return "", err
	}
	if len(resp) == 0 {
		slog.Error(`[TRANSLATE] Translate returned empty response to text: ` + text)
		return "", fmt.Errorf("Translate returned empty response to text: %s", text)
	}

	slog.Info(`[TRANSLATE] Translation generated with success`)
	return resp[0].Text, nil
}

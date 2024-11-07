import argparse
from transformers import MarianMTModel, MarianTokenizer


def translate_text(text, target_lang, model_name):
    src_text = f">>{target_lang}<< {text}"
    tokenizer = MarianTokenizer.from_pretrained(model_name)
    model = MarianMTModel.from_pretrained(model_name)
    translated = model.generate(
        **tokenizer([src_text], return_tensors="pt", padding=True))
    return tokenizer.decode(translated[0], skip_special_tokens=True)


def main():
    parser = argparse.ArgumentParser(
        description="Tradução de textos usando MarianMT")
    parser.add_argument("--text", type=str, required=True,
                        help="Texto a ser traduzido")
    parser.add_argument("--target-lang", type=str, default="pt_BR",
                        help="Idioma de destino (ex: 'pt_BR' para português)")
    parser.add_argument(
        "--model", type=str, default="Helsinki-NLP/opus-mt-en-ROMANCE",
        help="Nome do modelo de tradução")
    parser.add_argument(
        "--output", type=str,
        help="Caminho do arquivo de saída para salvar a tradução")

    args = parser.parse_args()

    text_to_translate = args.text.strip('"')
    translated_text = translate_text(
        text_to_translate, args.target_lang, args.model)

    if args.output:
        with open(args.output, "w", encoding="utf-8") as f:
            f.write(translated_text)
        print(f"Tradução salva em {args.output}")
    else:
        print("Texto traduzido:", translated_text)


if __name__ == "__main__":
    main()

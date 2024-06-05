# goDeepL

`goDeepL` is a command-line tool for translating text using the DeepL API. It allows you to easily translate text between different languages, view the current API key, and update the API key.

## Build Instructions

To build the `goDeepL` executable, use the following command:

```sh
  go build -o ./Build/goDeepL.exe
```


This will generate the executable file at `./Build/goDeepL.exe`.

## Usage

### Translating Text

To translate text from one language to another, use the `translate` mode:

```sh
goDeepL -mode translate -from EN -to RU -text "hello world"
```

or you can use the shortcuts:

```sh
goDeepL translate -from EN -to RU -text "hello world"
```

### Viewing Help

To display the usage help, use the `help` flag:

```sh
goDeepL -help
```

### Updating the API Key

To use the `goDeepL` tool, you need a valid DeepL API key. You can update the API key using the `update-key` mode:

```sh
goDeepL -mode update-key
```

This command will prompt you to enter a new API key.

## Command Reference

- `translate`: Translates text from the source language to the target language.
    - `-from`: Specifies the source language (default: `SP`).
    - `-to`: Specifies the target language (default: `EN`).
    - `-text`: Specifies the text to be translated (default: `"Hello!"`).
- `view-key`: Displays the current API key.
- `update-key`: Updates the DeepL API key.
- `delete-key`: Deletes the DeepL API key. (not yet implemented!)
- `help`: Displays the usage help.

## Example Commands

Translate text from English to Russian:

```sh
goDeepL translate -from EN -to RU -text "hello world"
```

Translate text from Spanish to English:

```sh
goDeepL translate -from ES -to EN -text "hola mundo!"
```

View usage help:

```sh
goDeepL -help
```

Update the DeepL API key:
```sh
goDeepL -mode update-key
```

View the DeepL API key:

```sh
goDeepL -mode view-key
```

## Notes

- Ensure you have a valid DeepL API key to use the `translate` functionality.
- The API key is stored in a `key.json` file in the `.goDeepL` directory within your home directory.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [DeepL API](https://www.deepl.com/en/translator) for providing the translation service.

---
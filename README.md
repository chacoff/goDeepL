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


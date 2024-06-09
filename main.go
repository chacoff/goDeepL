/*
 * File:    main.go
 * Date:    June 05, 2024
 * Author:  J.
 * Email:   jaime.gomez@usach.cl
 * Project: goDeepL
 * Description:
 *   Uses the REST API from DeepL Translator to translate text and phrases via CLI
 *
 *
 * - Basic Build:
 *    go build -o ./Build/goDeepL.exe
 *
 */

package main

import (
	"flag"
	"fmt"
	"os"
)

type Args struct {
	Help bool   // help mode
	Mode string // functionality
	From string // original language
	To   string // target language
	Text string // text to translate
}

type Config struct {
	APIKey string `json:"apikey"` // api key from deepL
}

var (
	config Config
	args   Args
	text   []string
)

// init
func init() {
	readJson()
}

// main
func main() {

	flag.BoolVar(&args.Help, "help", false, "usage help")
	flag.StringVar(&args.Mode, "mode", "", "activate the translation mode")
	flag.StringVar(&args.From, "from", "", "original language")
	flag.StringVar(&args.To, "to", "", "target language")
	flag.StringVar(&args.Text, "text", "", "input text")
	flag.Usage = usage
	flag.Parse()

	if args.Help {
		flag.Usage()
		return
	}

	remainingArgs := flag.Args()

	// fmt.Println(remainingArgs)

	if len(remainingArgs) >= 3 {
		args.Mode = "translate"
		args.From = remainingArgs[0]
		args.To = remainingArgs[1]
		args.Text = remainingArgs[2]
	}

	if len(remainingArgs) == 2 && args.Mode == "" {
		args.Mode = "translate"
		args.To = remainingArgs[0]
		args.Text = remainingArgs[1]
	}

	if len(remainingArgs) == 1 {
		args.Mode = remainingArgs[0]
	}

	if len(remainingArgs) == 0 {
		flag.Usage()
		return
	}

	cliLogic()
}

// cliLogics handles the arguments to properly decide either to call the api to translate or use another functionality
func cliLogic() {

	switch args.Mode {

	case "help":
		flag.Usage()

	case "view-key":
		fmt.Println(getKey())

	case "update-key":
		updateKey()

	case "delete-key":
		fmt.Println("Delete key functionality is not yet implemented.")

	case "translate":
		if args.To == "" || args.Text == "" {
			fmt.Println("Error: -to, and -text arguments are required for translate mode. -from is recommended.")
			fmt.Print("godeepl -mode translate -from EN -to RU -text \"hello world\" ")
			return
		}

		if args.From == "" {
			fmt.Println(">> deepL API recommends to include the source language whenever is possible.")
		}

		translatedText, err := getTranslation()
		if err != nil {
			return
		}

		fmt.Println(translatedText)

	default:
		fmt.Println("Error: unknown mode. Please use goDeepL -help to see the usage.")
		// flag.Usage()
		return
	}
}

// usage displays information and examples on how to use the tool
func usage() {

	fmt.Fprintf(os.Stderr, "\nUsage: goDeepL [options] [arguments]\n")
	fmt.Fprintln(os.Stderr, "\nOptions:")

	flag.PrintDefaults()

	fmt.Fprintln(os.Stderr, "\nModes:")
	fmt.Fprintln(os.Stderr, "  translate          Triggers the functionalities to translate")
	fmt.Fprintln(os.Stderr, "  view-key           View the API key")
	fmt.Fprintln(os.Stderr, "  update-key         Update the API key")
	fmt.Fprintln(os.Stderr, "  delete-key         Delete the API key")
	fmt.Fprintf(os.Stderr, "\nExample: goDeepL -from EN -to RU -text hello world\n")
}

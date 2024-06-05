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
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type Args struct {
	Help bool   // help mode
	From string // original language
	To   string // target language
	Text string // text to translate
	Mode string
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

	commandMap := map[string]string{
		"translate": "-mode=translate",
		"view":      "-mode=view-key",
		"update":    "-mode=update-key",
		"delete":    "-mode=delete-key",
		"help":      "-help",
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: <command> [arguments]\nExample: goDeepL -help")
		os.Exit(1)
	}

	command := os.Args[1]
	remainingArgs := os.Args[2:]

	if fullCommand, exists := commandMap[command]; exists {
		// Split fullCommand to extract key-value pairs
		fullCommandParts := strings.Split(fullCommand, " ")
		os.Args = append([]string{os.Args[0]}, append(fullCommandParts, remainingArgs...)...)
	}

	flag.BoolVar(&args.Help, "help", false, "usage help")
	flag.StringVar(&args.Mode, "mode", "key", "show the used key")
	flag.StringVar(&args.From, "from", "SP", "original language")
	flag.StringVar(&args.To, "to", "EN", "target language")
	flag.StringVar(&args.Text, "text", "Hello!", "input text")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nUsage: goDeepL [options] [arguments]\n")
		fmt.Fprintln(os.Stderr, "\nOptions:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nModes:")
		fmt.Fprintln(os.Stderr, "  view-key           View the API key")
		fmt.Fprintln(os.Stderr, "  update-key         Update the API key")
		fmt.Fprintln(os.Stderr, "  delete-key         Delete the API key")
		fmt.Fprintln(os.Stderr, "  translate          Translate text (requires -from, -to, and -text)")
		fmt.Fprintf(os.Stderr, "\nExample: goDeepL -mode translate -from EN -to RU -text hello world\n")
	}

	flag.Parse()

	// Display help if the -help flag is set
	if args.Help {
		flag.Usage()
		return
	}

	// Handle the different modes
	switch args.Mode {

	case "view-key":
		fmt.Println(config.APIKey)

	case "update-key":
		updateKey()

	case "delete-key":
		fmt.Println("Delete key functionality is not yet implemented.")

	case "translate":
		if args.From == "" || args.To == "" || args.Text == "" {
			fmt.Println("Error: -from, -to, and -text arguments are required for translate mode.")
			flag.Usage()
			os.Exit(1)
		}

		translatedText, err := getTranslation()
		if err != nil {
			fmt.Println("Output error:", err)
			return
		}

		fmt.Println(translatedText)

	default:
		fmt.Println("Error: unknown mode. Please use goDeepL -help to see the usage.")
		// flag.Usage()
		os.Exit(1)
	}
}

// getTranslations make the http request to deepL. you need to have your own key
func getTranslation() (string, error) {

	text = append(text, args.Text)

	// payload
	payload := map[string]interface{}{
		"text":        text,
		"target_lang": args.To,
		"source_lang": args.From,
		// "glossary_id": "yourGlossaryId",
	}

	// fmt.Printf("payload: %s \n", payload)

	// payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling payload:", err)
		return "", err
	}

	// HTTP request
	req, err := http.NewRequest("POST", "https://api-free.deepl.com/v2/translate", bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	// HTTP headers
	auth := fmt.Sprintf("DeepL-Auth-Key %s", config.APIKey) // "DeepL-Auth-Key [yourAuthKey]"
	req.Header.Set("Authorization", auth)
	req.Header.Set("User-Agent", "goDeepL/0.0.1")
	req.Header.Set("Content-Type", "application/json")

	// fmt.Println("\n", req)

	// request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", err
	}

	translatedText, err := getResponseText(body)
	if err != nil {
		fmt.Println("No translation found", err)
		return "", err
	}

	return translatedText, nil
}

// getResponseText handles the response from deepL to extract only the translated text
func getResponseText(body []byte) (string, error) {

	var result map[string]interface{}

	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return "", nil
	}

	if translations, ok := result["translations"].([]interface{}); ok && len(translations) > 0 {

		if firstTranslation, ok := translations[0].(map[string]interface{}); ok {

			if translatedText, ok := firstTranslation["text"].(string); ok {
				return translatedText, nil
			}
		}
	}

	return "", errors.New("no translation found")
}

// getCurrentUser to create the folder .goDeepL and store there the api key
func getCurrentUser() *user.User {
	usr, err := user.Current()

	if err != nil {
		fmt.Println("Error getting current user:", err)
		os.Exit(1)
	}

	return usr
}

// createGoDeepLFolder
func createGoDeepLFolder() string {

	usr := getCurrentUser()

	dirPath := filepath.Join(usr.HomeDir, ".goDeepL")

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory .goDeepL under user folder:", err)
			os.Exit(1)
		}
	}

	return dirPath
}

// readJson reads the api key of the json file
func readJson() {

	dirPath := createGoDeepLFolder()

	keyPath := filepath.Join(dirPath, "key.json")

	if _, err := os.Stat(keyPath); os.IsNotExist(err) {

		content := map[string]string{"apikey": "apiapiapi"}
		file, err := os.Create(keyPath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			os.Exit(1)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		err = encoder.Encode(content)
		if err != nil {
			fmt.Println("Error writing JSON to file:", err)
			os.Exit(1)
		}
		fmt.Println("First use detected:\n  Do not forget to change the api key with your own: goDeepL -mode update-key")
	} else {
		file, err := os.Open(keyPath)
		if err != nil {
			fmt.Println("Error opening file:", err)
			os.Exit(1)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		err = decoder.Decode(&config)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			os.Exit(1)
		}
	}
}

// updateKey writes the new key
func updateKey() {

	dirPath := createGoDeepLFolder()

	keyPath := filepath.Join(dirPath, "key.json")

	fmt.Print("Enter the new API key for DeepL Translator: ")

	reader := bufio.NewReader(os.Stdin)
	newAPIKey, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}
	newAPIKey = strings.TrimSpace(newAPIKey) // Remove the newline character

	content := map[string]string{"apikey": newAPIKey}
	file, err := os.Create(keyPath)
	if err != nil {
		fmt.Println("Error creating or opening file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	err = encoder.Encode(content)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return
	}

	fmt.Println("api Key updated successfully:", keyPath)
}

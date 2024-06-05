package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"io"
	"net/http"
)

type Args struct{
	Help 		bool		// help mode
	From		string		// original language
	To			string		// target language
	Text		string		// text to translate
	Mode 		string
	Verbose		bool		// verbose mode
}

type Config struct {
	APIKey string `json:"apikey"`	// api key from deepL
}

var (
	config Config
	args Args
	text []string
)

//init 
func init(){
	readJson()
}

//main
func main(){
	flag.BoolVar(&args.Help, "help", false, "usage help")
	flag.StringVar(&args.Mode, "mode", "key", "show the used key")
	flag.StringVar(&args.From, "from", "SP", "original language")
	flag.StringVar(&args.To, "to", "EN", "target language")
	flag.StringVar(&args.Text, "text", "Hello!", "input text")
	flag.BoolVar(&args.Verbose, "v", false, "enable verbose mode")


	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nUsage: goDeepL [options] [arguments]\n")
		fmt.Fprintln(os.Stderr, "\nOptions:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nModes:")
		fmt.Fprintln(os.Stderr, "  view-key           View the API key")
		fmt.Fprintln(os.Stderr, "  update-key         Update the API key")
		fmt.Fprintln(os.Stderr, "  delete-key         Delete the API key")
		fmt.Fprintln(os.Stderr, "  translate          Translate text (requires -from, -to, and -text)")
		fmt.Fprintf(os.Stderr, "\nExample: goDeepL -mode translate -From EN -TO RU -text hello world\n")
	}

	// Parse the flags
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
		fmt.Println("Update key functionality is not yet implemented.")

	case "delete-key":
		fmt.Println("Delete key functionality is not yet implemented.")

	case "translate":
		// Validate mandatory arguments for translate mode
		if args.From == "" || args.To == "" || args.Text == "" {
			fmt.Println("Error: -from, -to, and -text arguments are required for translate mode.")
			flag.Usage()
			os.Exit(1)
		}

		// Implement your translation logic here
		fmt.Printf("Translating from %s to %s\n", args.From, args.To)

		translatedText, err := getTranslation()
		if err != nil {
			fmt.Println("Output error:", err)
			return
		}
		fmt.Println(translatedText)

	default:
		fmt.Println("Error: unknown mode. Please use -help to see the usage.")
		flag.Usage()
		os.Exit(1)
	}
		// fmt.Printf("Translating from %s to %s\n", args.From, args.To)

		// translatedText, err := getTranslation()
		// if err !=nil{
		// 	fmt.Println("Output error:", err)
		// 	return
		// }
		// fmt.Println(translatedText)

}

//getTranslations make the http request to deepL. you need to have your own key
func getTranslation() (string, error){

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
	auth := fmt.Sprintf("DeepL-Auth-Key %s", config.APIKey)  // "DeepL-Auth-Key [yourAuthKey]"
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

	translatedText := getReponseText(body)

	return translatedText, nil
}

//getResponseText handles the response from deepL to extract only the translated text
func getReponseText(body []byte) string{

	var result map[string]interface{}
	var translatedText string

	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return ""
	}

	if translations, ok := result["translations"].([]interface{}); ok && len(translations) > 0 {

		if firstTranslation, ok := translations[0].(map[string]interface{}); ok {

			if translatedText, ok := firstTranslation["text"].(string); ok {
				fmt.Println("Translated Text:", translatedText)
			}
		}
	}

	return translatedText
}

//readJson reads the api key of the json file
func readJson(){
	file, err := os.Open("key.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
}

//writeJson writes the new key
func writeJson(){
	file, err := os.Create("key.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
}

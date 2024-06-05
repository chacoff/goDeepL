package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
)

type Args struct{
	From		string		// original language
	To			string		// target language
	Text		string		// text to translate
	ApiKey		string		// api key from deepL
	Verbose		bool		// verbose mode
}

type TranslationResponse struct {
	Translations []struct {
		DetectedSourceLanguage string `json:"detected_source_language"`
		Text                   string `json:"text"`
	} `json:"translations"`
}

var (
	args Args
	text []string
)

//init 
func init(){
	args.ApiKey = "8ba5bf2b-6276-4cb5-bbf8-ade500f70285:fx"
}

//main
func main(){
	flag.StringVar(&args.From, "from", "SP", "original language")
	flag.StringVar(&args.To, "to", "EN", "target language")
	flag.StringVar(&args.Text, "text", "Hello!", "input text")
	flag.BoolVar(&args.Verbose, "v", false, "enable verbose mode")

	flag.Parse()

	fmt.Printf("Translating from %s to %s\n", args.From, args.To)
	
	translatedText, err := getTranslation()
	if err !=nil{
		fmt.Println("Output error:", err)
		return
	}

	fmt.Println(translatedText)
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
	auth := fmt.Sprintf("DeepL-Auth-Key %s", args.ApiKey)  // "DeepL-Auth-Key [yourAuthKey]"
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

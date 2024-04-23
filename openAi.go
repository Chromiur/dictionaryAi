package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
)

const (
	apiEndpoint      = "https://api.openai.com/v1/chat/completions"
	russianMessage   = "Сгенерируй пожалуйста сообщение на русском с этими словами: %+q"
	engMessage       = "Hi, can you please generate the sentence with thees words: %+q"
	translateMessage = "Hi, can you please translate this sentence into English: %s"
)

func GenerateRussianMessageRequest(words []string) (string, error) {
	message := fmt.Sprintf(russianMessage, words)
	return GenerateRequest(message)
}

func TranslateMessageRequest(messageOnRussian string) (string, error) {
	message := fmt.Sprintf(translateMessage, messageOnRussian)
	return GenerateRequest(message)
}

func GenerateRequest(message string) (string, error) {
	apiKey := openApiKey
	client := resty.New()
	//engMessage := fmt.Sprintf("Hi, can you please generate the sentence with thees words: %+q", words)

	response, err := client.R().
		SetAuthToken(apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model":      "gpt-3.5-turbo",
			"messages":   []interface{}{map[string]interface{}{"role": "system", "content": message}},
			"max_tokens": 50,
		}).
		Post(apiEndpoint)

	if err != nil {
		log.Fatalf("Failed to send the request: %v", err)
	}
	// Read the response and parse the JSON
	body := response.Body()

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error while decoding JSON response:", err)
		return "", err
	}

	if val, ok := data["error"]; ok {
		log.Fatalf("Error : %v", val.(map[string]interface{})["message"])

	}

	// Extract the content from the JSON response
	content := data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)

	return content, nil
}

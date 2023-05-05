package notification

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"log"
	"github.com/joho/godotenv"
	"os"
)

func SendSlack(msg string) {
	err := godotenv.Load(".env")
	WEBHOOK := os.Getenv("WEBHOOK")
	payload := fmt.Sprintf(`{
		"text": "Custom controller notification",
		"blocks": [
			{
				"type": "section",
				"text": {
					"type": "mrkdwn",
					"text": "New Deployment %s has created."
				}
			},
			{
				"type": "section",
				"block_id": "section567",
				"text": {
					"type": "mrkdwn",
					"text": "<https://kubernetes.io|Kubernetes> \n :heart: \n Slack notification from Kubernetes."
				},
				"accessory": {
					"type": "image",
					"image_url": "https://i.imgur.com/MRyQova.png",
					"alt_text": "Server Monk"
				}
			},
			{
				"type": "section",
				"block_id": "section789",
				"fields": [
					{
						"type": "mrkdwn",
						"text": "*Yours,*\nServer Monk"
					}
				]
			}
		]
	}`, msg)
	var jsonData = []byte(payload)
	request, err := http.NewRequest("POST", WEBHOOK, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Panic(err)
		}
	}(response.Body)

	log.Println(fmt.Sprintf("Response Status: %s", response.Status))
	log.Println(fmt.Sprintf("Response Headers: %s", response.Header))
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	log.Println(fmt.Sprintf("Response Body: %s", string(body)))
}

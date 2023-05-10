package usecase

import (
	"context"
	speechkit "fsm/internal/api"
	"fsm/internal/models"
	"fsm/pkg/repository"
	"github.com/sashabaranov/go-openai"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

func GenerateTale(text string, body models.User) string {
	var result string

	c := openai.NewClient(os.Getenv("GPT"))
	resp, err := c.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: text,
				},
			},
		},
	)

	if err != nil {
		log.Println(err)
	}

	result = resp.Choices[0].Message.Content
	if body.Sounder == "Yandex" {
		repository.UpdateCounter(body.UserId, 1)
		if body.Counter+1 > 15 {
			return "Простите, но Вы превысили количество запросов к Yandex SpeechKit."
		}
		yandexGenerate(result)
		return "."

	}
	return result
}

func yandexGenerate(text string) string {
	API_KEY := os.Getenv("YANDEX")
	client := &http.Client{Timeout: 60 * time.Second}
	apiParams := speechkit.APIParams{APIKey: API_KEY, Client: client}

	// define folder for audio
	currentDir, _ := os.Getwd()
	pathToFiles := path.Join(currentDir, "temp-folder")

	speechParams := speechkit.SpeechParams{
		Voice:       "female",
		Speed:       1.0,
		PathToFiles: pathToFiles,
	}

	client1 := speechkit.NewSpeechKitClient(apiParams, speechParams)

	err := client1.CreateAudio(text)
	if err != nil {
		log.Println(err)
	}
	path, _ := os.Getwd()
	err = os.Remove(path + "/temp-folder/output.txt")
	if err != nil {
		log.Println(err)
	}
	return ""
}

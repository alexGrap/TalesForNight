package api

/*This part of code was */

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

const URL = "https://tts.api.cloud.yandex.net/speech/v1/tts:synthesize"

var (
	speechSpeed    = 1.0
	speechLanguage = "ru-RU"
	speechFormat   = "oggopus"
	speechEmotion  = "neutral"
	textMaxLen     = 2000
	output         = "output.txt"
)

type SpeechKitClient struct { //nolint
	APIParams
	SpeechParams
}

type APIParams struct {
	Client *http.Client
	APIKey string
}

type SpeechParams struct {
	Emotion     string
	Voice       string
	Speed       float64
	PathToFiles string
}

func NewSpeechKitClient(apiParams APIParams, speechParams SpeechParams) *SpeechKitClient {
	return &SpeechKitClient{
		APIParams:    apiParams,
		SpeechParams: speechParams,
	}
}

func (c *SpeechKitClient) CreateAudio(text string) error {
	output, err := c.createFile()
	if err != nil {
		return errors.Wrap(err, "error: while creating output.txt file")
	}
	defer output.Close()

	if err != nil {
		return errors.Wrap(err, "error: occurred while splitting the text")
	}

	var stringToFile string
	stringToFile += fmt.Sprintf("file '%s'", "file.ogg")
	_, err = output.WriteString(stringToFile)
	if err != nil {
		return errors.Wrap(err, "error: occurred while writing to file")
	}

	err = c.doRequest(text, "file.ogg")
	if err != nil {
		return err
	}
	return nil
}

func (c *SpeechKitClient) createFile() (*os.File, error) {
	output := path.Join(c.PathToFiles, output)
	var _, err = os.Stat(output)
	if os.IsNotExist(err) {
		var file, err = os.Create(output)
		if err != nil {
			return nil, err
		}
		return file, err
	}
	return nil, errors.New("error: file already exists")
}

func (c *SpeechKitClient) generateURL(text string) string {
	if c.SpeechParams.Speed == 0.0 {
		c.SpeechParams.Speed = speechSpeed
	}

	if c.SpeechParams.Voice == "female" {
		c.SpeechParams.Voice = "alena"
	} else if c.SpeechParams.Voice == "male" {
		c.SpeechParams.Voice = "filipp"
	}

	if c.SpeechParams.Emotion == "" {
		c.SpeechParams.Emotion = speechEmotion
	}

	v := url.Values{}
	v.Add("text", text)
	v.Add("speed", fmt.Sprintf("%.2f", c.SpeechParams.Speed))
	v.Add("emotion", c.SpeechParams.Emotion)
	v.Add("voice", c.SpeechParams.Voice)
	v.Add("lang", speechLanguage)
	v.Add("format", speechFormat)
	return v.Encode()
}

func (c *SpeechKitClient) doRequest(text, fileName string) error {
	body := strings.NewReader(c.generateURL(text))
	req, err := http.NewRequest(http.MethodPost, URL, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Add("Authorization", fmt.Sprintf("Api-Key %s", c.APIParams.APIKey))

	response, err := c.APIParams.Client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("error: api occurred with status: %v", response.StatusCode))
	}

	fullFilePath := path.Join(c.PathToFiles, fileName)
	outputFile, err := os.Create(fullFilePath)
	if err != nil {
		return errors.Wrap(err, "error: occurred while creating audio file")
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, response.Body)
	if err != nil {
		return errors.Wrap(err, "error: occurred while copying response to file")
	}
	return nil
}

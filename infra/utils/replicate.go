package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/yusifubot/infra/vars"
)

const GPT41NanoAPI = "https://api.replicate.com/v1/models/openai/gpt-4.1-nano/predictions"

type replicateReq struct {
	Input replicateInput `json:"input"`
}

type replicateResp struct {
	ID          string         `json:"id"`
	Model       string         `json:"model"`
	Version     string         `json:"version"`
	Input       replicateInput `json:"input"`
	Logs        string         `json:"logs"`
	Output      []string       `json:"output"`
	DataRemoved bool           `json:"data_removed"`
	Error       *string        `json:"error"`
	Status      string         `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
}
type replicateInput struct {
	FrequencyPenalty    int     `json:"frequency_penalty"`
	MaxCompletionTokens int     `json:"max_completion_tokens"`
	PresencePenalty     int     `json:"presence_penalty"`
	Prompt              string  `json:"prompt"`
	SystemPrompt        string  `json:"system_prompt"`
	Temperature         float64 `json:"temperature"`
	TopP                int     `json:"top_p"`
}

func CallGPT41Nano(prompt, input string, temp float64, maxTokens int, timeout time.Duration) (string, error) {
	reqBody := replicateReq{Input: replicateInput{
		Prompt:              input,
		SystemPrompt:        prompt,
		Temperature:         temp,
		MaxCompletionTokens: maxTokens,
		TopP:                1,
	}}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	logrus.Debugf("replicate req: %s", string(reqBytes))

	req, err := http.NewRequest(http.MethodPost, GPT41NanoAPI, bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+vars.ReplicateToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "wait")

	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respRaw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	logrus.Debugf("replicate resp: %s", string(respRaw))

	respBody := replicateResp{}
	err = json.Unmarshal(respRaw, &respBody)
	if err != nil {
		return "", err
	}
	if respBody.Error != nil {
		return "", errors.New(*respBody.Error)
	}
	return strings.Join(respBody.Output, ""), nil
}

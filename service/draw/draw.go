package draw

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/yusifubot/infra/utils"
	"github.com/zjyl1994/yusifubot/infra/vars"
)

const (
	IMAGE_MODEL_IDENTIFIER  = "prunaai/z-image-turbo"
	PROMPT_MODEL_IDENTIFIER = "openai/gpt-5-nano"
	replicateAPIBaseURL     = "https://api.replicate.com/v1/models/"
)

type replicatePredictionRequest struct {
	Input map[string]interface{} `json:"input"`
}

type replicatePredictionResponse struct {
	Output interface{} `json:"output"`
	Error  interface{} `json:"error"`
	Status string      `json:"status"`
}

type replicateProblemResponse struct {
	Detail interface{} `json:"detail"`
	Status int         `json:"status"`
	Title  string      `json:"title"`
}

func DrawImageHandler(msg *models.Message) error {
	if msg.From.ID != vars.AdminUserId && !checkChatIdAllowed(msg.Chat.ID) {
		return utils.ReplyTextToTelegram(msg, "该聊天未开启绘图能力", false)
	}
	commandArgs := utils.ParseCommandArguments(msg.Text)
	prompt := strings.TrimSpace(strings.Join(commandArgs, " "))
	if prompt == "" {
		return utils.ReplyTextToTelegram(msg, "请输入提示词", false)
	}

	if msg.From.ID != vars.AdminUserId {
		key := "draw:" + strconv.FormatInt(msg.From.ID, 10)
		if !vars.ReplicateCooldown.CheckAndSetCooldown(key, 1*time.Minute) {
			remaining := vars.ReplicateCooldown.RemainingTime(key)
			return utils.ReplyTextToTelegram(msg, fmt.Sprintf("你的绘图能力冷却中，请等待 %d 秒", int(remaining.Seconds())), false)
		}
	}

	go func(ctx context.Context, m *models.Message) {
		preparedPrompt, err := preparePrompt(context.Background(), prompt)
		if err != nil {
			utils.ReplyTextToTelegram(msg, err.Error(), false)
			return
		}
		url, err := drawWithReplicate(ctx, preparedPrompt)
		if err != nil {
			utils.ReplyTextToTelegram(msg, err.Error(), false)
			return
		}

		params := &bot.SendPhotoParams{
			ChatID: msg.Chat.ID,
			ReplyParameters: &models.ReplyParameters{
				MessageID: msg.ID,
			},
			Photo:               &models.InputFileString{Data: url},
			DisableNotification: true,
		}

		_, err = vars.BotInstance.SendPhoto(ctx, params)
		if err != nil {
			utils.ReplyTextToTelegram(msg, err.Error(), false)
			return
		}
	}(context.Background(), msg)
	return nil
}

func drawWithReplicate(ctx context.Context, prompt string) (string, error) {
	output, err := runReplicateModel(ctx, IMAGE_MODEL_IDENTIFIER, map[string]interface{}{
		"prompt": prompt,
	})
	if err != nil {
		logrus.WithError(err).Errorln("run image replicate model failed")
		return "", utils.NewBizErr("run replicate model failed")
	}
	logrus.Debugf("replicate output: %s", utils.MarshalToJsonNoError(output))
	outputURL, err := extractImageURL(output)
	if err != nil {
		return "", err
	}
	return outputURL, nil
}

func extractImageURL(output interface{}) (string, error) {
	if outputURL, ok := output.(string); ok {
		outputURL = strings.TrimSpace(outputURL)
		if outputURL == "" {
			return "", utils.NewBizErr("image output is empty")
		}
		return outputURL, nil
	}

	if outputs, ok := output.([]interface{}); ok {
		for _, item := range outputs {
			if outputURL, ok := item.(string); ok {
				outputURL = strings.TrimSpace(outputURL)
				if outputURL != "" {
					return outputURL, nil
				}
			}
		}
		return "", utils.NewBizErr("image output is empty")
	}

	return "", utils.NewBizErr(fmt.Sprintf("unsupported image output type: %s", reflect.TypeOf(output)))
}

func preparePrompt(ctx context.Context, prompt string) (string, error) {
	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		return "", utils.NewBizErr("prompt is empty")
	}
	// 标记为raw的prompt，直接返回不做处理
	if cut, found := strings.CutPrefix(prompt, "raw:"); found {
		return cut, nil
	}
	// 使用gpt-5-nano进行拓展
	output, err := runReplicateModel(ctx, PROMPT_MODEL_IDENTIFIER, map[string]interface{}{
		"system_prompt":         "你是一个专业摄影师，解读输入的文字，转化为包含镜头、主题、环境、灯光、风格、摄影参数的明确绘图指令并输出，无需输出更多其他内容。",
		"prompt":                prompt,
		"max_completion_tokens": 1024,
	})
	if err != nil {
		logrus.WithError(err).Errorln("run prompt replicate model failed")
		return "", utils.NewBizErr("run replicate model failed")
	}
	logrus.Debugf("replicate output: %s", utils.MarshalToJsonNoError(output))

	var outputText string
	if str, ok := output.(string); ok {
		outputText = str
	} else if tokens, ok := output.([]interface{}); ok {
		sb := strings.Builder{}
		for _, token := range tokens {
			if s, ok := token.(string); ok {
				sb.WriteString(s)
			}
		}
		outputText = sb.String()
	} else {
		return "", utils.NewBizErr("output format error")
	}
	return outputText, nil
}

func runReplicateModel(ctx context.Context, identifier string, input map[string]interface{}) (interface{}, error) {
	requestBody, err := json.Marshal(replicatePredictionRequest{Input: input})
	if err != nil {
		return nil, fmt.Errorf("marshal replicate request failed: %w", err)
	}

	for attempt := 0; attempt < 2; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, replicateAPIBaseURL+identifier+"/predictions", bytes.NewReader(requestBody))
		if err != nil {
			return nil, fmt.Errorf("create replicate request failed: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+vars.ReplicateToken)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Prefer", "wait")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("execute replicate request failed: %w", err)
		}

		respBody, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return nil, fmt.Errorf("read replicate response failed: %w", readErr)
		}

		if resp.StatusCode == http.StatusTooManyRequests && attempt == 0 {
			waitDuration := parseReplicateRetryAfter(resp.Header)
			if waitDuration > 0 {
				select {
				case <-time.After(waitDuration):
					continue
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			}
		}

		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
			return nil, formatReplicateAPIError(resp.StatusCode, respBody)
		}

		var predictionResp replicatePredictionResponse
		if err := json.Unmarshal(respBody, &predictionResp); err != nil {
			return nil, fmt.Errorf("unmarshal replicate response failed: %w", err)
		}
		if predictionResp.Error != nil {
			return nil, fmt.Errorf("replicate model error: %v", predictionResp.Error)
		}
		if predictionResp.Status != "" && predictionResp.Status != "succeeded" {
			return nil, fmt.Errorf("replicate prediction status: %s", predictionResp.Status)
		}
		return predictionResp.Output, nil
	}

	return nil, fmt.Errorf("replicate request exhausted retries")
}

func formatReplicateAPIError(statusCode int, respBody []byte) error {
	var problem replicateProblemResponse
	if err := json.Unmarshal(respBody, &problem); err == nil {
		if detail := strings.TrimSpace(fmt.Sprint(problem.Detail)); detail != "" && detail != "<nil>" {
			return fmt.Errorf("replicate api status %d: %s", statusCode, detail)
		}
		if title := strings.TrimSpace(problem.Title); title != "" {
			return fmt.Errorf("replicate api status %d: %s", statusCode, title)
		}
	}
	return fmt.Errorf("replicate api status %d: %s", statusCode, strings.TrimSpace(string(respBody)))
}

func parseReplicateRetryAfter(header http.Header) time.Duration {
	for _, key := range []string{"Retry-After", "ratelimit-reset"} {
		value := strings.TrimSpace(header.Get(key))
		if value == "" {
			continue
		}
		seconds, err := strconv.Atoi(value)
		if err == nil && seconds > 0 {
			return time.Duration(seconds) * time.Second
		}
	}
	return 10 * time.Second
}

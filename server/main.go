package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
)

func main() {
	r := gin.Default()

	r.POST("/data", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("FormFile error: %v", err))
			return
		}
		err = c.SaveUploadedFile(file, "input_1.wav")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("SaveUploadedFile error: %v", err))
			return
		}

		client := openai.NewClient("xxx")
		text, err := client.CreateTranscription(
			context.Background(),
			openai.AudioRequest{
				Model:    openai.Whisper1,
				FilePath: "input_1.wav",
			},
		)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("CreateTranscription error: %v", err))
			return
		}
		log.Println(text.Text)

		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: "你现在是温柔知心的女性，请在30个字以内完成这个回答，问题内容如下：" + text.Text, // 使用获取到的问题内容
					},
				},
			},
		)

		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("ChatCompletion error: %v", err))
			return
		}

		res_file, err := client.CreateSpeech(
			context.Background(),
			openai.CreateSpeechRequest{
				Model: openai.TTSModel1,
				Input: resp.Choices[0].Message.Content,
				Voice: openai.VoiceNova,
			},
		)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("CreateSpeech error: %v", err))
			return
		}
		buf, err := io.ReadAll(res_file)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("CreateSpeech error: %v", err))
			return
		}
		err = os.WriteFile("output_1.mp3", buf, 0644)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("CreateSpeech error: %v", err))
			return
		}

		cmdArgs := []string{"-i", "output_1.mp3", "output_1.wav", "-y"} // 覆盖输出文件
		cmd := exec.Command("ffmpeg", cmdArgs...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			log.Fatal("执行FFmpeg命令出错:", err)
		}

		c.File("output_1.wav")
	})

	r.Run(":8080")
}

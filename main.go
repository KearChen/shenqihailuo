package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const OpenAIAPIURL = "https://open.bigmodel.cn/api/paas/v4/chat/completions"

func main() {
	apiKey := ("d3d31f1d454e4cdea4fb829de406e147.rjBZyv4QRvO5oeu6")
	if apiKey == "" {
		log.Fatal("请设置环境变量 OPENAI_API_KEY")
	}

	// 初始化Gin路由
	router := gin.Default()

	// 静态文件服务，指定public目录为静态文件目录
	router.Static("/", "./public")

	// 定义API路由，用于获取答案
	router.POST("/api/getAnswer", func(c *gin.Context) {
		var request struct {
			Question string `json:"question"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
			return
		}

		answer, err := getChatGPTAnswer(apiKey, request.Question)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法生成回答，请稍后重试"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"answer": answer})
	})

	// 启动服务器
	router.Run(":2000")
}

// getChatGPTAnswer 调用OpenAI API生成回答
func getChatGPTAnswer(apiKey, question string) (string, error) {
	prompt := "你是一只神秘的海螺，智慧深邃、充满神秘。请以神秘而简洁的风格回答问题："

	requestBody, err := json.Marshal(map[string]interface{}{
		"model": "GLM-4-Air",
		"messages": []map[string]string{
			{"role": "user", "content": prompt + question},
		},
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", OpenAIAPIURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var responseData struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &responseData); err != nil {
		return "", err
	}

	if len(responseData.Choices) > 0 {
		return responseData.Choices[0].Message.Content, nil
	}
	return "无法生成回答", nil
}

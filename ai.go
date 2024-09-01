package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

const (
	API = "https://www.blackbox.ai/api/chat"
)

var (
	rxpCleanReply = regexp.MustCompile(`\$\@\$(.*?)\$\@\$(.*?)`)
)	

type Message struct {
	Author  string
	Content string
}

type AgentMode struct {
	Id   string `json:"id"`
	Mode bool   `json:"mode"`
}

type RequestBody struct {
	CodeModelMode     bool             `json:"codeModelMode"`
	ClickedAnswer3    bool             `json:"clickedAnswer3"`
	AgentMode         AgentMode        `json:"agentMode"`
	TrendingAgentMode AgentMode        `json:"trendingAgentMode"`
	Messages          []RequestMessage `json:"messages"`
}

type RequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatBot struct {
	History []Message
	Client  *http.Client
}

func (bot *ChatBot) Send(body RequestBody, Username string) string {

	byteData, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
		return "```\n" + err.Error() + "\n```"
	}

	bot.History = append(bot.History, Message{
		Author:  Username,
		Content: body.Messages[0].Content,
	})

	req, err := http.NewRequest("POST", API, bytes.NewBuffer(byteData))
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := bot.Client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "```\n" + err.Error() + "\n```"
	}

	defer resp.Body.Close()

	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}

	reply := rxpCleanReply.ReplaceAllString(string(bodyByte), "")

	bot.History = append(bot.History, Message{
		Author:  "Assistant",
		Content: reply,
	})

	return reply
}

func (bot *ChatBot) GetHistory() string {
	var history string
	for _,msg := range bot.History {
		history += msg.Author + " messaged \"" + msg.Content + "\";\n"
	}
	return history
}

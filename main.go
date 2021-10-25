package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"net/http"
	"log"
	"strconv"
	"github.com/thoas/go-funk"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/joho/godotenv"
)

func main() {
	mainRoutin()
}

func mainRoutin(){
	err := godotenv.Load(".env")
	
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	} 

	url := "https://api.todoist.com/rest/v1/tasks"
	authHeaderName := "Authorization"
	authHeaderValue := "Bearer " + os.Getenv("TODOIST_TOKEN")

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(authHeaderName, authHeaderValue)


	client := new(http.Client)
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error Request:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
			fmt.Println("Error Response:", resp.Status)
			return
	}

	body, _ := io.ReadAll(resp.Body)

	var articles []Task
	if err := json.Unmarshal(body, &articles); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
	}
	bot, err := linebot.New(
		os.Getenv("LINE_BOT_CHANNEL_SECRET"),
		os.Getenv("LINE_BOT_CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	inboxTasks := funk.Filter(articles, func(x Task) bool {
    return x.ProjectID == 1528596419
	})
	inboxTasksNum := len(inboxTasks.([]Task))

	text := strconv.Itoa(inboxTasksNum) + "個、タスクが残ってるよ！"
	message := linebot.NewTextMessage(text)

	if _, err := bot.BroadcastMessage(message).Do(); err != nil {
		log.Fatal(err)	
	}
}

type Task struct {
	ID           int64  `json:"id"`
	ProjectID    int64  `json:"project_id"`
	SectionID    int    `json:"section_id"`
	ParentID     int64  `json:"parent_id"`
	Content      string `json:"content"`
	Description  string `json:"description"`
	CommentCount int    `json:"comment_count"`
	Assignee     int    `json:"assignee"`
	Assigner     int    `json:"assigner"`
	Order        int    `json:"order"`
	Priority     int    `json:"priority"`
	URL          string `json:"url"`
}
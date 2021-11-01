package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"net/http"
	"strconv"
	"github.com/thoas/go-funk"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/joho/godotenv"
)

func main() {
	mainRoutin()
}

// メイン
func mainRoutin(){
	if err := initLoadEnv(); err != nil {
		fmt.Errorf("Error: load env:", err)
		return 
	}

	response, err := getTasks()
	if err != nil {
		fmt.Println("Error: getTasks: ", err)
		return
	}

	tasks, err := parseToTasks(response)
	if err != nil {
		fmt.Println("Error: parse to tasks:", err)
		return
	}
	
	bot, err := initLineBot()
	if err != nil {
		fmt.Println("Error: init line bot:", err)
		return
	}
	
	inboxTasksNum := getTasksNum(tasks)
	text := strconv.Itoa(inboxTasksNum) + "個、タスクが残ってるよ！"

	if err := sendMessageToLine(text, bot); err != nil {
		fmt.Println("Error: line bot: send message:", err)
		return 
	}
	return
}

// 環境変数のロードの初期設定
func initLoadEnv() error {
	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("Error: load env:", err)
	}
	return nil
}

// tasksをTodoistから取得
func getTasks() ([]byte, error) {
	url := "https://api.todoist.com/rest/v1/tasks"
	authHeaderName := "Authorization"
	authHeaderValue := "Bearer " + os.Getenv("TODOIST_TOKEN")

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(authHeaderName, authHeaderValue)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error: Request:", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error: Response:", resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)

	return body, nil
}

// レスポンスを []Task に変換
func parseToTasks(response []byte) ([]Task, error)  {
	var tasks []Task
	err := json.Unmarshal(response, &tasks)
	if err != nil {
		return nil, fmt.Errorf("Error: JSON Unmarshal:", err)
	}

	return tasks, nil
}

// Linebotの初期設定
func initLineBot() (*linebot.Client, error) {
	bot, err := linebot.New(
		os.Getenv("LINE_BOT_CHANNEL_SECRET"),
		os.Getenv("LINE_BOT_CHANNEL_TOKEN"),
	)
	if err != nil {
		return nil, fmt.Errorf("Error: creaate line bot:", err)
	}
	return bot, nil
}

// INBOXのタスク個数の取得
func getTasksNum(tasks []Task) int {
	inboxTasks := funk.Filter(tasks, func(x Task) bool {
    return x.ProjectID == 1528596419
	})
	return len(inboxTasks.([]Task))
}

// メッセージをLINEに送る
func sendMessageToLine(text string, bot *linebot.Client) error  {
	message := linebot.NewTextMessage(text)
	_, err := bot.BroadcastMessage(message).Do()
	if err != nil {
		return fmt.Errorf("Error: line bot: boradcast message:", err)
	}
	return nil
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
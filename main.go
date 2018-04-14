package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/go-github/github"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// TgAPI : telegram API url
var TgAPI = "https://api.telegram.org/bot"

// TgToken : telegram bot token
var TgToken = os.Getenv("TGTOKEN")

// ChatID : bot sends messages to this chat
var ChatID = os.Getenv("CHATID")

func main() {
	log.Println("Starting app...")

	// configure tg bot
	bot, err := tb.NewBot(tb.Settings{
		Token:  TgToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	// define handlers for telegram commands
	bot.Handle("/health", func(m *tb.Message) {
		log.Printf("Tg health request from %v", m.Sender)
		bot.Send(m.Chat, "Go!")
	})

	log.Println("Starting bot...")
	go bot.Start()

	// define handlers for web routes
	http.HandleFunc("/", webhookHandler)
	http.HandleFunc("/health", healthHandler)

	// starting web server
	log.Println("Starting web server...")

	Server := http.Server{Addr: ":80"}
	go func() {
		log.Fatal(Server.ListenAndServe())
	}()

	// graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Shutdown signal received, exiting...")
	Server.Shutdown(context.Background())
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	log.Printf("Receive request from: %v, Method: %v", r.RemoteAddr, r.Method)

	// only POST method is accepted
	if r.Method != "POST" {
		log.Println("Only POST method is accepted for webhook processing")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// read body of http request
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	// parse github webhook type
	event, err := github.ParseWebHook(github.WebHookType(r), b)
	if err != nil {
		log.Println(err)
	}

	// variable for message that will be sent to tg
	var Message string

	// handler based on type
	switch e := event.(type) {
	case *github.PullRequestEvent:
		log.Println("New pullrequest:",
			*e.Repo.FullName,
			*e.PullRequest.Title,
			*e.PullRequest.HTMLURL)
		Message = fmt.Sprintf(`Pull request!
Repo: %v
Pullrequest Title: %v
url: %v`, *e.Repo.FullName, *e.PullRequest.Title, *e.PullRequest.URL)

	case *github.ReleaseEvent:
		log.Println("New release:",
			*e.Repo.FullName,
			*e.Release.Author.Login,
			*e.Release.Name,
			*e.Release.TagName,
			*e.Release.Body,
			*e.Release.HTMLURL)

		Message = fmt.Sprintf(`Release! 
Repo: %v
Author: %v
Title: %v
Tag: %v
Description: %v
url: %v`, *e.Repo.FullName,
			*e.Release.Author.Login,
			*e.Release.Name,
			*e.Release.TagName,
			*e.Release.Body,
			*e.Release.HTMLURL)

	case *github.IssuesEvent:
		log.Println(*e.Repo.FullName, *e.Issue.Title, *e.Issue.URL)
		Message = fmt.Sprintf(`Issue!
Repo: %v
Issue: %v
url: %v
`, *e.Repo.FullName, *e.Issue.Title, *e.Issue.URL)

	default:
		log.Printf("Event type %v isn't supported", github.WebHookType(r))
		return
	}

	// send http code 200
	w.WriteHeader(http.StatusOK)

	// send message to tg
	go sendMsg(Message)

	log.Printf("Processing: %.2fs elapsed\n", time.Since(start).Seconds())
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receive health request from: %v", r.RemoteAddr)
	healthMsg := "Go!"
	w.Header().Set("Server", "Golang is awesome!")
	fmt.Fprintf(w, healthMsg)
}

func sendMsg(msg string) {
	start := time.Now()

	log.Println("Sending message to tg...")

	url := TgAPI + TgToken + "/sendMessage"

	// build json message
	m := fmt.Sprintf(`{"chat_id": %v,"disable_web_page_preview": 1,"text": "%v"}`, ChatID, msg)
	var jsonStr = []byte(m)

	// post request to tg api
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		log.Println("Message hasn't been sent. Status:", resp.Status)
		return
	}

	log.Println("Message has been successfully sent to tg")

	log.Printf("Processing: %.2fs elapsed\n", time.Since(start).Seconds())
}

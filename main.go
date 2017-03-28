package main

import (
	"github.com/RecastAI/SDK-golang/recast"
	"github.com/bmizerany/pat"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"github.com/yageek/recast-go-bot-connector"
	"log"
	"os"
)

var (
	rend     *render.Render
	conn     *botconn.Connector
	aiClient *recast.Client
)

func init() {
	rend = render.New()
	aiClient = recast.NewClient(os.Getenv("RECAST_TOKEN"), "fr")
	conf := botconn.ConnConfig{
		Domain:    botconn.RecastAPIDomain,
		BotID:     os.Getenv("BOT_ID"),
		UserSlug:  os.Getenv("USER_SLUG"),
		UserToken: os.Getenv("USER_TOKEN"),
	}
	conn = botconn.New(conf)
}
func main() {

	// Check DB
	_, err := NewStopDB()
	if err != nil {
		panic(err)
	}
	// Message routing
	conn.UseHandler(botconn.MessageHandlerFunc(nextBus))
	// Router
	mux := pat.New()
	mux.Post("/chatbot", conn)

	// Mux
	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":" + os.Getenv("PORT"))
}

func nextBus(w botconn.MessageWriter, m botconn.InputMessage) {

	response, err := aiClient.TextRequest(m.Attachment.Content, nil)
	if err != nil {
		log.Println("Impossible to contact RecastAPI:", err)
		return
	}
	intent, err := response.Intent()
	if err != nil {
		log.Println("Unknown intent:", err)
	} else if intent.Slug == "next-bus" {
		log.Println("Searching element")
	} else if intent.Slug == "greetings" {
		log.Println("Bonjour")
	}
}

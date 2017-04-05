package main

import (
	"fmt"
	"github.com/RecastAI/SDK-Golang/recast"
	"github.com/bmizerany/pat"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"log"
	"os"
)

var (
	rend     *render.Render
	conn     *recast.ConnectClient
	aiClient *recast.RequestClient
	stopDB   *StopDB
)

func init() {
	rend = render.New()
	aiClient = &recast.RequestClient{Token: os.Getenv("RECAST_TOKEN"), Language: "fr"}
}
func main() {

	// Check DB
	db, err := NewStopDB()
	if err != nil {
		panic(err)
	}
	stopDB = db

	conn = recast.NewConnectClient(os.Getenv("RECAST_TOKEN"))
	// Message routing
	conn.UseHandler(recast.MessageHandlerFunc(nextBus))
	// Router
	mux := pat.New()
	mux.Post("/chatbot", conn)

	// Mux
	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":3001")
}

func nextBus(w recast.MessageWriter, m recast.Message) {
	log.Println("Message Received")
	response, err := aiClient.ConverseText(m.Attachment.Content, &recast.ConverseOpts{ConversationToken: m.ConversationId})
	if err != nil {
		log.Println("Impossible to contact RecastAPI:", err)
		return
	}
	fmt.Printf("Response: %+v \n", response)

	var text string
	if len(response.Replies) > 0 {
		text = response.Replies[0]
	} else {
		text = "?????"
	}
	reply := recast.NewTextMessage(text)

	if err := conn.SendMessage(m.ConversationId, reply); err != nil {
		fmt.Println("Error by sending message:", err)
	}

	// Now check what we need to do
	if response.Action.Done && response.Action.Slug == "next-bus-stop" {

	}
}

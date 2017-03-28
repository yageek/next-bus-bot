package main

import (
	"fmt"
	"github.com/bmizerany/pat"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"github.com/yageek/recast-go-bot-connector"
	"os"
)

var (
	rend *render.Render
	conn *botconn.Connector
)

func init() {
	rend = render.New()

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

	output := botconn.OutputMessage{
		Content: "Coucou",
		Kind:    botconn.TextKind,
	}
	err := conn.Broadcast(output)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Response succeeded")
	}
}

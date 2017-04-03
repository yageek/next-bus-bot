package main

import (
	"github.com/RecastAI/SDK-Golang/recast"
	"github.com/bmizerany/pat"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"os"
)

var (
	rend     *render.Render
	conn     *recast.ConnectClient
	aiClient *recast.RequestClient
)

func init() {
	rend = render.New()
	aiClient = &recast.RequestClient{os.Getenv("RECAST_TOKEN"), "fr"}
}
func main() {

	// // Check DB
	// _, err := NewStopDB()
	// if err != nil {
	// 	panic(err)
	// }
	conn = recast.NewConnectClient(os.Getenv("BOT_ID"))
	// Message routing
	conn.UseHandler(recast.MessageHandlerFunc(nextBus))
	// Router
	mux := pat.New()
	mux.Post("/", conn)
	mux.Get("/", http.HandlerFunc(testHandler))

	// Mux
	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":" + os.Getenv("PORT"))
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world \n"))
}
func nextBus(w recast.MessageWriter, m recast.Message) {
	log.Println("Message Received")
	response, err := aiClient.ConverseText(m.Attachment.Content, &recast.ConverseOpts{ConversationToken: m.ConversationId})
	if err != nil {
		log.Println("Impossible to contact RecastAPI:", err)
		return
	}

	var text string
	if len(response.Replies) > 1 {
		text = response.Replies[0]
	} else {
		text = "?????"
	}
	message := recast.Attachment{
		Content: text,
		Type:    "text",
	}
	if err := w.Reply(message); err != nil {
		fmt.Println("Error:", err)
	}
}

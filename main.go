package main
import (
	"fmt"
	"samurai/router"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	r := router.Router() // create router. it will be used to register routes.

	// Create a CORS handler with desired options.
	// it will allow api to be accessed from any origin
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // All origins
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})

	// wrapping router with the CORS handler.
	// wrapping is done to allow api to be accessed from any origin
	handler := c.Handler(r)

	http.Handle("/api/", handler) // registering router with http Handle.
	// it will handle all the incoming requests. "/" means all incoming requests.
	// second parameter is the router. here it is wrapped with CORS handler.

	fmt.Println("Server is running at port 8000...")
	http.ListenAndServe(":8000", nil) // this will start the server.
	// second parameter is the handler. nil means use default handler.
	// default handler is router. so it will use router to handle all the incoming requests.
}
package main

// Import our dependencies. We'll use the standard HTTP library as well as the gorill router for this app
import (
	"github.com/gorilla/mux"

	"net/http"
)

func main() {
	// Hrer we are instantliation the gorilla/mux router
	r := mux.NewRouter()

	// On the default page we will simply serve our static index page

	r.Handle("/", http.FileServer(http.Dir("./views/")))

	// We will setup our server so we can serve static assest like images, css from the /static/{file} froute

	r.PathPrefix("/static/").Handler(http.StripPrefix("/statis/", http.FileServer(http.Dir("./static/"))))

	// Our application will run on port 8080. Here we declare the port and pass in our router.
	http.ListenAndServe(":8080", r)
}
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

	// Our API is going to consist of three routes
	// /status - which we will call to make sure that our API is up and runninng
	// /products - which will retrieve a list of product that the user can leave feedback on
	// /products/{slug}/feedback - which will capture user feedback on products

	r.Handle("/status", NotImplemented).Methods("GET")

	r.Handle("/products", NotImplemented).Methods("GET")

	r.Handle("/products/{slug}/feedback", NotImplemented).Methods("POST")



	// Our application will run on port 8080. Here we declare the port and pass in our router.
	http.ListenAndServe(":8080", r)
}

// Here we are implementing the NotImplemented handler. Whener and API endpoint is hit
// we will simply return the message "Not Implemented"

var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

	w.Write([]byte("Not Implemented"))
})
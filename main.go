package main

// Import our dependencies. We'll use the standard HTTP library as well as the gorill router for this app
import (
	"encoding/json"
  "errors"
  "github.com/auth0/go-jwt-middleware"
  "github.com/dgrijalva/jwt-go"
  "github.com/gorilla/mux"
  "github.com/rs/cors"
  "net/http"

	"github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"context"
	"log"
)

type Response struct {
	Message string `json:"message"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type Product struct {
	Id int

	Name string

	Slug string

	Description string
}

// Create catalog of VR experiences and store them in a slice

var products = []Product{
	Product{Id: 1, Name: "World of Authcraft", Slug: "world-of-authcraft", Description : "Battle bugs and protect yourself from invaders while you explore a scary world with no security"},
  Product{Id: 2, Name: "Ocean Explorer", Slug: "ocean-explorer", Description : "Explore the depths of the sea in this one of a kind underwater experience"},
  Product{Id: 3, Name: "Dinosaur Park", Slug : "dinosaur-park", Description : "Go back 65 million years in the past and ride a T-Rex"},
  Product{Id: 4, Name: "Cars VR", Slug : "cars-vr", Description: "Get behind the wheel of the fastest cars in the world."},
  Product{Id: 5, Name: "Robin Hood", Slug: "robin-hood", Description : "Pick up the bow and arrow and master the art of archery"},
  Product{Id: 6, Name: "Real World VR", Slug: "real-world-vr", Description : "Explore the seven wonders of the world in VR"},
}

// var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 	claims := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)

// 	payload, err := json.Marshal(claims)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(payload)
// })


func main() {

	keyFunc := func(ctx context.Context) (interface{}, error) {
		// Our token must be signed using this data.
		return []byte("secret"), nil
	}

	// Set up the validator.
	jwtValidator, err := validator.New(
		keyFunc,
		validator.HS256,
		"https://<issuer-url>/",
		[]string{"<audience>"},
	)
	if err != nil {
		log.Fatalf("failed to set up the validator: %v", err)
	}

	// Set up the middleware.
	jwtMiddleware := jwtmiddleware.New(jwtValidator.ValidateToken)

	// jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
	// 	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
	// 		//Verify 'aud' claim

	// 		aud := "https://golang-vr"

	// 		checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)

	// 		if !checkAud {
	// 			return token, errors.New("Invalid audience.")
	// 		}

	// 		// Verify 'iss' claim

	// 		iss := "https://golang-vr/"

	// 		checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)

	// 		if !checkIss {
	// 			return token, errors.New("Invalid issuer.")
	// 		}

	// 		cert, err := getPemCert(token)
	// 		if err != nil {
	// 			panic(err.Error())
	// 		}
	// 		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))

	// 		return result, nil

	// 	},

	// 	SigningMethod: jwt.SigningMethodRS256,

	// })

	// Here we are instantliation the gorilla/mux router
	r := mux.NewRouter()

	// On the default page we will simply serve our static index page

	r.Handle("/", http.FileServer(http.Dir("./views/")))

	// We will setup our server so we can serve static assest like images, css from the /static/{file} froute

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Our API is going to consist of three routes
	// /status - which we will call to make sure that our API is up and runninng
	// /products - which will retrieve a list of product that the user can leave feedback on
	// /products/{slug}/feedback - which will capture user feedback on products

	r.Handle("/status", StatusHandler).Methods("GET")

	r.Handle("/products", jwtMiddleware.CheckJWT(ProductsHandler)).Methods("GET")

	r.Handle("/products/{slug}/feedback", jwtMiddleware.CheckJWT(AddFeedbackHandler)).Methods("POST")

	// For dev only - Set up CORS so React client can consume API

	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},

		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})



	// Our application will run on port 8080. Here we declare the port and pass in our router.
	http.ListenAndServe(":8080", r)
}

// Here we are implementing the NotImplemented handler. Whener and API endpoint is hit
// we will simply return the message "Not Implemented"

var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

	w.Write([]byte("Not Implemented"))
})

/* The status handler will be invoked when the user calls the /status route
		it will simply return a string with the message "API is up and running" */

var StatusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("API is up and running"))
})

/* The products handler will be called when the user makes a GET request to the /products endpoint.
This handler will return a list of products available for users to review */

var ProductsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	// Here we are converting the slice of products to JSON
	payload, _ := json.Marshal(products)

	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(payload))
})

/* The feedback handler will add either positive or negative feedback to the product
We would normally save this data to the database - but for this demo, we'll fake it
so that as long as the request is successful and we can match a product to our catalog of products we'll return an OK status. */
var AddFeedbackHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

	var product Product

	vars := mux.Vars(r)

	slug := vars["slug"]

	for _, p := range products {
		if p.Slug == slug {
			product = p
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if product.Slug != "" {
		payload, _ := json.Marshal(product)

		w.Write([]byte(payload))
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
})

func getPemCert(token *jwt.Token)(string, error) {
	cert := ""

	resp, err := http.Get("https://golang-vr/.well-known/jwks.json")

	if err != nil {
		return cert, err
	}

	defer resp.Body.Close()

	var jwks = Jwks{}

	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}
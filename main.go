package main

import (
	"encoding/json" // to encode the data into json when send it to the postman
	"fmt"           // to printing out stuff (your software is connected, etc)
	"log"           // to log the errors if there's any for connecting to the server
	"math/rand"     // create new id for new user
	"net/http"      // create a server in golang
	"strconv"       // to convert the math/rand (int) to str so we able to assign the value to the struct (bcs the id type is string)

	"github.com/gorilla/mux" // import the gorilla mux
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

// passing a pointer of the request that we'll send from postman to this func
func getMovies(w http.ResponseWriter, r *http.Request) {
	// set content type to json
	w.Header().Set("Content-Type", "application/json")
	// Resending all the movies
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			// create a new slice that excludes the movie at the current index.
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	// Resending all the remaining movies
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	// use _ (blank indentifier) because we wont using the index
	for _, item := range movies {
		if item.ID == params["id"] {
			// Sending back the selected movie
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	// get the value to the movie var after decoding it
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(100000000))
	// assign new movie (var movie) to movies slice
	movies = append(movies, movie)
	// return the movie to the front/user
	json.NewEncoder(w).Encode(movie)
}

// just to make it simple for practicing
// should modify only the data that we want to modify if we code irl / work
func updateMovie(w http.ResponseWriter, r *http.Request) {
	// set json content type
	w.Header().Set("Content-Type", "application/json")
	// params
	params := mux.Vars(r)
	// loop over the movie, range
	// delete the movie with the id that we've sent (from the from/user/postman)
	// add a new movie that we'ev also sent
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
		}
	}
}

func main() {
	r := mux.NewRouter()

	// &Director, because we want the reference of the address of the director object
	// & to give the address, * to access that addresss of the pointer
	// just to know that there's a movie when we hit it in the postman for the first time
	movies = append(movies, Movie{ID: "1", Isbn: "438227", Title: "Movie One", Director: &Director{Firstname: "Andrea", Lastname: "Pirlo"}})
	movies = append(movies, Movie{ID: "2", Isbn: "458224", Title: "Movie Two", Director: &Director{Firstname: "Kevin", Lastname: "De Bruyne"}})

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}

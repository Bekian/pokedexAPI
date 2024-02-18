package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

// Data structure for user information
type UserData struct {
	User          string `json:"user"`
	FavoriteColor string `json:"favoritecolor"`
}

type Users struct {
	Users []UserData
}

type pokemon struct {
	ID         string `json:"ID"`
	Name       string `json:"Name"`
	Form       string `json:"Form"`
	Type1      string `json:"Type1"`
	Type2      string `json:"Type2"`
	Total      string `json:"Total"`
	HP         string `json:"HP"`
	Attack     string `json:"Attack"`
	Defense    string `json:"Defense"`
	Sp_Atk     string `json:"Sp_Atk"`
	Sp_Def     string `json:"Sp_Def"`
	Speed      string `json:"Speed"`
	Generation string `json:"Generation"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]string{"message": "Hello root!"}
	tmpl.Execute(w, data)
}

func getUserData(userID string) (UserData, error) {
	// open the file and defer its closing
	file, err := os.Open("data.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// create a reader to read the file and get all the data
	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	// convert the string row number to an int
	rowNumber, err := strconv.Atoi(userID)
	if err != nil {
		return UserData{}, fmt.Errorf("invalid user ID: %w", err)
	}
	// pull the data from the csv using the parsed indexes
	user := data[rowNumber-1][0]
	favoriteColor := data[rowNumber-1][1]
	// return the data as an object
	return UserData{user, favoriteColor}, nil
}

func handleGetUser(w http.ResponseWriter, r *http.Request) {
	// get the data after the user route
	userID := r.URL.Path[len("/user/"):]
	// turn get the stored data at that path
	userData, err := getUserData(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	// turn the returned data into json
	jsonData, err := json.Marshal(userData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshalling data: %v", err)
		return
	}
	// set the content and write the json as a response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	// open the file and defer its closing
	file, err := os.Open("data.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// create a reader to read the file and get all the data
	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	// turn the returned data into json
	jsonData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshalling data: %v", err)
		return
	}
	// set the content and write the json as a response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func getPokemonData(pokeID string) (pokemon, error) {
	// convert the string row number to an int
	rowNumber, err := strconv.Atoi(pokeID)
	// check for row errors
	if err != nil {
		return pokemon{}, fmt.Errorf("invalid pokemon ID parsing error: %w", err)
	}
	var dataFile string
	// due to how the rows are calculated from the path arguments, alternate forms as separate entries aren't possible yet
	switch {
	case rowNumber >= 1 && rowNumber <= 151:
		dataFile = "./trimmedData/gen01trimmed.csv"
	case rowNumber >= 152 && rowNumber <= 251:
		rowNumber -= 151
		dataFile = "./trimmedData/gen02trimmed.csv"
	case rowNumber >= 252 && rowNumber <= 386:
		rowNumber -= 251
		dataFile = "./trimmedData/gen03trimmed.csv"
	case rowNumber >= 387 && rowNumber <= 493:
		rowNumber -= 386
		dataFile = "./trimmedData/gen04trimmed.csv"
	case rowNumber >= 494 && rowNumber <= 649:
		rowNumber -= 493
		dataFile = "./trimmedData/gen05trimmed.csv"
	case rowNumber >= 650 && rowNumber <= 721:
		rowNumber -= 649
		dataFile = "./trimmedData/gen06trimmed.csv"
	case rowNumber >= 722 && rowNumber <= 809:
		rowNumber -= 721
		dataFile = "./trimmedData/gen07trimmed.csv"
	case rowNumber >= 810 && rowNumber <= 905:
		rowNumber -= 809
		dataFile = "./trimmedData/gen08trimmed.csv"
	case rowNumber >= 906 && rowNumber <= 1010:
		rowNumber -= 905
		dataFile = "./trimmedData/gen09trimmed.csv"
	default:
		dataFile = "./trimmedData/gen01trimmed.csv"
	}
	// open the file and read the data
	file, err := os.Open(dataFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	if rowNumber < 1 || rowNumber > 1010 {
		return pokemon{}, fmt.Errorf("pokemon ID out of range 1-1010")
	}

	// pull the data from the csv using the parsed indexes
	p := pokemon{
		ID:         data[rowNumber-1][0],
		Name:       data[rowNumber-1][1],
		Form:       data[rowNumber-1][2],
		Type1:      data[rowNumber-1][3],
		Type2:      data[rowNumber-1][4],
		Total:      data[rowNumber-1][5],
		HP:         data[rowNumber-1][6],
		Attack:     data[rowNumber-1][7],
		Defense:    data[rowNumber-1][8],
		Sp_Atk:     data[rowNumber-1][9],
		Sp_Def:     data[rowNumber-1][10],
		Speed:      data[rowNumber-1][11],
		Generation: data[rowNumber-1][12],
	}
	return p, nil
}

func handleGetPokemon(w http.ResponseWriter, r *http.Request) {
	// get the data after the 'pokemon' path
	pokeID := r.URL.Path[len("/pokemon/"):]
	// get the stored data at that path
	userData, err := getPokemonData(pokeID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	// turn the returned data into json
	jsonData, err := json.Marshal(userData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshalling data: %v", err)
		return
	}
	// set the content and write the json as a response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func handlePanic(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Exit endpoint used, exiting...")
	os.Exit(1)
}

func handleHandlers() {
	http.HandleFunc("/", rootHandler)

	http.HandleFunc("GET /panic", handlePanic)

	http.HandleFunc("GET /users", handleGetUsers)
	http.HandleFunc("GET /user/{id}", handleGetUser)

	http.HandleFunc("GET /pokemon/{id}", handleGetPokemon)
}

func main() {
	// setup handlers
	handleHandlers()

	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

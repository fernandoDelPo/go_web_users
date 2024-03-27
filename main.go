package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

var users []User

func init() {
	users = []User{{
		ID:        "1",
		FirstName: "fernando",
		LastName:  "Del Pozzi",
		Email:     "fernandodelpozzi@example.com",
	}, {
		ID:        "2",
		FirstName: "Leticia",
		LastName:  "Caceres",
		Email:     "leticia@example.com",
	}, {
		ID:        "3",
		FirstName: "Francesca",
		LastName:  "Del Pozzi",
		Email:     "fran@example.com",
	}, {
		ID:        "4",
		FirstName: "Maxima",
		LastName:  "Perez",
		Email:     "mperez@example.com",
	}}
}

func main() {

	http.HandleFunc("/users", UserServer)
	fmt.Println("Starting server on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func UserServer(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		GetAllUser(w)
	case http.MethodPost:
		decode := json.NewDecoder(r.Body)
		var u User
		err := decode.Decode(&u)
		if err != nil {
			MsgResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		PostUser(w, u)

	default:
		InvalidMethod(w)
	}

}

func GetAllUser(w http.ResponseWriter) {
	DataResponse(w, http.StatusOK, users)
}

func PostUser(w http.ResponseWriter, data interface{}) {

	user := data.(User)

	if user.FirstName  == "" {
		MsgResponse(w, http.StatusBadRequest, "First name is required.")
		return
	} 
	
	if user.LastName  == "" {
		MsgResponse(w, http.StatusBadRequest, "Last name is required.")
		return
	}
	if user.Email  == "" {
		MsgResponse(w, http.StatusBadRequest, "Email is required.")
		return
	}

	id := len(users) + 1
	idString := strconv.Itoa(id)
	user.ID = idString
	users = append(users, user)
	DataResponse(w, http.StatusCreated, user)

}

func MsgResponse(w http.ResponseWriter, status int, message string) {

	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "message": "%s"}`, status, message)
}

func DataResponse(w http.ResponseWriter, status int, users interface{}) {
	value, err := json.Marshal(users)
	if err != nil {
		MsgResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "data": %s}`, status, value)
}

func InvalidMethod(w http.ResponseWriter) {

	status := http.StatusNotFound
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "message": "method doesn't exits"}`, status)
}

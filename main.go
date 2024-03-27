package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	},{
		ID:        "4",
		FirstName: "Maxima",
		LastName:  "Perez",
		Email:     "mperez@example.com",
	}}
}

func main() {

	http.HandleFunc("/users", UserServer)

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func UserServer(w http.ResponseWriter, r *http.Request) {
	var status int
	switch r.Method {
	case http.MethodGet:
		GetAllUser(w)
	case http.MethodPost:
		status = 200
		w.WriteHeader(status)
		fmt.Fprintf(w, `{"status": %d, "message": "%s"}`, status, "success method POST")

	default:
		status = 404
		w.WriteHeader(status)
		fmt.Fprintf(w, `{"status": %d, "message": "%s"}`, status, "NOT FOUND")
	}

}

func GetAllUser(w http.ResponseWriter) {
	DataResponse(w, http.StatusOK, users)
}

func DataResponse(w http.ResponseWriter, status int, users interface{}) {
	value, _ := json.Marshal(users)
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "data": %s}`, status, value)
}

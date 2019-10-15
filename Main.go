package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
)
var users []User
var path string
func main() {

	router := mux.NewRouter()
    path = "users.csv"

	users = append(users, User{"1", "login", "pass", "..."})
	users = append(users, User{"2", "login2", "pass2", "..."})
	users = append(users, User{"3", "login3", "pass3", "..."})

	fmt.Println(users)

	router.HandleFunc("/add_user", addUser).Methods("POST")
	router.HandleFunc("/user/{id}", getUser).Methods("GET")
	router.HandleFunc("/upd_user", updateUser).Methods("PUT")
	router.HandleFunc("/del_user/{id}", deleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))

}

type User struct {
	Id string `json: "ID"`
	Login string `json: "login"`
	Password string `json: "password"`
	Introduction string `json: "..."`

}

func getUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	parameters := mux.Vars(request)

	users := getAllUsersFromFile(path)
	for _, user := range users {
		if user.Id == parameters["id"] {
			json.NewEncoder(writer).Encode(user)
			return
		}
	}
}

func addUser(writer http.ResponseWriter, request *http.Request){
	writer.Header().Set("Content-Type", "application/json")
	var user User
	reqBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Print(err)
	}
	user.Id = strconv.Itoa(rand.Intn(100))

	json.Unmarshal(reqBody, &user)
	users = getAllUsersFromFile(path)
	users[len(users)-1].Introduction = users[len(users)-1].Introduction + "\n"
	users = append(users, user)

	err = writeAllUsersToFile(users, path)
	if err!=nil {
		fmt.Println(err)
	}
	json.NewEncoder(writer).Encode(user)
}



func updateUser(writer http.ResponseWriter, request *http.Request){
	writer.Header().Set("Content-Type", "application/json")
	var user User
	reqBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Print(err)
	}
	user.Id = strconv.Itoa(rand.Intn(100))

	json.Unmarshal(reqBody, &user)

	users := getAllUsersFromFile(path)
	for i, us := range  users {
		if us.Id == user.Id{
			users = append(users[:i], users[i+1:]...)
			users = append(users, user)
		}

	}
	err = writeAllUsersToFile(users, path)
	if err!=nil{
		fmt.Println(err)
	}
	json.NewEncoder(writer).Encode(user)
}

func deleteUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	parameters := mux.Vars(request)

	users := getAllUsersFromFile(path)

	for i, user := range users {
		if user.Id == parameters["id"] {
			users = append(users[:i], users[i+1:]...)
		}
	}

	chars := []byte(users[len(users)-1].Introduction)
	chars = chars[:len(chars)-1]
	users[len(users)-1].Introduction = string(chars)

	err := writeAllUsersToFile(users, path)
	if err!=nil {
		fmt.Println(err)
	}
}

func getAllUsersFromFile(path string) []User {
	var users []User
	file, _ := os.Open("users.csv")
	reader := bufio.NewReader(file)
	for {
		text, err := reader.ReadString('\n')
		params := strings.Split(text, ",")
		users = append(users, User{params[0], params[1], params[2], params[3] })
		if err!=nil{
			break
		}
	}
	file.Close()
	return users

}

func writeAllUsersToFile(users []User, path string) error{
	file, err := os.Create(path)
	for _, user := range users {
		file.WriteString(user.Id + "," + user.Login + "," + user.Password + "," + user.Introduction)
	}
	file.Close()
	return err
}
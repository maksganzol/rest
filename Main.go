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
	"time"
)

var path string

func main() {

	fmt.Println("Use command `start` to start server")
	var command string
	for command != "start" {
		fmt.Scan(&command)
	}
	fmt.Println("Server started...")

	router := mux.NewRouter()
	path = "users.csv"

	router.HandleFunc("/user", addUser).Methods("POST")
	router.HandleFunc("/user/{id}", getUser).Methods("GET")
	router.HandleFunc("/user/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/user/{id}", deleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))

}

type User struct {
	Id           string `json: "ID"`
	Login        string `json: "login"`
	Password     string `json: "password"`
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
	writer.WriteHeader(404)
}

func addUser(writer http.ResponseWriter, request *http.Request) {

	var users []User

	writer.Header().Set("Content-Type", "application/json")
	var user User
	reqBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Print(err)
	}

	users = getAllUsersFromFile(path)
	if len(users) != 0 {
		users[len(users)-1].Introduction = users[len(users)-1].Introduction + "\n"
	}

	json.Unmarshal(reqBody, &user)
	if !isValidData(user.Login, user.Password) {
		writer.WriteHeader(415)
		return
	}

	user.Id = strconv.Itoa(rand.Intn(1000)) //Генерирум уникальное id пользователю
	for !isUnicId(user.Id, users) {
		rand.Seed(time.Now().UnixNano())
		user.Id = strconv.Itoa(rand.Intn(1000))
	}

	if isUserExists(user, users) {
		writer.WriteHeader(409)
		return
	}
	users = append(users, user)
	err = writeAllUsersToFile(users, path)
	if err != nil {
		fmt.Println(err)
	}
	json.NewEncoder(writer).Encode(user)
}

func updateUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	parameters := mux.Vars(request)
	var user User
	reqBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Print(err)
	}

	users := getAllUsersFromFile(path)

	json.Unmarshal(reqBody, &user)
	if !isUserExists(user, users) {
		writer.WriteHeader(404)
		return
	}

	user.Id = parameters["id"]

	for i, us := range users {
		if us.Id == user.Id {
			users = append(users[:i], users[i+1:]...)
			users = append(users, user)
			json.NewEncoder(writer).Encode(user)
		}

	}
	err = writeAllUsersToFile(users, path)
	if err != nil {
		fmt.Println(err)
	}

}

func deleteUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	parameters := mux.Vars(request)

	users := getAllUsersFromFile(path)

	delete := false
	for i, user := range users {
		if user.Id == parameters["id"] {
			users = append(users[:i], users[i+1:]...)
			delete = true
		}
	}

	if !delete {
		writer.WriteHeader(404)
		return
	}

	users[len(users)-1].Introduction = cutLastChar(users[len(users)-1].Introduction)

	err := writeAllUsersToFile(users, path)
	if err != nil {
		fmt.Println(err)
	}
}

func getAllUsersFromFile(path string) []User {
	var users []User
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		file, _ = os.Create(path)
	}
	reader := bufio.NewReader(file)
	for {
		text, err := reader.ReadString('\n')
		params := strings.Split(text, ",")
		if len(params) < 4 {
			break
		}
		users = append(users, User{params[0], params[1], params[2], params[3]})
		if err != nil {
			break
		}
	}
	file.Close()
	return users

}

func writeAllUsersToFile(users []User, path string) error {
	file, err := os.Create(path)
	for _, user := range users {
		file.WriteString(user.Id + "," + user.Login + "," + user.Password + "," + user.Introduction)
	}
	file.Close()
	return err
}

func isUserExists(user User, users []User) bool {
	exists := false
	for _, us := range users {
		if us.Login == user.Login {
			exists = true
		}
	}
	return exists
}

func isValidData(data ...string) bool {
	valid := "abcdefghigklmnopqrstuvwxyz_1234567890"
	for _, line := range data {
		for _, char := range line {
			if !strings.Contains(valid, string(char)) {
				return false
			}
		}
	}
	return true
}

func cutLastChar(s string) string {
	chars := []byte(s)
	chars = chars[:len(chars)-1]
	s = string(chars)
	return s
}

func isUnicId(id string, users []User) bool {
	for _, user := range users {
		if user.Id == id {
			return false
		}
	}
	return true
}

package httpserver

import (
	"encoding/json"
	"entrytask/internal/communication/pb"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	Username     string `json:"username"`
	Nickname     string `json:"nickname"`
	ProfileImage []byte `json:"profile_image"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	client := GetClient()
	request := pb.LoginRequest{Username: req.Username, Password: req.Password}
	response := pb.LoginResponse{}
	err = client.Call("UserService.Login", &request, &response)
	if err != nil {
		log.Fatal("call UserService.Login error:", err.Error())
	}
	fmt.Println(response.Message)
	json.NewEncoder(w).Encode(response)
}

package tcp

import (
	"entrytask/internal/communication/pb"
	"log"
)

type UserService struct{}

func (us *UserService) Login(req pb.LoginRequest, resp *pb.LoginResponse) error {
	user := getUser(req.Username)
	err := compare(user.Password, req.Password)
	if err != nil {
		resp.Success = false
		resp.Message = "Wrong Password"
	} else {
		resp.Success = true
		resp.Message = "success"
	}
	return nil
}

func getUser(username string) *User {
	// Retrieve the user profile from the database (dummy implementation for demonstration)
	// Replace this with your actual database query
	row := db.QueryRow("SELECT * FROM users WHERE username = ?", username)

	var user User
	err := row.Scan(&user.Username, &user.Password, &user.Nickname, &user.ProfileImage)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &user
}

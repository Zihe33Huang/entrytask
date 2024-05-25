package tcp

import (
	"entrytask/internal/communication/pb"
	"errors"
	"fmt"
	"log"
)

type UserService struct{}

func (us *UserService) Login(req pb.LoginRequest, resp *pb.LoginResponse) error {
	user := getUserByUsername(req.Username)
	if user == nil {
		return errors.New("User not found!")
	}
	err := compare(user.Password, req.Password)

	if err != nil {
		resp.Success = false
		return fmt.Errorf("Password doesn't match!")
	}

	jwt, err := generateJwt(user.Id)
	if err != nil {
		resp.Success = false
		return err
	}
	resp.Token = jwt
	resp.Success = true
	return nil
}

func (us *UserService) GetProfile(req pb.ProfileRequest, resp *pb.ProfileResponse) error {
	id := req.UserId
	user := getUserById(id)
	if user == nil {
		err := errors.New("user not found")
		return err
	}
	resp.Username = user.Username
	resp.Nickname = user.Nickname
	resp.ProfileImg = user.ProfileImage
	return nil
}

func (us *UserService) UpdateNickname(req pb.NicknameUpdateRequest, resp *pb.NicknameUpdateResponse) error {

	err := updateNickname(req.UserId, req.Nickname)
	if err != nil {
		return err
	}
	resp.Success = true
	return nil
}

func getUserByUsername(username string) *User {
	row := db.QueryRow("SELECT * FROM users WHERE username = ?", username)

	var user User
	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Nickname, &user.ProfileImage)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &user
}

// AuthMiddleware is a middleware function to intercept and validate the JWT

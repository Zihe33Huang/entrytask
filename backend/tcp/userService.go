package tcp

import (
	"entrytask/backend/communication/pb"
	"errors"
	"fmt"
	"regexp"
)

type UserService struct{}

func (us *UserService) Login(req pb.LoginRequest, resp *pb.LoginResponse) error {
	// Validate request fields
	if req.Username == "" {
		return errors.New("Username is required")
	}
	if req.Password == "" {
		return errors.New("Password is required")
	}

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
	// 1. Validate Data
	ok := isValidNickname(req.Nickname)
	if !ok {
		return errors.New("Invalid Nickname")
	}

	err := updateNickname(req.UserId, req.Nickname)
	if err != nil {
		return err
	}
	resp.Success = true
	return nil
}

func (us *UserService) UpdateProfileImg(req pb.ProfileImgUpdateRequest, resp *pb.ProfileImgUpdateResponse) error {

	err := updateProfileImg(req.UserId, req.ProfileImg)
	if err != nil {
		return err
	}
	resp.Success = true
	return nil
}

// isValidNickname validates the given nickname according to the defined rules.
func isValidNickname(nickname string) bool {
	// Define a regular expression pattern for common nicknames
	pattern := `.*`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Check if the length of the nickname is between 1 and 20 characters
	if len(nickname) < 3 || len(nickname) > 10 {
		return false
	}

	// Use MatchString to check if the nickname matches the pattern
	return re.MatchString(nickname)
}

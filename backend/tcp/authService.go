package tcp

import (
	"entrytask/backend/communication/pb"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type AuthService struct{}

func (authService *AuthService) ValidateToken(req pb.AuthRequest, resp *pb.AuthResponse) error {
	// Parse and validate the token
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil // Validate the token with the secret key
	})

	if err != nil || !token.Valid {
		resp.Success = false
		return fmt.Errorf("Token Validation Fails")
	}

	resp.Success = true
	claims := token.Claims.(jwt.MapClaims)
	userId := uint64(claims["userId"].(float64))
	//if !ok {
	//	return fmt.Errorf("jwt error")
	//}
	resp.UserId = userId
	return nil
}

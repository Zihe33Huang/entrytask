package httpserver

import (
	"context"
	"encoding/json"
	"entrytask/backend/communication/pb"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

var ServerUrl = "localhost:8888"

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ProfileResponse struct {
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	ProfileImg string `json:"profile_img"`
}

type UpdateNicknameRequest struct {
	Nickname string `json:"nickname"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		error(w, http.StatusBadRequest, "Invalid request method", nil)
		return
	}

	var logReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&logReq)
	if err != nil {
		error(w, http.StatusInternalServerError, "", nil)
		return
	}

	if logReq.Username == "" || logReq.Password == "" {
		error(w, http.StatusBadRequest, "Username or Password can not be blank", nil)
		return
	}

	pbReq := pb.LoginRequest{Username: logReq.Username, Password: logReq.Password}
	pbResp := pb.LoginResponse{}
	err = client.Call("UserService.Login", &pbReq, &pbResp)
	if err != nil {
		log.Println("rpc client: call UserService.Login error:", err.Error())
		error(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	w.Header().Set("Authorization", "Bearer "+pbResp.Token)
	success(w, "Login successfully", nil)
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		error(w, http.StatusBadRequest, "Invalid request method", nil)
		return
	}
	ctx := r.Context()
	userId, ok := ctx.Value("userId").(uint64)
	if !ok {
		error(w, http.StatusBadRequest, "Invalid user id", nil)
		return
	}

	pbReq := pb.ProfileRequest{UserId: userId}
	pbResp := pb.ProfileResponse{}
	err := client.Call("UserService.GetProfile", &pbReq, &pbResp)
	if err != nil {
		log.Println("rpc client: call UserService.Profile error:", err.Error())
		error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	profileResp := ProfileResponse{}
	profileResp.Username = pbResp.Username
	profileResp.Nickname = pbResp.Nickname
	profileResp.ProfileImg = pbResp.ProfileImg
	success(w, "ok", profileResp)
}

func UpdateNicknameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		error(w, http.StatusBadRequest, "Invalid request method", nil)
		return
	}
	ctx := r.Context()
	userId, ok := ctx.Value("userId").(uint64)
	if !ok {
		error(w, http.StatusBadRequest, "Invalid User Id", nil)
		return
	}

	var req UpdateNicknameRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		error(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	pbReq := pb.NicknameUpdateRequest{UserId: userId, Nickname: req.Nickname}
	pbResp := pb.NicknameUpdateResponse{}
	err = client.Call("UserService.UpdateNickname", &pbReq, &pbResp)
	if err != nil {
		log.Println("rpc client: call UserService.Profile error:", err.Error())
		error(w, http.StatusBadGateway, err.Error(), nil)
		return
	}
	success(w, "Update Nickname Successfully", nil)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the JWT token from the request header
		tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		pbReq := pb.AuthRequest{Token: tokenString}
		pbResp := pb.AuthResponse{}
		err := client.Call("AuthService.ValidateToken", &pbReq, &pbResp)
		if err != nil {
			log.Println("rpc client: call AuthService.ValidateToken error:", err.Error())
			error(w, http.StatusUnauthorized, err.Error(), nil)
			return
		}
		ctx := context.WithValue(r.Context(), "userId", pbResp.UserId)
		r = r.WithContext(ctx)

		// If the token is valid, call the next handler
		next(w, r)
	}
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// only POST method is allowed
	if r.Method != http.MethodPost {
		error(w, http.StatusBadRequest, "Invalid request method", nil)
		return
	}

	// Get userId from context
	ctx := r.Context()
	userId, ok := ctx.Value("userId").(uint64)
	if !ok {
		//http.Error(w, "Invalid user id", http.StatusBadRequest)
		error(w, http.StatusBadRequest, "Invalid user id", nil)
		return
	}

	// 1. Read image from form
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Println("Read file error")
		error(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	defer file.Close()

	filePath := "/Users/zihehuang/Downloads/img/" + strconv.FormatUint(userId, 10) + path.Ext(fileHeader.Filename)
	// 2. Create image
	newFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Create file error", http.StatusInternalServerError)
		return
	}
	defer newFile.Close()

	// 3. Write image into file system
	_, err = io.Copy(newFile, file)
	if err != nil {
		error(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// 4. Store file path into database
	pbReq := pb.ProfileImgUpdateRequest{UserId: userId, ProfileImg: filePath}
	pbResp := pb.ProfileImgUpdateResponse{}
	err = client.Call("UserService.UpdateProfileImg", &pbReq, &pbResp)
	if err != nil {
		log.Println("rpc client: call UserService.UpdateProfileImg error:", err.Error())
		//http.Error(w, err.Error(), http.StatusBadRequest)
		error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	success(w, "ok", nil)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	// Only GET method is allowed
	if r.Method != http.MethodGet {
		error(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	ctx := r.Context()
	userId, ok := ctx.Value("userId").(uint64)
	if !ok {
		error(w, http.StatusBadRequest, "Invalid user id", nil)
		return
	}
	// 1. Ask tcp server for image path
	pbReq := pb.ProfileRequest{UserId: userId}
	pbResp := pb.ProfileResponse{}
	err := client.Call("UserService.GetProfile", &pbReq, &pbResp)
	if err != nil {
		log.Println("rpc client: call UserService.Profile error:", err.Error())
		error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// 2. Read image from file system
	file, err := os.Open(pbResp.ProfileImg)
	if err != nil {
		error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	defer file.Close()

	w.Header().Add("Content-type", "application/octet-stream")
	w.Header().Add("content-disposition", "attachment;")

	// 3. Write image into network
	_, err = io.Copy(w, file)
	if err != nil {
		error(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
}

func success(w http.ResponseWriter, message string, data interface{}) {
	loginResp := Response{Status: http.StatusOK, Message: message, Data: data}
	json.NewEncoder(w).Encode(loginResp)
}

func error(w http.ResponseWriter, status int, message string, data interface{}) {
	loginResp := Response{Status: status, Message: message, Data: data}
	json.NewEncoder(w).Encode(loginResp)
}

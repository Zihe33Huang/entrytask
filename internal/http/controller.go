package httpserver

import (
	"context"
	"encoding/json"
	"entrytask/internal/communication/pb"
	"fmt"
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
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

type UpdateNicknameRequest struct {
	Nickname string `json:"nickname"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	var logReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&logReq)
	if err != nil {
		loginResp := Response{Status: http.StatusInternalServerError, Message: ""}
		json.NewEncoder(w).Encode(loginResp)
		return
	}
	// todo 验证数据

	pbReq := pb.LoginRequest{Username: logReq.Username, Password: logReq.Password}
	pbResp := pb.LoginResponse{}
	err = client.Call("UserService.Login", &pbReq, &pbResp)
	if err != nil {
		log.Println("rpc client: call UserService.Login error:", err.Error())
		loginResp := Response{Status: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(loginResp)
		return
	}
	w.Header().Set("Authorization", "Bearer "+pbResp.Token)
	success(w, "Login successfully", nil)
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	userId, ok := ctx.Value("userId").(uint64)
	if !ok {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	pbReq := pb.ProfileRequest{UserId: userId}
	pbResp := pb.ProfileResponse{}
	err := client.Call("UserService.GetProfile", &pbReq, &pbResp)
	if err != nil {
		log.Println("rpc client: call UserService.Profile error:", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profileResp := ProfileResponse{}
	profileResp.Username = pbResp.Username
	profileResp.Nickname = pbResp.Nickname
	success(w, "ok", profileResp)
}

func UpdateNicknameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	userId, ok := ctx.Value("userId").(uint64)
	if !ok {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	var req UpdateNicknameRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	pbReq := pb.NicknameUpdateRequest{UserId: userId, Nickname: req.Nickname}
	pbResp := pb.NicknameUpdateResponse{}
	err = client.Call("UserService.UpdateNickname", &pbReq, &pbResp)
	if err != nil {
		log.Println("rpc client: call UserService.Profile error:", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	success(w, "ok", nil)
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
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), "userId", pbResp.UserId)
		r = r.WithContext(ctx)
		// If the token is valid, call the next handler
		next(w, r)
	}
}

// Todo 改成response
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	//文件上传只允许POST方法
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	userId, ok := ctx.Value("userId").(uint64)
	fmt.Println()
	if !ok {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	//从表单中读取文件
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Println("Read file error")
		http.Error(w, "Read file error", http.StatusInternalServerError)
		return
	}
	//defer 结束时关闭文件
	defer file.Close()
	log.Println("filename: " + fileHeader.Filename)

	//创建文件
	newFile, err := os.Create("/Users/zihehuang/Downloads/img/" + strconv.FormatUint(userId, 10) + path.Ext(fileHeader.Filename))
	if err != nil {
		http.Error(w, "Create file error", http.StatusInternalServerError)
		return
	}
	//defer 结束时关闭文件
	defer newFile.Close()

	//将文件写到本地
	_, err = io.Copy(newFile, file)
	if err != nil {
		http.Error(w, "Write file error", http.StatusInternalServerError)
		return
	}
	success(w, "ok", nil)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	//文件上传只允许GET方法
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	userId, ok := ctx.Value("userId").(uint64)
	if !ok {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	//文件名
	filename := strconv.FormatUint(userId, 10)
	if filename == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	log.Println("filename: " + filename)
	//打开文件
	// todo 从数据库拿到文件地址
	file, err := os.Open("/Users/zihehuang/Downloads/img/" + filename + ".png")
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	//结束后关闭文件
	defer file.Close()

	//设置响应的header头
	w.Header().Add("Content-type", "application/octet-stream")
	w.Header().Add("content-disposition", "attachment; filename=\""+filename+"\"")
	//将文件写至responseBody
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Write file error", http.StatusInternalServerError)
		return
	}
}

func success(w http.ResponseWriter, message string, data interface{}) {
	loginResp := Response{Status: http.StatusOK, Message: message, Data: data}
	json.NewEncoder(w).Encode(loginResp)
}

package tcp

import (
	"database/sql"
	"log"
)

var (
	stmtGetUserById       *sql.Stmt
	stmtGetUserByUsername *sql.Stmt
	stmtUpdateProfileImg  *sql.Stmt
	stmtUpdateNickname    *sql.Stmt
)

func SetStmtGetUserById(stmt *sql.Stmt) {
	stmtGetUserById = stmt
}

func SetStmtGetUserByUsername(stmt *sql.Stmt) {
	stmtGetUserByUsername = stmt
}

func SetStmtUpdateProfileImg(stmt *sql.Stmt) {
	stmtUpdateProfileImg = stmt
}

func SetStmtUpdateNickname(stmt *sql.Stmt) {
	stmtUpdateNickname = stmt
}

type User struct {
	Id           uint64 `json:"id"`
	Username     string `json:"username"`
	Nickname     string `json:"nickname"`
	Password     string `json:"password"`
	ProfileImage string `json:"profile_image"`
}

func getUserById(id uint64) *User {
	var user User
	err := stmtGetUserById.QueryRow(id).Scan(&user.Username, &user.Nickname, &user.ProfileImage)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle not found error
			return nil
		}
		log.Println("Error fetching user:", err)
		return nil
	}

	return &user
}

func getUserByUsername(username string) *User {
	row := stmtGetUserByUsername.QueryRow(username)

	var user User
	//err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Nickname, &user.ProfileImage)
	err := row.Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		log.Println("database Error:", err)
		return nil
	}

	return &user
}

func updateNickname(id uint64, nickname string) error {
	_, err := stmtUpdateNickname.Exec(nickname, id)
	return err
}

func updateProfileImg(id uint64, profileImg string) error {
	_, err := stmtUpdateProfileImg.Exec(profileImg, id)
	return err
}

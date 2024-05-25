package tcp

import "log"

type User struct {
	Id           uint64 `json:"id"`
	Username     string `json:"username"`
	Nickname     string `json:"nickname"`
	Password     string `json:"password"`
	ProfileImage string `json:"profile_image"`
}

func getUserById(id uint64) *User {
	row := db.QueryRow("SELECT username, nickname, profile_image FROM users WHERE id = ?", id)

	var user User
	err := row.Scan(&user.Username, &user.Nickname, &user.ProfileImage)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &user
}

func updateNickname(id uint64, nickname string) error {
	stmt, err := db.Prepare("UPDATE users SET nickname = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(nickname, id)
	return err

}

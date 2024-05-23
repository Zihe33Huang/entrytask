package tcp

type User struct {
	Username     string `json:"username"`
	Nickname     string `json:"nickname"`
	Password     string `json:"password"`
	ProfileImage string `json:"profile_image"`
}

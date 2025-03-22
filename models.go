package main

type ShortLink struct {
	Handle string `json:"handle"`
	Target string `json:"target"`
	Owner  string `json:"owner"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

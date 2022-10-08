package entity

type Message struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	TokenString string `json:"token"`
}

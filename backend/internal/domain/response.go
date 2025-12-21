package domain

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ResponseToken struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

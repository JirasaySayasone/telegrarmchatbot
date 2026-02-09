package model

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Message struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

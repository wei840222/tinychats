// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Message struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt"`
	User      *User  `json:"user"`
}

type NewMessage struct {
	Text string `json:"text"`
}

type User struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	AvatarURL *string `json:"avatarUrl"`
}

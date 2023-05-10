package models

type User struct {
	UserId  int64  `json:"userid"`
	Genre   string `json:"genre"`
	Sounder string `json:"sounder"`
	Book    string `json:"book"`
	Counter int    `json:"counter"`
}

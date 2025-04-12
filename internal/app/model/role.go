package model

type Role struct {
	ID   int    `json:"id" `
	Projectlist []Project  `json:"projectlist"`
}


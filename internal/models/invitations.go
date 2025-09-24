package models

import "time"

type InvitationModel struct {
	UserId   int           `json:"user_id"`
	Token    string        `json:"token"`
	ExpireAt time.Duration `json:"expire_at"`
}

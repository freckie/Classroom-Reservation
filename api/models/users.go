package models

type UsersGetResponse struct {
	UsersCount int            `json:"users_count"`
	Users      []UsersGetItem `json:"users"`
}

type UsersGetItem struct {
	UserID    int64  `json:"user_id"`
	UserEmail string `json:"user_email"`
	IsSuper   bool   `json:"is_super"`
}

type UsersPostResponse struct {
	UserID int64 `json:"user_id"`
}

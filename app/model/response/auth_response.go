package response

import "dc-transaction-service/app/model/entity"

type AuthRegisterResponse struct {
	GeneralResponse
	UID string `json:"uid,omitempty"`
}

type AuthLoginResponse struct {
	GeneralResponse
	Token        string      `json:"token,omitempty"`
	RefreshToken string      `json:"refresh_token,omitempty"`
	Data         entity.User `json:"data,omitempty"`
}

type AuthLogoutResponse struct {
	GeneralResponse
}

type UserResponse struct {
	GeneralResponse
}

type UserDataResponse struct {
	GeneralResponse
	Data entity.User `json:"data,omitempty"`
}

package response

import "wh-auth-service/app/model/entity"

type RoleResponse struct {
	GeneralResponse
}

type RoleDataResponse struct {
	GeneralResponse
	Data entity.Role `json:"data,omitempty"`
}

type RoleListResponse struct {
	GeneralResponse
	Data []entity.Role `json:"data,omitempty"`
}

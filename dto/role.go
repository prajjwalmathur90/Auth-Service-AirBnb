package dto

type CreateRoleRequestDto struct {
	Name        string `json:"name" validate:"required,min=2,max=50"`
	Description string `json:"description" validate:"max=255"`
}

type UpdateRoleRequestDto struct {
	Name        string `json:"name" validate:"required,min=2,max=50"`
	Description string `json:"description" validate:"max=255"`
}

type AddPermissionToRoleRequestDto struct {
	PermissionId int64 `json:"permission_id" validate:"required,gt=0"`
}

type AssignRoleToUserRequestDto struct {
	RoleId int64 `json:"role_id" validate:"required,gt=0"`
}

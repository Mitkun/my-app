package repository

import (
	"github.com/google/uuid"
	"my-app/module/user/domain"
	"my-app/utils"
	"time"
)

type UserDTO struct {
	Id        uuid.UUID `gorm:"column:id;"`
	FirstName string    `gorm:"column:first_name;"`
	LastName  string    `gorm:"column:last_name;"`
	Email     string    `gorm:"column:email;"`
	Password  string    `gorm:"column:password;"`
	Salt      string    `gorm:"column:salt;"`
	Role      string    `gorm:"column:role;"`
	Status    string    `gorm:"column:status;"`
	Avatar    *string   `gorm:"column:avatar;"`
}

func (dto *UserDTO) ToEntity() (user *domain.User, err error) {
	return domain.NewUser(dto.Id, dto.FirstName, dto.LastName, dto.Email, dto.Password, dto.Salt, domain.GetRole(dto.Role), dto.Status, utils.StringFromPointer(dto.Avatar))
}

type SessionDTO struct {
	Id           uuid.UUID `gorm:"column:id;"`
	UserId       uuid.UUID `gorm:"column:user_id;"`
	RefreshToken string    `gorm:"column:refresh_token;"`
	AccessExpAt  time.Time `gorm:"column:access_exp_at;"`
	RefreshExpAt time.Time `gorm:"column:refresh_exp_at;"`
}

func (dto SessionDTO) ToEntity() (*domain.Session, error) {
	s := domain.NewSession(dto.Id, dto.UserId, dto.RefreshToken, dto.AccessExpAt, dto.RefreshExpAt)

	return s, nil
}

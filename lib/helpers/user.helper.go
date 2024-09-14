package helpers

import (
	"github.com/cbstorm/wyrstream/lib/entities"
	"github.com/cbstorm/wyrstream/lib/utils"
)

type UserHelper struct {
	user *entities.UserEntity
}

func NewUserHelper(user *entities.UserEntity) *UserHelper {
	return &UserHelper{user: user}
}

func (h *UserHelper) HashPassword() error {
	hashed, err := utils.BcryptHash(h.user.Password)
	if err != nil {
		return err
	}
	h.user.Password = hashed
	return nil
}

package helpers

import (
	"github.com/cbstorm/wyrstream/lib/entities"
	"github.com/cbstorm/wyrstream/lib/repositories"
	"github.com/cbstorm/wyrstream/lib/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewVODsHelper(vods *[]*entities.VodEntity) *VODsHelper {
	return &VODsHelper{
		vods: vods,
	}
}

type VODsHelper struct {
	vods *[]*entities.VodEntity
}

func (h *VODsHelper) ResolveOwner() error {
	user_ids := utils.Map(h.vods, func(a *entities.VodEntity, b int) primitive.ObjectID {
		return a.OwnerId
	})
	users := make([]*entities.UserEntity, 0)
	if err := repositories.GetUserRepository().FindManyByIds(*user_ids, &users); err != nil {
		return err
	}
	users_key_by_id := utils.KeyBy(&users, func(a *entities.UserEntity) string {
		return a.Id.Hex()
	})
	for _, v := range *h.vods {
		v.Owner = (*users_key_by_id)[v.OwnerId.Hex()]
	}
	return nil
}

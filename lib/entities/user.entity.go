package entities

type UserEntity struct {
	PersonEntity `bson:",inline"`
	Name         string `bson:"name" json:"name"`
	Email        string `bson:"email" json:"email"`
	Password     string `bson:"password" json:"-"`
}

func NewUserEntity() *UserEntity {
	user := &UserEntity{}
	user.New()
	return user
}

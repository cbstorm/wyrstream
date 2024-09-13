package entities

type UserEntity struct {
	BaseEntity `bson:",inline"`
	Name       string `bson:"name" json:"name"`
	Email      string `bson:"email" json:"email"`
	Password   string `bson:"password" json:"-"`
}

func NewUser() *UserEntity {
	user := &UserEntity{}
	user.NewId()
	user.SetTime()
	return user
}

package entities

type AdminEntity struct {
	BaseEntity `bson:",inline"`
	Email      string `bson:"email,omitempty" json:"email,omitempty"`
	Password   string `bson:"password,omitempty"`
	AvatarUrl  string `bson:"avatarUrl,omitempty" json:"avatarUrl,omitempty"`
}

func NewAdminEntity() *AdminEntity {
	admin := &AdminEntity{}
	admin.NewId()
	admin.SetTime()
	return admin
}

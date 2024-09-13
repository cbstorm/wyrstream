package enums

const (
	USER_ROLE_ADMIN = "ADMIN"
	USER_ROLE_USER  = "USER"
)

const (
	USER_STATUS_ACTIVE   = "ACTIVE"
	USER_STATUS_INACTIVE = "IN_ACTIVE"
)

type UserGenders string

const (
	USER_GENDER_MALE   UserGenders = "MALE"
	USER_GENDER_FEMALE UserGenders = "FEMALE"
)

package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

func AssertDefaultRef[T any](v *T, d *T) *T {
	if v == nil {
		return d
	}
	return v
}

func TernaryOp[T any](condition bool, then T, el T) T {
	if condition {
		return then
	}
	return el
}

func CopyMap(source *map[string]interface{}, destination *map[string]interface{}, ignore_keys []string) {
	ignore_keys_map := map[string]int{}
	for _, v := range ignore_keys {
		ignore_keys_map[v] = 1
	}
	for k, v := range *source {
		if ignore_keys_map[k] == 1 {
			continue
		}
		(*destination)[k] = v
	}
}

func Haversine(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	const earthRadiusKm = 6371e3

	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad

	a := (math.Sin(dlat/2) * math.Sin(dlat/2)) + (math.Cos(lat1Rad) * math.Cos(lat2Rad) * math.Sin(dlon/2) * math.Sin(dlon/2))
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadiusKm * c

	return distance
}

func MD5Sum(str string) string {
	token_sum := md5.Sum([]byte(str))
	return hex.EncodeToString(token_sum[:])
}

func IsValidEmailAddress(email string) bool {
	if _, err := mail.ParseAddress(email); err != nil {
		return false
	}
	return true
}

func BcryptHash(text string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(text), 8)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
func BcryptMatch(hashed string, text string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(text)); err != nil {
		return false
	}
	return true
}

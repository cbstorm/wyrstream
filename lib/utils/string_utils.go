package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"
	"unicode"
	"unicode/utf8"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func DefaultString(val string, defaultVal string) string {
	if val != "" {
		return val
	}
	return defaultVal
}

func StringLength(text string) int {
	return utf8.RuneCountInString(text)
}

func Uniq(t *[]string) []string {
	str_map := map[string]int{}
	for _, v := range *t {
		str_map[v] = 1
	}
	out := []string{}
	for s := range str_map {
		out = append(out, s)
	}
	return out
}

func UniqObjectId(ids *[]primitive.ObjectID) []primitive.ObjectID {
	hex_ids := Map(ids, func(i primitive.ObjectID, idx int) string {
		return i.Hex()
	})
	uniq_ids := Uniq(hex_ids)
	uniq_obj_ids := []primitive.ObjectID{}
	for _, v := range uniq_ids {
		obj_id, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			obj_id = primitive.NewObjectID()
		}
		uniq_obj_ids = append(uniq_obj_ids, obj_id)
	}
	return uniq_obj_ids
}

func NgramTokenize(text string, min_gram int, max_gram int) []string {
	text = strings.Join(strings.Fields(text), "")
	text = strings.ToLower(text)
	text = NormalizeUTF8String(text)
	var ngrams []string
	textLength := utf8.RuneCountInString(text)
	for i := 0; i < textLength; i++ {
		for j := min_gram; j <= max_gram && i+j <= textLength; j++ {
			ngram := text[i : i+j]
			ngrams = append(ngrams, ngram)
		}
	}
	return ngrams
}

func NormalizeUTF8String(text string) string {
	text = strings.ToLower(text)
	trans := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(trans, text)
	if err != nil {
		return ""
	}
	result = strings.ReplaceAll(result, "đ", "d")
	result = strings.ReplaceAll(result, "Đ", "D")
	return result
}

func ContainPattern(source string, target string) bool {
	target_fields_count := len(strings.Fields(target))
	source_fields := strings.Fields(strings.ToLower(source))
	if len(source_fields) < target_fields_count {
		return false
	}
	for i := 0; i < len(source_fields)-target_fields_count+1; i++ {
		if strings.Join(source_fields[i:i+target_fields_count], " ") == target {
			return true
		}
	}
	return false
}

func StringRand(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		ret[i] = letters[num.Int64()]
	}

	return string(ret)
}

func Base64Encode(inp string) string {
	data := []byte(inp)
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(dst, data)
	return string(dst)
}

func Base64Decode(inp string) (string, error) {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(inp)))
	n, err := base64.StdEncoding.Decode(dst, []byte(inp))
	if err != nil {
		return "", err
	}
	dst = dst[:n]
	return fmt.Sprintf("%q\n", dst), nil
}

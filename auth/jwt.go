package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var Jwtkey = []byte(os.Getenv("SECRET"))

func TokenGeneration(id string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECERET")))
	if err != nil {
		panic(err)
	}
	return tokenString

}

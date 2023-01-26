package midilware

import (
	"fmt"

	"net/http"
	"os"

	"time"

	
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func UserAuth(c *gin.Context) {
	//Get the cookie off req
	tokenString, err := c.Cookie("Autherization")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"User": "logout",
		})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Decode/validate it
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECERET")), nil
	})

	if err != nil {
		c.JSON(500, gin.H{
			"erroe": "Bad request",
		})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check the exp

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Set("userid", claims["sub"])

		c.Next()

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)

	}
}

func AdminAuth(c *gin.Context) {
	//Get the cookie off req
	tokenString, err := c.Cookie("AdminAutherization")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invaid access",
		})

		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// Decode/validate it
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECERET")), nil
	})
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Admin logout",
		})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check the exp

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		
		// Atach to req
		c.Set("adminid",claims["sub"])

		// continuew
		c.Next()

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)

	}
}

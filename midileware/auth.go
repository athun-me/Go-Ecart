package midilware

import (
	"fmt"
	"net/http"
	"os"

	"time"

	"github.com/athunlal/config"
	"github.com/athunlal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func RequirAuth(c *gin.Context) {
	//Get the cookie off req
	tokenString, err := c.Cookie("Autherization")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"User" : "logout",
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

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check the exp

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// Find the user with token  sub
		var user models.User
		Db := config.DBconnect()
	
		result := Db.First(&user,"email LIKE ?", claims["sub"])
		if result.Error != nil {
			fmt.Println(result.Error)
		}
		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// Atach to req
		c.Set("user", user)

		// continuew
		c.Next()

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		
	}

}

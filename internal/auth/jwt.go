package auth

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

var secretKey = []byte("minha-chave-secreta")

func GenerateJWT(userId string) (string, error) {
	claims := Claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().UTC().Add(time.Minute * 10)},
			IssuedAt:  &jwt.NumericDate{Time: time.Now().UTC()},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (bson.ObjectID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return secretKey, nil
	})

	if err != nil {
		log.Println(err)
		return bson.ObjectID{}, err
	}

	if !token.Valid {
		err = errors.New("invalid jwt token")
		return bson.ObjectID{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("invalid claims")
		return bson.ObjectID{}, err
	}

	usrId, err := bson.ObjectIDFromHex(claims["user_id"].(string))
	if err != nil {
		log.Println(err)
		return bson.ObjectID{}, err
	}

	return usrId, nil
}

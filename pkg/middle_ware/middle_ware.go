package middle_ware

import (
	"github.com/dgrijalva/jwt-go"
	"rentServer/pkg/config"
	"time"
)

//var jwtSecret = []byte(config.GetConfig().JwtSecret)

type Claims struct {
	Id uint64 `json:"id"`
	jwt.StandardClaims
}

type AdminClaims struct {
	ID uint64 `json:"id"`
	jwt.StandardClaims
}

// ParseToken 验证用户token
func ParseToken(token string) (*Claims, error) {
	var jwtSecret = []byte(config.GetConfig().JwtSecret)
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (i interface{}, e error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

// GenerateToken 签发用户token
func GenerateToken(id uint64) (string, error) {
	var jwtSecret = []byte(config.GetConfig().JwtSecret)
	nowTime := time.Now()
	expireTime := nowTime.Add(720 * time.Hour)
	claims := Claims{
		Id: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "sys",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// ParseAdminToken 验证管理元token
func ParseAdminToken(token string) (*AdminClaims, error) {
	var jwtSecret = []byte(config.GetConfig().JwtSecret)
	tokenClaims, err := jwt.ParseWithClaims(token, &AdminClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*AdminClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

// GenerateAdminToken 签发管理员token
func GenerateAdminToken(id uint64) (string, error) {
	var jwtSecret = []byte(config.GetConfig().JwtSecret)
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour)
	claims := AdminClaims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "admin",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

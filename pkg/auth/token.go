package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/ramseyjiang/go_senior_to_principle/pkg/utils"
)

type TokenDetail struct {
	AccessToken     string
	RefreshToken    string
	AccessUUID      string
	RefreshUUID     string
	AccessExpireAt  int64
	RefreshExpireAt int64
}

type AccessDetail struct {
	AccessUUID string
	UserID     uint64
}

// GenerateToken is used to generate accessToken and refreshToken in the TokenDetail.
// The accessToken is expired at 1 hour later.
// The refreshToken is expired at 7 days later.
// The uuid also is used at here. Because uuid is unique each time.
// Nowadays, people always have several devices, using uuid, it is distinctive different devices.
func GenerateToken(rds *redis.Client, userID uint64) (*TokenDetail, error) {
	td := &TokenDetail{}
	td.AccessExpireAt = time.Now().Add(utils.AccessTokenExpiredAt).Unix()
	td.AccessUUID = uuid.New().String()
	td.RefreshExpireAt = time.Now().Add(utils.RefreshTokenExpiredAt).Unix()
	td.RefreshUUID = uuid.New().String()

	// MapClaims is a claims type that uses the map[string]interface{}. This is the default claims type if you don't supply one.
	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims["authorized"] = true
	accessTokenClaims["access_uuid"] = td.AccessUUID
	accessTokenClaims["user_id"] = userID
	accessTokenClaims["iat"] = time.Now().Unix() // iat means Issued At
	accessTokenClaims["exp"] = td.AccessExpireAt

	// jwt.NewWithClaims creates a new Token with the specified signing method and claims
	// Notice, the sign method must be the same when you use it to parse a token. As the VerifyToken method line 59.
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	var err error
	// td.AccessToken is a SignedString that creates and returns a complete, signed JWT
	td.AccessToken, err = accessToken.SignedString([]byte(os.Getenv("JWT_API_ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims["refresh_uuid"] = td.RefreshUUID
	refreshTokenClaims["user_id"] = userID
	refreshTokenClaims["iat"] = time.Now().Unix() // iat means Issued At
	refreshTokenClaims["exp"] = td.RefreshExpireAt
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	td.RefreshToken, err = refreshToken.SignedString([]byte(os.Getenv("JWT_API_REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	err = storeUUIDRedis(rds, userID, td)
	if err != nil {
		return nil, err
	}

	return td, nil
}

// storeUUIDRedis is used to save AccessUUID and RefreshUUID data, and their expired time in redis together.
func storeUUIDRedis(rds *redis.Client, userID uint64, td *TokenDetail) (err error) {
	ctx := context.Background()
	err = rds.Set(ctx, td.AccessUUID, userID, utils.AccessTokenExpiredAt).Err()
	if err != nil {
		return
	}

	err = rds.Set(ctx, td.RefreshUUID, userID, utils.RefreshTokenExpiredAt).Err()
	if err != nil {
		return
	}
	return
}

// ExtractTokenData is used to assign the data which will be stored in redis
func ExtractTokenData(c *gin.Context) (*AccessDetail, error) {
	token, err := ParseToken(c)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		UserID, uErr := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if uErr != nil {
			return nil, err
		}
		return &AccessDetail{
			AccessUUID: accessUUID,
			UserID:     UserID,
		}, nil
	}
	return nil, err
}

// FetchUIDFromRedis is used to get userID from redis, if it does not found. It means the token has expired.
func FetchUIDFromRedis(ad *AccessDetail, rds *redis.Client) (uint64, error) {
	ctx := context.Background()
	uid, err := rds.Get(ctx, ad.AccessUUID).Result()
	if err != nil || err != redis.Nil {
		return 0, err
	}

	userID, _ := strconv.ParseUint(uid, 10, 64)
	return userID, nil
}

// ParseToken is called ExtractToken inside, it will get the token string, then parse the token to using the signing method.
func ParseToken(c *gin.Context) (*jwt.Token, error) {
	tokenStr := extractToken(c)

	// Parse method verifies the signature and returns the parsed token.
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Notice here, the sign method must be the same when you use it to generate a token. As the GenerateToken method line 28.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_API_ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	return token, err
}

func ParseRefreshToken(refreshToken string) (*jwt.Token, error) {
	// Parse method verifies the signature and returns the parsed token.
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Notice here, the sign method must be the same when you use it to generate a token. As the GenerateToken method line 28.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_API_REFRESH_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	return token, err
}

// extractToken is used to get token string from header
func extractToken(c *gin.Context) string {
	// Bearer tokens come from the header in the format bearer <JWT>, split them to return a JWT string.
	bearerToken := c.Request.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}

func DelAuth(delUUID string, rds *redis.Client) (err error) {
	ctx := context.Background()
	intCmd := rds.Del(ctx, delUUID)
	if intCmd.Err() != nil {
		return intCmd.Err()
	}

	return nil
}

func RegenerateToken(token *jwt.Token, rds *redis.Client) (map[string]string, error) {
	claims, ok := token.Claims.(jwt.MapClaims) // the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUUID, _ := claims["refresh_uuid"].(string) // convert the interface to string
		userID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}

		// delete the previous refresh token base on the uuid
		if err := DelAuth(refreshUUID, rds); err != nil {
			return nil, err
		}

		// create new pairs of refresh and access tokens
		td, createErr := GenerateToken(rds, userID)
		if createErr != nil {
			return nil, createErr
		}

		// save uuids from tokens into redis
		saveErr := storeUUIDRedis(rds, userID, td)
		if saveErr != nil {
			return nil, saveErr
		}

		return map[string]string{
			"access_token":  td.AccessToken,
			"refresh_token": td.RefreshToken,
		}, nil
	} else {
		return nil, errors.New("refresh token expired")
	}
}

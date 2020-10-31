package redisjwt

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/twinj/uuid"
)

// TokenParams ...
type TokenParams struct {
	AccessSecret  string
	RefreshSecret string
	AccessTime    time.Duration
	RefreshTime   time.Duration
}

// TokenDetails ...
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

// AccessDetails ...
type AccessDetails struct {
	AccessUUID string
	UserID     string
	Role       string
}

type redisAuth struct {
	client        *redis.Client
	accessSecret  string
	refreshSecret string
}

func (r *redisAuth) CreateAuth(userid uint64, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := r.client.Set(td.AccessUUID, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := r.client.Set(td.RefreshUUID, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func (r *redisAuth) FetchAuth(authD *AccessDetails) (string, error) {
	userid, err := r.client.Get(authD.AccessUUID).Result()
	if err != nil {
		return "", err
	}
	return userid, nil
}

func (r *redisAuth) DeleteAuth(givenUUID string) (int64, error) {
	deleted, err := r.client.Del(givenUUID).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

// CreateToken ...
func CreateToken(userid, role string, params *TokenParams) (*TokenDetails, error) {

	td := &TokenDetails{}

	td.AtExpires = time.Now().Add(params.AccessTime).Unix() // ex: time.Minute * 15
	td.AccessUUID = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(params.RefreshTime).Unix() // ex: time.Hour * 24 * 7
	td.RefreshUUID = uuid.NewV4().String()

	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["role"] = role
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(params.AccessSecret))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(params.RefreshSecret))
	if err != nil {
		return nil, err
	}
	return td, nil
}

// ExtractTokenMetadata ...
func ExtractTokenMetadata(r *http.Request, accessSecret string) (*AccessDetails, error) {
	token, err := VerifyToken(r, accessSecret)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userID, ok := claims["user_id"].(string)
		if !ok {
			return nil, err
		}
		role, ok := claims["role"].(string)
		if !ok {
			return nil, err
		}
		return &AccessDetails{
			AccessUUID: accessUUID,
			UserID:     userID,
			Role:       role,
		}, nil
	}
	return nil, err
}

// VerifyToken ...
func VerifyToken(r *http.Request, accessSecret string) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(accessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// ExtractToken ...
func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

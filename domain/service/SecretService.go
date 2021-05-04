package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/smarest/smarest-common/domain/entity/exception"
	"github.com/smarest/smarest-portal/domain/entity"
)

// SecretService get restaurant info
type SecretService struct {
	JWTKey      []byte
	ExpiredTime time.Duration
}

// Claims store restaurantID and jwt keys
type RestaurantClaims struct {
	Restaurant *entity.Restaurant
	jwt.StandardClaims
}

// NewSecretService is SecretService's constructor
func NewSecretService(key string) *SecretService {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return &SecretService{
		JWTKey:      []byte(hex.EncodeToString(hasher.Sum(nil))),
		ExpiredTime: 60 * 24 * time.Minute,
	}
}

// GenerateToken create token
func (s *SecretService) GenerateRestaurantToken(restaurant *entity.Restaurant) (string, error) {
	expirationTime := time.Now().Add(s.ExpiredTime)
	claims := &RestaurantClaims{
		Restaurant: restaurant,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	return token.SignedString(s.JWTKey)
}

// CheckToken check token
func (s *SecretService) CheckRestaurantToken(cookie string) (*RestaurantClaims, *exception.Error) {
	// Initialize a new instance of `Claims`
	claims := &RestaurantClaims{}
	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (interface{}, error) {
		return s.JWTKey, nil
	})

	if err != nil {
		return nil, exception.CreateError(exception.CodeSignatureInvalid, "Cookie is invalid")
	}
	if !tkn.Valid {
		return nil, exception.CreateError(exception.CodeSignatureInvalid, "Cookie is invalid")
	}
	return claims, nil
}

func (s *SecretService) EncryptString(value string) (string, *exception.Error) {
	block, _ := aes.NewCipher(s.JWTKey)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", exception.CreateError(exception.CodeUnknown, err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", exception.CreateError(exception.CodeUnknown, err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(value), nil)
	return hex.EncodeToString(ciphertext), nil
}

func (s *SecretService) DecryptString(value string) (string, *exception.Error) {
	block, err := aes.NewCipher(s.JWTKey)
	if err != nil {
		return "", exception.CreateError(exception.CodeUnknown, err.Error())
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", exception.CreateError(exception.CodeUnknown, err.Error())
	}

	data, err := hex.DecodeString(value)
	if err != nil {
		return "", exception.CreateError(exception.CodeUnknown, err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", exception.CreateError(exception.CodeUnknown, err.Error())
	}
	return string(plaintext), nil
}

func (s *SecretService) AddHashFieldToSlice(fieldName string, item interface{}) {
	item.(map[string]interface{})[fieldName+"Hash"], _ = s.EncryptString(fmt.Sprint(item.(map[string]interface{})[fieldName].(float64)))
}
func (s *SecretService) AddHashFieldToSlices(fieldName string, items []interface{}) {
	for _, item := range items {
		s.AddHashFieldToSlice(fieldName, item)
	}
}

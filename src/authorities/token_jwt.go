package authorities

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"
)
const (
	TokenIssuer = "Do-not-involution@36625090"
)

type Claims struct {
	jwt.RegisteredClaims
	AccountRoles []string
	Principal    []byte
}

type jwtTokenHandler struct {
	pri *rsa.PrivateKey
	pub *rsa.PublicKey
	settings *Settings
}

func NewJwtTokenHandler(settings *Settings) (TokenHandler, error) {
	h := &jwtTokenHandler{
		pri: nil,
		pub:  nil,
		settings:   settings,
	}

	block,err := base64.StdEncoding.DecodeString(settings.PKCS1PublicKey)
	if err != nil {
		return nil, fmt.Errorf("decoding public key: %v", err)
	}
	h.pub, err = x509.ParsePKCS1PublicKey(block)
	if err != nil {
		return nil, fmt.Errorf("parser public key: %v", err)
	}

	block, err = base64.StdEncoding.DecodeString(settings.PKCS8PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("decoding private key: %v", err)
	}
	pri, err := x509.ParsePKCS8PrivateKey(block)
	if err != nil {
		return nil, fmt.Errorf("parser PKCS8 private key: %v", err)
	}
	var ok bool
	h.pri, ok = pri.(*rsa.PrivateKey)
	if !ok{
		return nil, fmt.Errorf("invalid PKCS8 private key")
	}

	return h, nil
}


// GenerateToken 产生token的函数
// 返回 Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9................................
func (m *jwtTokenHandler)GenerateToken(auth *Authorized) (string, error) {

	ciphertext, err := json.Marshal(auth.GetPrincipal())
	if nil != err{
		return "", err
	}

	enc, err := rsa.EncryptPKCS1v15(rand.Reader, m.pub, ciphertext)
	if err != nil {
		return "", fmt.Errorf("encrypt principal: %v", err)
	}

	claims := Claims{}
	claims.ID = auth.ID
	claims.Audience = jwt.ClaimStrings{auth.Account}
	claims.Issuer = TokenIssuer
	claims.AccountRoles = auth.AccountRoles
	claims.Principal = enc

	if m.settings.Timeout > 0{
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(m.settings.Timeout * time.Second))
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	token, err := tokenClaims.SignedString(m.pri)
	if err != nil {
		return "", fmt.Errorf("jwt signing failed: %v", err)
	}

	return "Bearer " + token, nil
}

// ParseToken
// 验证token的函数 Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9................................
func (m *jwtTokenHandler)parseToken(token string) (*Claims, error) {

	fields := strings.Split(token, " ")
	if len(fields) != 2 || fields[0] != "Bearer"{
		return nil, errors.New("invalid token, value must be: 'Bearer ......'")
	}

	tokenClaims, err := jwt.ParseWithClaims(fields[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return m.pub, nil
	})

	if err != nil {
		return nil, fmt.Errorf("token key with claims: %v", err)
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}else {
			return nil, tokenClaims.Claims.Valid()
		}
	}

	return nil, err
}


//ParseToken 解密验证信息
func (m *jwtTokenHandler)ParseToken(token string) (*Authorized, error) {
	claims, err := m.parseToken(token)
	if err != nil {
		return nil, err
	}

	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, m.pri, claims.Principal)
	if err != nil {
		return nil, err
	}

	var principal Principal
	if err := json.Unmarshal(plaintext, &principal); err != nil {
		return nil, err
	}

	authorized := &Authorized{
		ID:           claims.ID,
		AccountRoles: claims.AccountRoles,
		Principal:    principal,
	}
	if len(claims.Audience) >0{
		authorized.Account = claims.Audience[0]
	}
	return authorized, nil
}

package jwtservice

import (
	"errors"
	"strconv"
	"time"

	"github.com/dundunlabs/goauth/common"
	"github.com/dundunlabs/goauth/model"
	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	*common.BaseService
}

type Claims struct {
	jwt.RegisteredClaims
}

func (s *JWTService) Sign(session *model.Session) (string, error) {
	privKey, err := s.Config.JWT.PrivateKey()
	if err != nil {
		return "", err
	}
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        strconv.Itoa(int(session.ID)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: nil,
			Audience:  jwt.ClaimStrings{},
			Subject:   strconv.Itoa(int(session.Credential.IdentityID)),
			Issuer:    "",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privKey)
}

func (s *JWTService) Parse(str string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(str, &Claims{}, func(t *jwt.Token) (any, error) {
		return s.Config.JWT.PublicKey()
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid claimns")
}

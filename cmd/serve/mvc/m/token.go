package m

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	TokenTypeActivate      = TokenType(1)
	TokenTypeResetPassword = TokenType(2)
	TokenTypeMin           = TokenType(1)
	TokenTypeMax           = TokenType(2)

	TokenStatusValid   = TokenStatus(1)
	TokenStatusInvalid = TokenStatus(2)
	TokenStatusMin     = TokenStatus(1)
	TokenStatusMax     = TokenStatus(2)

	TokenSize = 8
)

type TokenType uint
type TokenStatus uint

type Token struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Duration  int64
	Code      string `query:"code" form:"code"`
	Type      TokenType
	Status    TokenStatus
	UserID    uint
}

// Load
func (token *Token) Load() (err error) {

	db, err := GetDB()
	if err != nil {
		return
	}
	defer db.Close()

	if db.Where("code = ?", token.Code).First(token).RecordNotFound() {
		err = fmt.Errorf("token not found")
	}
	return
}

// UsedAsActivate
func (token *Token) UsedAs(t TokenType) (err error) {
	if token.Type != t {
		return fmt.Errorf("invalid token type")
	}
	return token.Use()
}

// NewToken
func NewToken(t TokenType, user *User) (token *Token, err error) {

	// check user
	if user == nil {
		err = fmt.Errorf("invalid user")
		return
	}

	token = &Token{}
	token.UserID = user.ID

	if t < TokenTypeMin || t > TokenTypeMax {
		err = fmt.Errorf("invalid token type")
		return
	}
	token.Type = t
	token.Status = TokenStatusValid

	token.Duration = int64(24 * 60 * 60 * 2)
	token.Code = randString(TokenSize)

	db, err := GetDB()
	if err != nil {
		return
	}
	defer db.Close()

	tt := &Token{}
	if !db.Where("code = ?", token.Code).First(tt).RecordNotFound() {
		err = fmt.Errorf("create token value failed")
		return
	}

	db.Create(token)
	return
}

// Use
func (token *Token) Use() (err error) {
	if token.Status == TokenStatusInvalid {
		return fmt.Errorf("invalid token")
	}

	now := time.Now()
	expiredAt := token.CreatedAt.Add(time.Second * time.Duration(token.Duration))
	if now.After(expiredAt) {
		return fmt.Errorf("expired token")
	}

	db, err := GetDB()
	if err != nil {
		return
	}
	defer db.Close()

	token.Status = TokenStatusInvalid
	db.Save(token)
	return
}

// Owner
func (token *Token) Owner() (user *User, err error) {
	db, err := GetDB()
	if err != nil {
		return
	}
	defer db.Close()

	user = &User{}
	if db.Where("id = ?", token.UserID).First(user).RecordNotFound() {
		err = fmt.Errorf("no user found")
	}

	return
}

// RandString
func randString(size int) string {
	var src = rand.NewSource(time.Now().UnixNano())
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	return func(n int) string {
		b := make([]byte, n)
		// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
		for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
			if remain == 0 {
				cache, remain = src.Int63(), letterIdxMax
			}
			if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
				b[i] = letterBytes[idx]
				i--
			}
			cache >>= letterIdxBits
			remain--
		}

		return string(b)
	}(size)
}

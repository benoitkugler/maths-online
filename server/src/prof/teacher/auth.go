package teacher

import (
	"encoding/json"
	"fmt"
	"time"

	tc "github.com/benoitkugler/maths-online/sql/teacher"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	DeltaTokenJours = 3
	deltaToken      = DeltaTokenJours * 24 * time.Hour
)

// UserMeta are custom claims extending default ones.
type UserMeta struct {
	Teacher tc.Teacher
	jwt.StandardClaims
}

func (ct *Controller) JWTMiddleware() echo.MiddlewareFunc {
	config := middleware.JWTConfig{SigningKey: ct.teacherKey[:], Claims: &UserMeta{}}
	return middleware.JWTWithConfig(config)
}

// expects the token to be in the `token` query parameters
func (ct *Controller) JWTMiddlewareForQuery() echo.MiddlewareFunc {
	config := middleware.JWTConfig{SigningKey: ct.teacherKey[:], Claims: &UserMeta{}, TokenLookup: "query:token"}
	return middleware.JWTWithConfig(config)
}

func (ct *Controller) newToken(teacher tc.Teacher) (string, error) {
	// Set custom claims
	claims := &UserMeta{
		Teacher: teacher,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(deltaToken).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	return token.SignedString(ct.teacherKey[:])
}

// JWTTeacher expects a JWT authentified request, and must
// only be used in routes protected by `JWTMiddleware`
func JWTTeacher(c echo.Context) tc.Teacher {
	meta := c.Get("user").(*jwt.Token).Claims.(*UserMeta) // the token is valid here
	return meta.Teacher
}

// GetDevToken creates a new user and returns a valid token,
// so that client frontend doesn't have to use password when developping.
func (ct *Controller) GetDevToken() (string, error) {
	t, err := tc.Teacher{
		Mail:            fmt.Sprintf("%d", time.Now().Unix()),
		PasswordCrypted: ct.teacherKey.EncryptPassword("1234"),
	}.Insert(ct.db)
	if err != nil {
		return "", err
	}
	token, err := ct.newToken(t)
	if err != nil {
		return "", err
	}
	type meta struct {
		IdTeacher tc.IdTeacher
		Token     string
	}
	out, err := json.Marshal(meta{IdTeacher: t.Id, Token: token})
	return string(out), err
}

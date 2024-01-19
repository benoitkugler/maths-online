package teacher

import (
	"encoding/json"
	"time"

	tc "github.com/benoitkugler/maths-online/server/src/sql/teacher"
	"github.com/benoitkugler/maths-online/server/src/utils"
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
	IdTeacher tc.IdTeacher
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

func (ct *Controller) newToken(id tc.IdTeacher) (string, error) {
	// Set custom claims
	claims := &UserMeta{
		IdTeacher: id,
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
func JWTTeacher(c echo.Context) tc.IdTeacher {
	meta := c.Get("user").(*jwt.Token).Claims.(*UserMeta) // the token is valid here
	return meta.IdTeacher
}

// GetDevTokens creates a new user and returns a valid token,
// so that client frontend doesn't have to use password when developping.
func (ct *Controller) GetDevTokens() (string, string, error) {
	mail := utils.RandomString(false, 8) + "@dummy.com"
	t, err := tc.Teacher{
		Mail:            mail,
		PasswordCrypted: ct.teacherKey.EncryptPassword("1234"),
		FavoriteMatiere: tc.Mathematiques,
	}.Insert(ct.db)
	if err != nil {
		return "", "", err
	}
	type meta struct {
		IdTeacher tc.IdTeacher
		Token     string
	}
	token1, err := ct.newToken(t.Id)
	if err != nil {
		return "", "", err
	}
	out1, err := json.Marshal(meta{IdTeacher: t.Id, Token: token1})
	if err != nil {
		return "", "", err
	}

	token2, err := ct.newToken(ct.admin.Id)
	if err != nil {
		return "", "", err
	}
	out2, err := json.Marshal(meta{IdTeacher: ct.admin.Id, Token: token2})
	if err != nil {
		return "", "", err
	}
	return string(out1), string(out2), err
}

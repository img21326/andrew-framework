package helper

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
)

var SessionHelperInstance *SessionHelper

type SessionHelper struct {
	session *sessions.CookieStore
}

func newSessionHelper() *SessionHelper {
	serverKey := viper.GetViper().GetString("SERVER_KEY")
	return &SessionHelper{
		session: sessions.NewCookieStore([]byte(serverKey)),
	}
}

func GetSessionHelper() *SessionHelper {
	if SessionHelperInstance == nil {
		SessionHelperInstance = newSessionHelper()
	}
	return SessionHelperInstance
}

func (s *SessionHelper) UserLogin(ctx *gin.Context, userModel interface{}) {
	session, _ := s.session.Get(ctx.Request, "session")
	session.Values["user"] = userModel
	session.Options.MaxAge = 60 * 60 * 24 * 7 // 7 days
	session.Save(ctx.Request, ctx.Writer)
}

func (s *SessionHelper) UserLogout(ctx *gin.Context) {
	session, _ := s.session.Get(ctx.Request, "session")
	delete(session.Values, "user")
	session.Save(ctx.Request, ctx.Writer)
}

func (s *SessionHelper) GetCurrentUser(ctx *gin.Context) interface{} {
	session, _ := s.session.Get(ctx.Request, "session")
	if user, ok := session.Values["user"]; ok {
		return user
	}
	return nil
}

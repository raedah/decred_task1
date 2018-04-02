package controller

import (
	"gopkg.in/kataras/iris.v8"
	"github.com/kataras/iris/sessions"
	"time"
)

type (
	User struct {
		ID 				int			`storm:"increment"`
		Email			string		`validate:"required,email" storm:"unique"`
		Password 		string		`validate:"required"`
		Name			string		`validate:"required"`
		Telegram		string
		Skype			string
		WhatsApp		string
		Signal			string
		CreatedDate		time.Time
		ReferralId		int
	}
)


var (
	cookieNameForSessionID 	= "decred_task_1"
	sess                   	= sessions.New(sessions.Config{Cookie: cookieNameForSessionID})
	userStoreKey			= "user"
)

func BindRoute(app *iris.Application)  {
	app.UseGlobal(middleware)
	app.Get("/", getHome)
	app.Get("/register", getRegister)
	app.Post("/register", postRegister)
	app.Get("/login", getLogin)
	app.Post("/login", postLogin)
	app.Get("/logout", func(ctx iris.Context) {
		session := sess.Start(ctx)
		session.Delete("user")
		ctx.Redirect("/")
	})

	app.Get("/profile", needLogin, getProfile)
	app.Put("/profile", needLogin, updateProfile)
	app.Get("/change-password", needLogin, getChangePassword)
	app.Put("/change-password", needLogin, putChangePassword)
}

package ctrl

import (
	"gopkg.in/kataras/iris.v8"
	"github.com/kataras/iris/sessions"
	"time"
)

type (
	User struct {
		ID 				int			`storm:"increment"`
		UserName		string		`validate:"required,excludes= " storm:"unique"`
		Email			string		`validate:"omitempty,email"`
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
		ctx.IsAjax()
		session := sess.Start(ctx)
		session.Delete("user")
		ctx.Redirect("/")
	})

	app.Get("/profile", needLogin, getProfile)
	app.Put("/profile", needLogin, updateProfile)
	app.Get("/change-password", needLogin, getChangePassword)
	app.Put("/change-password", needLogin, putChangePassword)
}

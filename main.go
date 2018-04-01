package main

import (
	"gopkg.in/kataras/iris.v8"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
	"strings"
	"github.com/asdine/storm"
	"golang.org/x/crypto/bcrypt"
	"github.com/kataras/iris/sessions"
)

type (
	User struct {
		ID 				int			`storm:"increment"`
		Email			string		`validate:"required,email" storm:"unique"`
		Password 		string		`validate:"required"`
		Name			string		`validate:"required"`
	}
)

var (
	cookieNameForSessionID = "mycookiesessionnameid"
	sess                   = sessions.New(sessions.Config{Cookie: cookieNameForSessionID})
)

func main() {
	app := iris.Default()
	/* register view */
	app.RegisterView(iris.HTML("./templates", ".html").Layout("layout.html"))
	/* use global middleware */
	app.UseGlobal(middleware)
	/* serve static file */
	app.StaticWeb("/public", "./public")
	/* bind route */
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
	/* start server */
	app.Run(iris.Addr(":8080"),
		iris.WithCharset("UTF-8"),
			iris.WithoutVersionChecker)
}

func getHome(ctx iris.Context)  {
	ctx.ViewData("title", "Home")
	ctx.View("index.html")
}

func getLogin(ctx iris.Context)  {
	ctx.ViewData("title", "Login")
	ctx.View("login.html")
}

func getRegister(ctx iris.Context)  {
	ctx.ViewData("title", "Register")
	ctx.View("register.html")
}

func postLogin(ctx iris.Context)  {
	user := User{}
	if err := ctx.ReadForm(&user);err == nil {

		if err,errDetail := validateStruct(&user,"Name");err == nil {
			db, err := storm.Open("user.db")
			defer db.Close()
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.JSON(iris.Map{
					"error": err.Error(),
				})
			}else{
				storeUser := User{}
				if err := db.One("Email", user.Email, &storeUser);err == nil {
					if err := bcrypt.CompareHashAndPassword([]byte(storeUser.Password),[]byte(user.Password));err == nil {
						session := sess.Start(ctx)
						session.Set("user", storeUser)
						ctx.JSON(storeUser)
					}else{
						ctx.StatusCode(iris.StatusUnauthorized)
						ctx.JSON(iris.Map{
							"error": "Invalid password",
						})
					}
				}else{
					ctx.StatusCode(iris.StatusNotFound)
					ctx.JSON(iris.Map{
						"error": "Could not find user",
					})
				}
			}
		}else{
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{
				"error": errDetail,
			})
		}
	}else {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"error": err.Error(),
		})
	}
}

func postRegister(ctx iris.Context)  {
	user := User{}
	if err := ctx.ReadForm(&user);err == nil {

		if err,errDetail := validateStruct(&user);err == nil {
			db, err := storm.Open("user.db")
			defer db.Close()
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.JSON(iris.Map{
					"error": err.Error(),
				})
			}else{
				newPass,_ := bcrypt.GenerateFromPassword([]byte(user.Password),11)
				user.Password = string(newPass)
				if err := db.Save(&user);err != nil {
					if err == storm.ErrAlreadyExists {
						ctx.StatusCode(iris.StatusConflict)
						ctx.JSON(iris.Map{
							"error": iris.Map{
								"Email": "The email is already taken",
							},
						})
					}else{
						ctx.StatusCode(iris.StatusInternalServerError)
						ctx.JSON(iris.Map{
							"error": err.Error(),
						})
					}
				}else{
					ctx.StatusCode(200)
					ctx.JSON(user)
				}
			}
		}else{
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{
				"error": errDetail,
			})
		}
	}else {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"error": err.Error(),
		})
	}

}

func middleware(ctx iris.Context)  {
	session := sess.Start(ctx)
	if user,err := session.Get("user").(User);err == true {
		ctx.ViewData("authenticated",true)
		ctx.ViewData("user", user)
	}else {
		ctx.ViewData("user", nil)
		ctx.ViewData("authenticated",false)
	}
	ctx.ViewData("title", "test site")

	ctx.Next()
}

func validateStruct(target interface{},ignoreFields ...string) (error,validator.ValidationErrorsTranslations) {
	validate := validator.New()
	enLang := en.New()
	uni := ut.New(enLang)

	trans, _ := uni.GetTranslator("en")

	en_translations.RegisterDefaultTranslations(validate, trans)

	err := validate.StructExcept(target,ignoreFields...)

	if err != nil {
		errs := err.(validator.ValidationErrors)
		tran := errs.Translate(trans)
		return err,convertValidation(tran)
	}
	return nil,nil
}

func convertValidation(validMess validator.ValidationErrorsTranslations) validator.ValidationErrorsTranslations {
	if(validMess == nil){
		return validMess
	}
	newMess := make(validator.ValidationErrorsTranslations)
	for k,mess := range validMess{
		if index := strings.Index(k,".");index != -1{
			k = k[index+1:]
		}
		newMess[k] = mess
	}

	return newMess
}
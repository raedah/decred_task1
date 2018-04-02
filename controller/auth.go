package controller

import (
	"gopkg.in/kataras/iris.v8"
	"github.com/asdine/storm"
	"golang.org/x/crypto/bcrypt"
	"time"
	"strconv"
)

func getHome(ctx iris.Context)  {
	session := sess.Start(ctx)
	if session.Get("ref") == nil {
		if ref,err := strconv.Atoi(ctx.FormValue("ref"));err == nil {
			db, err := storm.Open("user.db")
			defer db.Close()
			if err == nil {
				user := User{}
				if err := db.One("ID", ref, &user);err == nil {
					session.Set("ref",user.ID)
				}
			}
		}
	}

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
						session.Set(userStoreKey, storeUser)
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
				user.CreatedDate = time.Now()
				session := sess.Start(ctx)
				if ref,ok := session.Get("ref").(int); ok {
					user.ReferralId = ref
				}
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
					ctx.JSON(user.ID)
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
	if user,ok := session.Get(userStoreKey).(User);ok == true {
		ctx.ViewData("authenticated",true)
		ctx.ViewData("user", user)
	}else {
		ctx.ViewData("user", nil)
		ctx.ViewData("authenticated",false)
	}
	ctx.ViewData("title", "test site")

	ctx.Next()
}

func getProfile(ctx iris.Context)  {
	ctx.ViewData("host", "http://"+ctx.Host())
	ctx.ViewData("title", "Profile")
	ctx.View("profile.html")
}

func updateProfile(ctx iris.Context)  {
	user := User{}
	if err := ctx.ReadForm(&user);err == nil {

		if err,errDetail := validateStruct(&user,"Password");err == nil {
			db, err := storm.Open("user.db")
			defer db.Close()
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.JSON(iris.Map{
					"error": err.Error(),
				})
			}else{
				session := sess.Start(ctx)
				oldUser,_ := session.Get(userStoreKey).(User)
				user.ID = oldUser.ID
				user.Password = oldUser.Password
				user.CreatedDate = oldUser.CreatedDate
				if err := db.Update(&user);err != nil {
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
					ctx.JSON(user.ID)
					session.Set(userStoreKey,user)
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

func getChangePassword(ctx iris.Context)  {
	ctx.ViewData("title", "change-password")
	ctx.View("change-password.html")
}
func putChangePassword(ctx iris.Context)  {
	session := sess.Start(ctx)
	user,_ := session.Get(userStoreKey).(User)
	password := struct {
		OldPassword		string
		NewPassword		string
	}{}

	if err := ctx.ReadForm(&password);err == nil {
		db, err := storm.Open("user.db")
		defer db.Close()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{
				"error": err.Error(),
			})
		}else{
			db.One("ID", user.ID, &user)
			if err := bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(password.OldPassword));err == nil {
				newPass,_ := bcrypt.GenerateFromPassword([]byte(password.NewPassword),11)
				user.Password = string(newPass)
				if err := db.UpdateField(&user, "Password", user.Password);err == nil {
					session.Set(userStoreKey,user)
					ctx.JSON(user.ID)
				}else{
					ctx.JSON(iris.Map{
						"error": err.Error(),
					})
				}
			}else{
				ctx.StatusCode(iris.StatusUnauthorized)
				ctx.JSON(iris.Map{
					"error": "Old Password is not match",
				})
			}
		}

	}else{
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"error": err.Error(),
		})
	}
}
func needLogin(ctx iris.Context)  {
	session := sess.Start(ctx)
	if _,ok := session.Get(userStoreKey).(User);ok != true {
		if ctx.IsAjax() {
			ctx.JSON(iris.Map{
				"error": "You need login to do this function",
			})
		}else{
			ctx.Redirect("/")
		}
		return
	}
	ctx.Next()
}
package main

import (
	"gopkg.in/kataras/iris.v8"
	"github.com/nobita0590/decred_task1/ctrl"
)

func main() {
	app := iris.Default()
	/* register view */
	app.RegisterView(iris.HTML("./templates", ".html").Layout("layout.html"))
	/* serve static file */
	app.StaticWeb("/public", "./public")
	/* bind route */
	ctrl.BindRoute(app)
	/* start server */
	app.Run(iris.Addr(":8080"),
		iris.WithCharset("UTF-8"),
			iris.WithoutVersionChecker)
}

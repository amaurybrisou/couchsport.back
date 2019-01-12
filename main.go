package main

import (
	"couchsport/api/handlers"
	"couchsport/api/stores"
	"couchsport/config"
	"couchsport/server"
	"flag"
)

func main() {
	env := flag.String("env", "dev", "select environment config file to use (will load config.[env].json")
	flag.Parse()

	c := config.Load(*env)

	srv := server.NewInstance(c)

	storeFactory := stores.NewStoreFactory(srv.Db, *c)
	storeFactory.Init()
	handlerFactory := handlers.NewHandlerFactory(storeFactory)

	srv.RegisterHandler("/languages", handlerFactory.LanguageHandler().All)
	srv.RegisterHandler("/activities", handlerFactory.ActivityHandler().All)
	srv.RegisterHandler("/pages", handlerFactory.PageHandler().All)

	srv.RegisterHandler("/pages/new", handlerFactory.UserHandler().IsLogged(
		handlerFactory.PageHandler().New),
	)
	srv.RegisterHandler("/pages/update", handlerFactory.UserHandler().IsLogged(
		handlerFactory.PageHandler().Update),
	)
	srv.RegisterHandler("/pages/publish", handlerFactory.UserHandler().IsLogged(
		handlerFactory.PageHandler().Publish),
	)
	srv.RegisterHandler("/pages/delete", handlerFactory.UserHandler().IsLogged(
		handlerFactory.PageHandler().Delete),
	)

	srv.RegisterHandler("/images/delete", handlerFactory.UserHandler().IsLogged(
		handlerFactory.ImageHandler().SoftDelete),
	)

	// srv.RegisterHandler("/users", handlerFactory.UserHandler().All)

	srv.RegisterHandler("/profiles/update", handlerFactory.UserHandler().IsLogged(
		handlerFactory.ProfileHandler().Update),
	)
	srv.RegisterHandler("/profiles/mine", handlerFactory.UserHandler().IsLogged(
		handlerFactory.ProfileHandler().Mine),
	)
	srv.RegisterHandler("/profiles/pages", handlerFactory.UserHandler().IsLogged(
		handlerFactory.PageHandler().ProfilePages),
	)

	srv.RegisterHandler("/login", handlerFactory.UserHandler().Login)
	srv.RegisterHandler("/signin", handlerFactory.UserHandler().SignIn)
	srv.RegisterHandler("/logout", handlerFactory.UserHandler().IsLogged(
		handlerFactory.UserHandler().Logout),
	)

	srv.ServePublic(c.PublicPath)

	srv.Start()

}

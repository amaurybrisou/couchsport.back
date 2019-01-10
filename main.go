package main

import (
	"couchsport/api/handlers"
	"couchsport/api/handlers/image"
	"couchsport/api/handlers/invitation"
	"couchsport/api/handlers/page"
	"couchsport/api/handlers/profile"
	"couchsport/api/handlers/user"
	"couchsport/api/stores"
	"couchsport/api/types"
	"couchsport/config"
	"couchsport/server"
	"flag"
)

func main() {
	env := flag.String("env", "dev", "select environment config file to use (will load config.[env].json")
	flag.Parse()

	c := config.Load(*env)

	srv := server.NewInstance(c)

	activityStore := stores.ActivityStore{Db: srv.Db}
	activityStore.Migrate()

	languageStore := stores.LanguageStore{Db: srv.Db}
	languageStore.Migrate()

	invitationStore := stores.InvitationStore{Db: srv.Db}
	invitationStore.Migrate()

	fileStore := stores.FileStore{
		FileSystem:    types.OsFS{},
		PublicPath:    c.PublicPath,
		ImageBasePath: c.ImageBasePath,
		FilePrefix:    c.FilePrefix,
	}

	imageStore := stores.ImageStore{Db: srv.Db}
	imageStore.Migrate()

	userStore := stores.UserStore{Db: srv.Db}
	userStore.Migrate()

	profileStore := stores.ProfileStore{Db: srv.Db, FileStore: fileStore}
	profileStore.Migrate()

	pageStore := stores.PageStore{Db: srv.Db, FileStore: fileStore, ProfileStore: profileStore}
	pageStore.Migrate()

	sessionStore := stores.SessionStore{Db: srv.Db}
	sessionStore.Migrate()

	pageHandler := page.PageHandler{
		Store: pageStore,
	}

	languageHandler := handlers.LanguageHandler{
		Store: languageStore,
	}

	activityHandler := handlers.ActivityHandler{
		Store: activityStore,
	}

	userHandler := user.UserHandler{
		Store:        userStore,
		SessionStore: &sessionStore,
	}

	imageHandler := image.ImageHandler{
		Store: imageStore,
	}

	profileHandler := profile.ProfileHandler{
		Store:     profileStore,
		UserStore: userStore,
	}

	srv.RegisterHandler("/languages", languageHandler.GetLanguages)
	srv.RegisterHandler("/activities", activityHandler.GetActivities)
	srv.RegisterHandler("/invitations", invitation.InvitationHandler{
		Store: invitationStore,
	}.IndexHandler)

	srv.RegisterHandler("/pages", pageHandler.IndexHandler)

	srv.RegisterHandler("/pages/new", userHandler.IsLogged(
		pageHandler.CreateHandler),
	)
	srv.RegisterHandler("/pages/publish", userHandler.IsLogged(
		pageHandler.PublishHandler),
	)
	srv.RegisterHandler("/pages/delete", userHandler.IsLogged(
		pageHandler.DeleteHandler),
	)

	srv.RegisterHandler("/images/delete", userHandler.IsLogged(
		imageHandler.SoftDeleteHandler),
	)

	srv.RegisterHandler("/users", userHandler.IndexHandler)
	// srv.RegisterHandler("/images/upload", userHandler.IsLogged(imageHandler.UploadHandler))

	srv.RegisterHandler("/profiles/update", userHandler.IsLogged(
		profileHandler.UpdateProfile),
	)
	srv.RegisterHandler("/profiles/mine", userHandler.IsLogged(
		profileHandler.GetProfileHandler),
	)
	srv.RegisterHandler("/profiles/pages", userHandler.IsLogged(
		pageHandler.GetProfilePagesHandler),
	)

	srv.RegisterHandler("/login", userHandler.Login)
	srv.RegisterHandler("/signin", userHandler.SignIn)
	srv.RegisterHandler("/logout", userHandler.IsLogged(
		userHandler.Logout),
	)

	srv.ServePublic(c.PublicPath)

	srv.Start()

}

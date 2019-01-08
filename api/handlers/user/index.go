package user

import (
	"couchsport/api/models"
	"couchsport/api/stores"
	"couchsport/api/utils"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"strconv"
)

type UserHandler struct {
	Store        stores.UserStore
	SessionStore *stores.SessionStore
}

func (app UserHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {

	keys := r.URL.Query()

	users := app.Store.GetUsers(keys)

	json, err := json.Marshal(users)

	if err != nil {
		log.Error(err)
	}

	fmt.Fprintf(w, string(json))

}

func (app UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {

	r.Close = true

	if r.Body != nil {
		defer r.Body.Close()
	}

	formUser, err := parseBody(r.Body)
	if err != nil {
		http.Error(w, fmt.Errorf("could not parse body %s", err).Error(), http.StatusBadRequest)
		return
	}

	user, err := app.Store.FindOrCreate(formUser)

	if err != nil {
		http.Error(w, fmt.Errorf("could not fetch user %s", err).Error(), http.StatusInternalServerError)
	}

	json, err := json.Marshal(user)

	if err != nil {
		http.Error(w, fmt.Errorf("could not marshal object %s", err).Error(), http.StatusInternalServerError)
	}

	fmt.Fprintf(w, string(json))
}

func (app UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	r.Close = true

	if r.Body != nil {
		defer r.Body.Close()
	}

	formUser, err := parseBody(r.Body)
	if err != nil {
		http.Error(w, fmt.Errorf("could not extract session %s", err).Error(), http.StatusBadRequest)
		return
	}

	dbUser, err := app.Store.GetUser(models.User{Email: formUser.Email})
	if err != nil {
		http.Error(w, fmt.Errorf("could not fetch user %s", err).Error(), http.StatusBadRequest)
		return
	}

	if r := comparePasswords(dbUser.Password, []byte(formUser.Password)); !r {
		http.Error(w, fmt.Errorf("wrong credentials").Error(), http.StatusUnauthorized)
		return
	}

	isLogged, err := app.SessionStore.Create(dbUser.ID)
	if err != nil {
		// app.SessionStore.DestroyAllByUserID(dbUser.ID)
		http.Error(w, fmt.Errorf("could not create session %s", err).Error(), http.StatusInternalServerError)
		return
	}

	if isLogged == false {
		http.Error(w, fmt.Errorf("could not log in").Error(), http.StatusUnauthorized)
		return
	}

	cookie, err := app.SessionStore.CreateCookie()

	if err != nil {
		http.Error(w, fmt.Errorf("could not create cookie %s", err).Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, cookie)

	type res struct {
		Token string
		Email string
	}

	responseBody := res{Token: app.SessionStore.GetToken(), Email: dbUser.Email}

	json, err := json.Marshal(responseBody)

	if err != nil {
		http.Error(w, fmt.Errorf("could not marshal object %s", err).Error(), http.StatusInternalServerError)
	}

	fmt.Fprintf(w, string(json))
}

func (app UserHandler) IsLogged(pass func(userID uint, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		session, err := app.SessionStore.GetSession(r)
		if err != nil {
			http.Error(w, fmt.Errorf("could not extract session %s", err).Error(), http.StatusUnauthorized)
			return
		}

		if session.HasExpired() {
			app.SessionStore.Destroy(r)

			if err != nil {
				http.Error(w, fmt.Errorf("could not logout properly %s", err).Error(), http.StatusInternalServerError)
				return
			}

			http.Error(w, fmt.Errorf("session has expired").Error(), http.StatusUnauthorized)
			return
		}

		pass(session.UserID, w, r)
	}
}

func (app UserHandler) Logout(_ uint, w http.ResponseWriter, r *http.Request) {
	success, err := app.SessionStore.Destroy(r)
	if err != nil {
		http.Error(w, fmt.Errorf("could not logout properly %s", err).Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, `{ "Result" : `+strconv.FormatBool(success)+` }`)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return false
	}

	return true
}

func parseBody(body io.Reader) (*models.User, error) {
	tmp, err := utils.ParseBody(&models.User{}, body)
	if err != nil {
		return &models.User{}, err
	}

	r, ok := tmp.(*models.User)

	if !ok {
		return &models.User{}, err
	}

	return r, nil
}

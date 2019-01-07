package stores

import (
	"couchsport/api/models"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const tokenKey = "user-token"
const sessionValidity = 60 * 60

type SessionStore struct {
	Db     *gorm.DB
	token  string
	userID uint
}

func (app SessionStore) Migrate() {
	app.Db.AutoMigrate(&models.Session{})
}

func (app *SessionStore) Create(userID uint) (bool, error) {
	app.userID = userID

	out := models.Session{}
	if err := app.Db.Where("user_id = ?", userID).Where("(UNIX_TIMESTAMP(expires) - UNIX_TIMESTAMP()) > 0").First(&out).Error; gorm.IsRecordNotFoundError(err) {
		// record not found => remove all from user and create fresh session
		app.DestroyAllByUserID(userID)

		token, err := uuid.NewV4()
		if err != nil {
			return false, err
		}

		session := models.Session{
			SessionID: token.String(),
			UserID:    userID,
			Expires:   time.Now().Add(time.Duration(sessionValidity) * time.Second),
			Validity:  sessionValidity,
		}

		if errs := app.Db.Create(&session).GetErrors(); len(errs) > 0 {
			for err := range errs {
				log.Errorln(err)
			}

			return false, errs[0]
		}

		app.token = session.SessionID
		return true, nil
	}

	app.token = out.SessionID

	return true, nil
}

func (app *SessionStore) GetSession(r *http.Request) (*models.Session, error) {

	cookie, err := app.GetCookieFromRequest(r)
	if err != nil {
		return nil, err
	}

	if cookie.Value == "" {
		return nil, http.ErrNoCookie
	}

	var session = models.Session{}
	if errs := app.Db.Where("session_id = ?", cookie.Value).First(&session).GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorln(err)
		}
		return nil, errs[0]
	}

	app.token = session.SessionID
	app.userID = session.UserID

	return &session, nil
}

func (app *SessionStore) GetCookieFromRequest(r *http.Request) (*http.Cookie, error) {

	c, err := r.Cookie(tokenKey)
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status

			return nil, http.ErrNoCookie
		}
		// For any other type of error, return a bad request status
		return nil, http.ErrNoCookie
	}

	return c, nil

}

func (app *SessionStore) Destroy(r *http.Request) (bool, error) {
	if app.token == "" {
		return false, http.ErrNoCookie
	}

	if errs := app.Db.Where("user_id = ?", app.userID).Delete(&models.Session{}).GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorln(err)
		}
		return false, errs[0]
	}

	return true, nil
}

func (app *SessionStore) DestroyAllByUserID(userID uint) (bool, error) {
	if errs := app.Db.Where("user_id = ?", userID).Delete(&models.Session{}).GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorln(err)
		}
		return false, errs[0]
	}

	return true, nil
}

func (app *SessionStore) CreateCookie() (*http.Cookie, error) {
	if app.token == "" {
		return nil, fmt.Errorf("cannot generate cookie without token")
	}

	return &http.Cookie{
		Name:    tokenKey,
		Value:   app.token,
		Expires: time.Now().Add(sessionValidity * time.Second),
	}, nil
}

func (app SessionStore) GetToken() string {
	return app.token
}

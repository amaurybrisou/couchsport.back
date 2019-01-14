package stores

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/goland-amaurybrisou/couchsport/api/models"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const tokenKey = "user-token"
const sessionValidity = 60 * 60

type sessionStore struct {
	Db     *gorm.DB
	token  string
	userID uint
}

func (me sessionStore) Migrate() {
	me.Db.AutoMigrate(&models.Session{})
	me.Db.Model(&models.Session{}).AddForeignKey("owner_id", "users(id)", "CASCADE", "CASCADE")

}

func (me *sessionStore) Create(userID uint) (bool, error) {
	me.userID = userID

	out := models.Session{}
	if err := me.Db.Where("owner_id = ?", userID).Where("(UNIX_TIMESTAMP(expires) - UNIX_TIMESTAMP()) > 0").First(&out).Error; gorm.IsRecordNotFoundError(err) {
		// record not found => remove all from user and create fresh session
		me.DestroyAllByUserID(userID)

		token, err := uuid.NewV4()
		if err != nil {
			return false, err
		}

		session := models.Session{
			SessionID: token.String(),
			OwnerID:   userID,
			Expires:   time.Now().Add(time.Duration(sessionValidity) * time.Second),
			Validity:  sessionValidity,
		}

		if err := me.Db.Create(&session).Error; err != nil {
			return false, err
		}

		me.token = session.SessionID
		return true, nil
	}

	me.token = out.SessionID

	return true, nil
}

func (me *sessionStore) GetSession(r *http.Request) (*models.Session, error) {

	cookie, err := me.GetCookieFromRequest(r)
	if err != nil {
		return nil, err
	}

	if cookie.Value == "" {
		return nil, http.ErrNoCookie
	}

	var session = models.Session{}
	if errs := me.Db.Where("session_id = ?", cookie.Value).First(&session).GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorln(err)
		}
		return nil, errs[0]
	}

	me.token = session.SessionID
	me.userID = session.OwnerID

	return &session, nil
}

func (me *sessionStore) GetCookieFromRequest(r *http.Request) (*http.Cookie, error) {

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

func (me *sessionStore) Destroy(r *http.Request) (bool, error) {
	if me.token == "" {
		return false, http.ErrNoCookie
	}

	if errs := me.Db.Where("owner_id = ?", me.userID).Delete(&models.Session{}).GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorln(err)
		}
		return false, errs[0]
	}

	return true, nil
}

func (me *sessionStore) DestroyAllByUserID(userID uint) (bool, error) {
	if errs := me.Db.Where("owner_id = ?", userID).Delete(&models.Session{}).GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorln(err)
		}
		return false, errs[0]
	}

	return true, nil
}

func (me *sessionStore) CreateCookie() (*http.Cookie, error) {
	if me.token == "" {
		return nil, fmt.Errorf("cannot generate cookie without token")
	}

	return &http.Cookie{
		Name:    tokenKey,
		Value:   me.token,
		Expires: time.Now().Add(sessionValidity * time.Second),
	}, nil
}

func (me sessionStore) GetToken() string {
	return me.token
}

package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/goland-amaurybrisou/couchsport/api/models"
	"github.com/goland-amaurybrisou/couchsport/api/stores"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
)

type conversationHandler struct {
	Store *stores.StoreFactory
}

func (me conversationHandler) HandleMessage(w http.ResponseWriter, r *http.Request) {
	r.Close = true

	if r.Body != nil {
		defer r.Body.Close()
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	jsonBody := models.SendMessageBodyModel{}
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	_, err = jsonBody.Validate()
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	toProfile, err := me.Store.UserStore().GetProfile(jsonBody.ToID)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	fromUser, err := me.Store.UserStore().GetByEmail(jsonBody.Email, true)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	fromProfile, err := me.Store.UserStore().GetProfile(fromUser.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	//TODO send Email "account_created" with set password

	conversation, err := me.Store.ConversationStore().GetByReferents(fromProfile, toProfile)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	message := conversation.AddMessage(fromProfile.ID, jsonBody.Text)

	err = me.Store.ConversationStore().Save(conversation)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	json, err := json.Marshal(&message)

	if err != nil {
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(json))
}

func (me conversationHandler) ProfileConversations(userID uint, w http.ResponseWriter, r *http.Request) {
	profileID, err := me.Store.UserStore().GetProfileID(userID)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	conversations, err := me.Store.ConversationStore().ProfileConversations(profileID)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	json, err := json.Marshal(&conversations)

	if err != nil {
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(json))
}

func (me conversationHandler) Delete(userID uint, w http.ResponseWriter, r *http.Request) {
	r.Close = true

	if r.Body != nil {
		defer r.Body.Close()
	}

	query := r.URL.Query()

	tmp := query.Get("id")
	if tmp == "" {
		log.Println("id mising")
		http.Error(w, fmt.Errorf("id missing %s", tmp).Error(), http.StatusBadRequest)
		return
	}

	conversationID, err := strconv.Atoi(tmp)
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusUnprocessableEntity)
		return
	}

	owns, err := me.Store.UserStore().OwnConversation(userID, uint(conversationID))
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusUnprocessableEntity)
		return
	}

	if !owns {
		log.Error(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusForbidden)
		return
	}

	result, err := me.Store.ConversationStore().Delete(uint(conversationID))
	if err != nil {
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	json, err := json.Marshal(struct{ Result bool }{Result: result})

	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(json))
}

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

	if fromProfile.ID == toProfile.ID {
		http.Error(w, fmt.Errorf("%s", "invalid request").Error(), http.StatusBadRequest)
		return
	}

	//TODO send Email "account_created" with set password

	conversation, err := me.Store.ConversationStore().GetByReferents(fromProfile, toProfile)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	conversation, message, err := me.Store.ConversationStore().AddMessage(conversation, fromProfile.ID, toProfile.ID, jsonBody.Email, jsonBody.Text)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	if fromUser.New {
		go me.Store.MailStore().AccountAutoCreated(fromUser.Email, fromUser.PasswordTmp)
	}

	j, err := json.Marshal(&message)

	if err != nil {
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusInternalServerError)
		return
	}

	if !conversation.New {
		me.Store.WsStore().EmitToMutationNamespace(message.ToID, "CONVERSATION_ADD_MESSAGE", string(j), "conversations")
	} else {
		c, err := conversation.ToJSON()
		if err != nil {
			http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusInternalServerError)
			return
		}
		me.Store.WsStore().EmitToMutationNamespace(message.ToID, "NEW_CONVERSATION", c, "conversations")
	}

	fmt.Fprint(w, string(j))
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

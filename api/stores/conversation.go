package stores

import (
	"github.com/goland-amaurybrisou/couchsport/api/models"
	"github.com/jinzhu/gorm"
)

type conversationStore struct {
	Db *gorm.DB
}

//Migrate creates the db table
func (me conversationStore) Migrate() {
	me.Db.AutoMigrate(&models.Message{})
	me.Db.AutoMigrate(&models.Conversation{})
	me.Db.Model(&models.Message{}).AddForeignKey("owner_id", "conversations(id)", "CASCADE", "CASCADE")
	me.Db.Model(&models.Message{}).AddForeignKey("from_id", "profiles(id)", "CASCADE", "CASCADE")
	me.Db.Model(&models.Message{}).AddForeignKey("to_id", "profiles(id)", "CASCADE", "CASCADE")
	me.Db.Model(&models.Conversation{}).AddForeignKey("from_id", "profiles(id)", "CASCADE", "CASCADE")
	me.Db.Model(&models.Conversation{}).AddForeignKey("to_id", "profiles(id)", "CASCADE", "CASCADE")
}

//Delete a conversation by convID (softdelete)
func (me conversationStore) Delete(conversationID uint) (bool, error) {
	if err := me.Db.Exec("DELETE FROM conversations WHERE id = ?", conversationID).Error; err != nil {
		return false, err
	}
	return true, nil
}

//ProfileConversations fetch a profileID conversations
func (me conversationStore) ProfileConversations(profileID uint) ([]models.Conversation, error) {
	var outConversations []models.Conversation
	if err := me.Db.Preload("To").Preload("From").Preload("Messages").Where("from_id = ?", profileID).Or("to_id = ?", profileID).Find(&outConversations).Error; err != nil {
		return []models.Conversation{}, nil
	}
	return outConversations, nil
}

func (me conversationStore) GetByReferents(fromProfile, toProfile models.Profile) (models.Conversation, error) {
	outConversation := models.Conversation{}
	if err := me.Db.Model(&models.Conversation{}).Where("from_id = ? AND to_id = ?", fromProfile.ID, toProfile.ID).Or("from_id = ? AND to_id = ?", toProfile.ID, fromProfile.ID).FirstOrCreate(&outConversation, models.Conversation{
		FromID: fromProfile.ID,
		ToID:   toProfile.ID,
		New:    true,
	}).Error; err != nil {
		return models.Conversation{}, err
	}
	return outConversation, nil
}

//AddMessage in database
func (me conversationStore) AddMessage(conversation models.Conversation, fromID, toID uint, fromEmail, text string) (models.Conversation, models.Message, error) {
	m := models.Message{Text: text, OwnerID: conversation.ID, FromID: fromID, ToID: toID, Email: fromEmail}
	if err := me.Db.Create(&m).Error; err != nil {
		return models.Conversation{}, models.Message{}, err
	}

	conversation.Messages = append(conversation.Messages, m)

	return conversation, m, nil
}

func (me conversationStore) Save(conversation models.Conversation) (models.Conversation, error) {
	if err := me.Db.Model(&models.Conversation{}).Update(&conversation).Error; err != nil {
		return models.Conversation{}, err
	}
	return conversation, nil
}

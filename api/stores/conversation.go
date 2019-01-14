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
	me.Db.Model(&models.Conversation{}).AddForeignKey("from_id", "profiles(id)", "RESTRICT", "CASCADE")
	me.Db.Model(&models.Conversation{}).AddForeignKey("to_id", "profiles(id)", "RESTRICT", "CASCADE")
}

//Delete a conversation by convID
func (me conversationStore) Delete(convID uint) (bool, error) {
	return true, nil
}

//ProfileConversations fetch a profileID conversations
func (me conversationStore) ProfileConversations(profileID uint) ([]models.Conversation, error) {
	return []models.Conversation{}, nil
}

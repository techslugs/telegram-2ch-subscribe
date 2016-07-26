package storage

import (
	"gopkg.in/mgo.v2/bson"
)

type FailedChat struct {
	ChatID              int64 `bson:"chatID"`
	FailedAttemptsCount int64 `bson:"failedAttemptsCount"`
}

const (
	MaxFailedAttempts = 5
)

func chatByIDQuery(chatID int64) bson.M {
	return bson.M{"chatID": chatID}
}

func (storage *Storage) ReportChatSuccess(chatID int64) error {
	query := chatByIDQuery(chatID)
	change := bson.M{
		"$set": bson.M{"failedAttemptsCount": 0},
	}
	_, err := storage.FailedChats.UpdateAll(query, change)
	return err
}

func (storage *Storage) ReportChatError(chatID int64) error {
	query := chatByIDQuery(chatID)
	change := bson.M{
		"$inc": bson.M{"failedAttemptsCount": 1},
	}
	_, err := storage.FailedChats.Upsert(query, change)
	if err != nil {
		return err
	}
	return storage.unsubscribeFailedChat(chatID)
}

func (storage *Storage) unsubscribeFailedChat(chatID int64) error {
	query := chatByIDQuery(chatID)
	var failedChat FailedChat
	err := storage.FailedChats.Find(query).One(&failedChat)
	if err != nil {
		return err
	}
	if failedChat.FailedAttemptsCount < MaxFailedAttempts {
		return nil
	}
	err = storage.UnsubscribeChat("", failedChat.ChatID)
	if err != nil {
		return err
	}
	return storage.removeFailedChat(chatID)
}

func (storage *Storage) removeFailedChat(chatID int64) error {
	query := chatByIDQuery(chatID)
	return storage.FailedChats.Remove(query)
}

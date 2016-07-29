package storage

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type BoardSubscription struct {
	ChatID        int64    `bson:"chatID"`
	BoardName     string   `bson:"boardName"`
	MinScore      float64  `bson:"minScore"`
	SentThreadIDs []string `bson:"sentThreadIDs"`
	Timestamp     int64    `bson:"timestamp"`
}

func defaultTimestamp() int64 {
	return time.Now().Unix() - 30 // 30 seconds ago
}

func subscriptionQuery(boardName string, chatID int64) bson.M {
	return bson.M{"boardName": boardName, "chatID": chatID}
}

func (storage *Storage) SubscribeChat(boardName string, chatID int64, minScore float64) error {
	query := subscriptionQuery(boardName, chatID)
	change := bson.M{
		"$set":         bson.M{"minScore": minScore},
		"$setOnInsert": bson.M{"timestamp": defaultTimestamp()},
	}
	_, err := storage.BoardSubscriptions.Upsert(query, change)
	return err
}

func (storage *Storage) UnsubscribeChat(boardName string, chatID int64) error {
	var query bson.M
	if boardName != "" {
		query = subscriptionQuery(boardName, chatID)
	} else {
		query = bson.M{"chatID": boardName}
	}

	_, err := storage.BoardSubscriptions.RemoveAll(query)
	return err
}

func (storage *Storage) AllBoardSubscriptions(boardName string) ([]BoardSubscription, error) {
	query := bson.M{"boardName": boardName}
	var subscriptions []BoardSubscription
	err := storage.BoardSubscriptions.Find(query).All(&subscriptions)
	return subscriptions, err
}

func (storage *Storage) AllBoardNames() ([]string, error) {
	var names []string
	err := storage.BoardSubscriptions.Find(nil).Distinct("boardName", &names)
	if err != nil {
		return nil, err
	}

	return names, nil
}

func (storage *Storage) LogSentThread(boardName string, chatID int64, threadID string) error {
	query := subscriptionQuery(boardName, chatID)
	change := bson.M{
		"$addToSet": bson.M{"sentThreadIDs": threadID},
	}
	_, err := storage.BoardSubscriptions.UpdateAll(query, change)
	return err
}

func (storage *Storage) ClearStaleThreadIDs(boardName string, threadIDs []string) error {
	query := bson.M{"boardName": boardName}
	change := bson.M{
		"$pull": bson.M{"$nin": bson.M{"sentThreadIDs": threadIDs}},
	}
	_, err := storage.BoardSubscriptions.UpdateAll(query, change)
	return err
}

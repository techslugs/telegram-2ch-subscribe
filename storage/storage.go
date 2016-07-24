package storage

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Storage struct {
	DB                 *mgo.Database
	BoardSubscriptions *mgo.Collection
	FailedChats        *mgo.Collection
	BoardDescriptions  *mgo.Collection
}

var (
	BoardsNameIndex = mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	ChatIDIndex = mgo.Index{
		Key:        []string{"chatID"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
)

func NewStorage(DB *mgo.Database) (*Storage, error) {
	storage := Storage{
		DB:                 DB,
		BoardSubscriptions: DB.C("board_subscriptions"),
		FailedChats:        DB.C("failed_chats"),
		BoardDescriptions:  DB.C("board_descriptions"),
	}

	err := storage.BoardSubscriptions.EnsureIndex(BoardsNameIndex)
	if err != nil {
		return nil, err
	}
	err = storage.BoardDescriptions.EnsureIndex(BoardsNameIndex)
	if err != nil {
		return nil, err
	}
	err = storage.FailedChats.EnsureIndex(ChatIDIndex)
	if err != nil {
		return nil, err
	}

	return &storage, nil
}

func boardByNameQuery(boardName string) bson.M {
	return bson.M{"name": boardName}
}

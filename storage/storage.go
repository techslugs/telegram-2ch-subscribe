package storage

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Storage struct {
	DB                 *mgo.Database
	BoardSubscriptions *mgo.Collection
	BoardDescriptions  *mgo.Collection
}

func NewStorage(DB *mgo.Database) (*Storage, error) {
	storage := Storage{
		DB:                 DB,
		BoardSubscriptions: DB.C("board_subscriptions"),
		BoardDescriptions:  DB.C("board_descriptions"),
	}

	boardsNameIndex := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := storage.BoardSubscriptions.EnsureIndex(boardsNameIndex)
	if err != nil {
		return nil, err
	}
	err = storage.BoardDescriptions.EnsureIndex(boardsNameIndex)
	if err != nil {
		return nil, err
	}

	return &storage, nil
}

func boardByNameQuery(boardName string) bson.M {
	return bson.M{"name": boardName}
}

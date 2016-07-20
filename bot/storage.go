package bot

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Storage struct {
	DB     *mgo.Database
	Boards *mgo.Collection
}

type Board struct {
	Name      string
	ChatIDs   []string
	Timestamp int64
}

func NewStorage(DB *mgo.Database) (*Storage, error) {
	storage := Storage{DB: DB, Boards: DB.C("boards")}

	boardsNameIndex := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := storage.Boards.EnsureIndex(boardsNameIndex)
	if err != nil {
		return nil, err
	}

	return &storage, nil
}

func boardByNameQuery(boardName string) bson.M {
	return bson.M{"name": boardName}
}

func (storage *Storage) SubscribeChat(boardName string, chatID string) error {
	query := boardByNameQuery(boardName)
	change := bson.M{"$addToSet": bson.M{"chatIDs": chatID}}
	_, err := storage.Boards.Upsert(query, change)
	return err
}

func (storage *Storage) UnsubscribeChat(boardName string, chatID string) error {
	var query bson.M
	if boardName != "" {
		query = nil
	} else {
		query = bson.M{}
	}

	change := bson.M{"$pull": bson.M{"chatIDs": chatID}}
	_, err := storage.Boards.UpdateAll(query, change)
	return err
}

func (storage *Storage) AllBoardNames() ([]string, error) {
	var boardsWithNames []Board
	err := storage.Boards.Find(nil).Select(bson.M{"name": 1}).All(&boardsWithNames)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(boardsWithNames))
	for i, board := range boardsWithNames {
		names[i] = board.Name
	}
	return names, nil
}

func (storage *Storage) BoardDetails(boardName string) (*Board, error) {
	query := boardByNameQuery(boardName)
	var board Board
	err := storage.Boards.Find(query).One(&board)
	return &board, err
}

func (storage *Storage) UpdateBoardTimestamp(boardName string, timestamp int64) error {
	query := boardByNameQuery(boardName)
	change := bson.M{"$set": bson.M{"timestamp": timestamp}}
	_, err := storage.Boards.Upsert(query, change)
	return err
}

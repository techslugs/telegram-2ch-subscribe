package storage

import (
	"gopkg.in/mgo.v2/bson"
)

type BoardDescription struct {
	Name        string `bson:"name"`
	FullName    string `bson:"fullName"`
	Description string `bson:"description"`
}

func (storage *Storage) SaveBoardDescription(
	boardName, fullName, description string,
) error {
	query := boardByNameQuery(boardName)
	change := bson.M{
		"$set": bson.M{"fullName": fullName, "description": description},
	}
	_, err := storage.BoardDescriptions.Upsert(query, change)
	return err
}

func (storage *Storage) FilterInvalidBoardNames(boardNames []string) ([]string, error) {
	var boardsWithNames []BoardDescription
	query := bson.M{"name": bson.M{"$in": boardNames}}
	err := storage.BoardDescriptions.Find(query).Select(bson.M{"name": 1}).All(&boardsWithNames)

	if err != nil {
		return nil, err
	}

	names := make([]string, len(boardsWithNames))
	for i, board := range boardsWithNames {
		names[i] = board.Name
	}
	return names, nil
}

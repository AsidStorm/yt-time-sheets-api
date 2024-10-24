package jsonapi

import (
	"encoding/json"
	"yandex.tracker.api/domain/models"
)

const boardType = "board"

type board struct {
	Id         int64           `json:"id"`
	Type       string          `json:"type"`
	Attributes boardAttributes `json:"attributes"`
}

type boardAttributes struct {
	Name string `json:"name"`
}

func makeBoard(b models.Board) board {
	return board{
		Id:   b.Id,
		Type: boardType,
		Attributes: boardAttributes{
			Name: b.Name,
		},
	}
}

func MarshalBoards(in []models.Board) ([]byte, error) {
	response := struct {
		Data []board `json:"data"`
	}{make([]board, len(in))}

	for i, b := range in {
		response.Data[i] = makeBoard(b)
	}

	return json.Marshal(response)
}

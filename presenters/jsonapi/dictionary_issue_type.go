package jsonapi

import (
	"encoding/json"
	"yandex.tracker.api/domain/models"
)

type dictionaryIssueType struct {
	Id   int64  `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

func makeDictionaryIssueType(t models.DictionaryIssueType) dictionaryIssueType {
	return dictionaryIssueType{
		Id:   t.Id,
		Key:  t.Key,
		Name: t.Name,
	}
}

func MarshalDictionaryIssueTypes(types []models.DictionaryIssueType) ([]byte, error) {
	out := struct {
		Data []dictionaryIssueType `json:"data"`
	}{make([]dictionaryIssueType, len(types))}

	for i, t := range types {
		out.Data[i] = makeDictionaryIssueType(t)
	}

	return json.Marshal(out)
}

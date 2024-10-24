package jsonapi

import (
	"encoding/json"
	"fmt"
)

type FilterTasks struct {
	Query string `json:"query"`
}

func UnmarshalFilterTasks(in []byte) (*FilterTasks, error) {
	out := &FilterTasks{}

	if err := json.Unmarshal(in, &out); err != nil {
		return nil, fmt.Errorf("unable to unmarshal filter tasks due [%s]", err)
	}

	return out, nil
}

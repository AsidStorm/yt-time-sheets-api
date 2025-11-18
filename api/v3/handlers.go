package v3

import (
	"encoding/json"
	"io"
	"net/http"

	"yandex.tracker.api/domain"
	"yandex.tracker.api/domain/cases/result_v3"
)

func Result(_ map[string]string, c domain.Context, _ http.ResponseWriter, r *http.Request) (int, []byte, error) {
	var request result_v3.Request

	in, err := io.ReadAll(r.Body)
	if err != nil {
		return BadRequest(err)
	}

	if err := json.Unmarshal(in, &request); err != nil {
		return BadRequest(err)
	}

	response, err := result_v3.Run(c, request)
	if err != nil {
		return InternalServerError(err)
	}

	out, err := json.Marshal(response)
	if err != nil {
		return InternalServerError(err)
	}

	return OK(out)
}

package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/birat/restapi/internal/types"
	"github.com/birat/restapi/internal/utils/response"
)

func NewStudent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			// JSON Response
			response.WriteJson(w, http.StatusBadRequest, err.Error())

			return
		}

		slog.Info("creating a student")

		response.WriteJson(w, http.StatusCreated, map[string]interface{}{"success": "ok", "student": student})
	}
}
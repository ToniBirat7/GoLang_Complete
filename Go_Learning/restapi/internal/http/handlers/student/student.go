package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/birat/restapi/internal/types"
	"github.com/birat/restapi/internal/utils/response"
	vld "github.com/go-playground/validator/v10"
)

func NewStudent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			// JSON Response
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))

			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Validate the data
		err = vld.New().Struct(student)

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		slog.Info("creating a student")

		response.WriteJson(w, http.StatusCreated, map[string]interface{}{"success": "ok", "student": student})
	}
}

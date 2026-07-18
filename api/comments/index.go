package handler

import (
	"encoding/json"
	"net/http"

	"newspaper-site/lib/shared"
)

type commentInput struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

// Handler serves:
//
//	GET  /api/comments  -> list of reader letters, oldest first
//	POST /api/comments  -> add a new reader letter {name, text}
func Handler(w http.ResponseWriter, r *http.Request) {
	if shared.HandleCORSPreflight(w, r) {
		return
	}

	switch r.Method {
	case http.MethodGet:
		list, err := shared.LoadComments()
		if err != nil {
			shared.WriteError(w, http.StatusInternalServerError, "Could not load comments right now")
			return
		}
		shared.WriteJSON(w, http.StatusOK, map[string]interface{}{"comments": list})

	case http.MethodPost:
		var in commentInput
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			shared.WriteError(w, http.StatusBadRequest, "Malformed request body")
			return
		}
		if len(in.Text) == 0 || len(in.Text) > 4000 {
			shared.WriteError(w, http.StatusBadRequest, "A comment needs some text (under 4000 characters)")
			return
		}
		c, err := shared.AddComment(in.Name, in.Text)
		if err != nil {
			shared.WriteError(w, http.StatusInternalServerError, "Could not save your comment right now")
			return
		}
		shared.WriteJSON(w, http.StatusCreated, c)

	default:
		shared.WriteError(w, http.StatusMethodNotAllowed, "Use GET or POST")
	}
}

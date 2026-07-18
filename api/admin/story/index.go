package handler

import (
	"encoding/json"
	"net/http"

	"newspaper-site/lib/shared"
)

type storyInput struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// Handler serves POST /api/admin/story {title, body} to overwrite the
// published manuscript. Requires a valid Authorization: Bearer <token> header.
func Handler(w http.ResponseWriter, r *http.Request) {
	if shared.HandleCORSPreflight(w, r) {
		return
	}
	if !shared.RequireAdmin(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		shared.WriteError(w, http.StatusMethodNotAllowed, "Use POST")
		return
	}

	var in storyInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		shared.WriteError(w, http.StatusBadRequest, "Malformed request body")
		return
	}
	if len(in.Body) == 0 {
		shared.WriteError(w, http.StatusBadRequest, "The story body can't be empty")
		return
	}

	if err := shared.SaveStory(in.Title, in.Body); err != nil {
		shared.WriteError(w, http.StatusInternalServerError, "Could not save the story right now")
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]bool{"saved": true})
}

package handler

import (
	"net/http"

	"newspaper-site/lib/shared"
)

// Handler serves GET /api/story - the current title + body of the piece.
func Handler(w http.ResponseWriter, r *http.Request) {
	if shared.HandleCORSPreflight(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		shared.WriteError(w, http.StatusMethodNotAllowed, "Use GET")
		return
	}

	title, body, err := shared.GetStory()
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, "Could not load the story right now")
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]string{
		"title": title,
		"body":  body,
	})
}

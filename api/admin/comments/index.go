package handler

import (
	"net/http"

	"newspaper-site/lib/shared"
)

// Handler serves:
//
//	GET    /api/admin/comments           -> all comments (admin view)
//	DELETE /api/admin/comments?id=<id>   -> remove one comment
//
// Both require a valid Authorization: Bearer <token> header.
func Handler(w http.ResponseWriter, r *http.Request) {
	if shared.HandleCORSPreflight(w, r) {
		return
	}
	if !shared.RequireAdmin(w, r) {
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

	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			shared.WriteError(w, http.StatusBadRequest, "Missing id")
			return
		}
		found, err := shared.DeleteComment(id)
		if err != nil {
			shared.WriteError(w, http.StatusInternalServerError, "Could not delete that comment right now")
			return
		}
		if !found {
			shared.WriteError(w, http.StatusNotFound, "That comment no longer exists")
			return
		}
		shared.WriteJSON(w, http.StatusOK, map[string]bool{"deleted": true})

	default:
		shared.WriteError(w, http.StatusMethodNotAllowed, "Use GET or DELETE")
	}
}

package handler

import (
	"net/http"
	"strconv"

	"newspaper-site/lib/shared"
)

const viewCountKey = "view_count"

// Handler serves GET /api/admin/stats -> {views, comments}
func Handler(w http.ResponseWriter, r *http.Request) {
	if shared.HandleCORSPreflight(w, r) {
		return
	}
	if !shared.RequireAdmin(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		shared.WriteError(w, http.StatusMethodNotAllowed, "Use GET")
		return
	}

	raw, ok, err := shared.Get(viewCountKey)
	views := int64(0)
	if err == nil && ok {
		if n, convErr := strconv.ParseInt(raw, 10, 64); convErr == nil {
			views = n
		}
	}

	comments, err := shared.LoadComments()
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, "Could not load stats right now")
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"views":    views,
		"comments": len(comments),
	})
}

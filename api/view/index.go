package handler

import (
	"net/http"

	"newspaper-site/lib/shared"
)

const viewCountKey = "view_count"

// Handler serves POST /api/view - fired once per page load from the
// frontend to bump the public view counter. Intentionally has no read
// path here; the count is only surfaced to the author via /api/admin/stats.
func Handler(w http.ResponseWriter, r *http.Request) {
	if shared.HandleCORSPreflight(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		shared.WriteError(w, http.StatusMethodNotAllowed, "Use POST")
		return
	}

	if _, err := shared.Incr(viewCountKey); err != nil {
		// Don't fail the page load over a stats hiccup.
		shared.WriteJSON(w, http.StatusOK, map[string]bool{"ok": false})
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

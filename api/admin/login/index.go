package handler

import (
	"encoding/json"
	"net/http"

	"newspaper-site/lib/shared"
)

type loginInput struct {
	Password string `json:"password"`
}

// Handler serves POST /api/admin/login {password} -> {token}
func Handler(w http.ResponseWriter, r *http.Request) {
	if shared.HandleCORSPreflight(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		shared.WriteError(w, http.StatusMethodNotAllowed, "Use POST")
		return
	}

	var in loginInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		shared.WriteError(w, http.StatusBadRequest, "Malformed request body")
		return
	}

	if !shared.CheckPassword(in.Password) {
		shared.WriteError(w, http.StatusUnauthorized, "Incorrect password")
		return
	}

	token, err := shared.CreateSession()
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, "Could not start a session right now")
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

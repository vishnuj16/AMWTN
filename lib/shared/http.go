package shared

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

// WriteJSON writes v as a JSON response with the given status code and the
// permissive-but-scoped headers every function in this project needs.
func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// WriteError is a small convenience wrapper around WriteJSON for error bodies.
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}

// HandleCORSPreflight writes the response for an OPTIONS preflight request
// and reports whether the caller should return immediately afterwards.
func HandleCORSPreflight(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return true
	}
	return false
}

// NewToken generates a random 32-byte hex session token.
func NewToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// sessionKey namespaces admin session tokens in Redis.
func sessionKey(token string) string {
	return "session:" + token
}

// CreateSession stores a new admin session token with a 12 hour TTL.
func CreateSession() (string, error) {
	token, err := NewToken()
	if err != nil {
		return "", err
	}
	if err := SetEx(sessionKey(token), "1", 12*60*60); err != nil {
		return "", err
	}
	return token, nil
}

// RequireAdmin checks the Authorization: Bearer <token> header against the
// active sessions stored in Redis. It writes a 401 response and returns
// false if authentication fails, so callers can just `if !RequireAdmin(...) { return }`.
func RequireAdmin(w http.ResponseWriter, r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		WriteError(w, http.StatusUnauthorized, "Not authenticated")
		return false
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	if token == "" {
		WriteError(w, http.StatusUnauthorized, "Not authenticated")
		return false
	}
	ok, err := Exists(sessionKey(token))
	if err != nil || !ok {
		WriteError(w, http.StatusUnauthorized, "Session expired, please log in again")
		return false
	}
	return true
}

// CheckPassword does a constant-time comparison against ADMIN_PASSWORD.
func CheckPassword(candidate string) bool {
	want := os.Getenv("ADMIN_PASSWORD")
	if want == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(candidate), []byte(want)) == 1
}

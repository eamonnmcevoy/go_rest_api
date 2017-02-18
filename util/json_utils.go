package util

import (
  "net/http"
  "encoding/json"
)

func RespondWithError(w http.ResponseWriter, code int, message string) {
  RespondWithJSON(w, code, map[string]string{"error": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
  response, _ := json.Marshal(payload)

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(code)
  w.Write(response)
}
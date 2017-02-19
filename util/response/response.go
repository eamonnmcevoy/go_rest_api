package response

import (
  "net/http"
  "encoding/json"
)

func Error(w http.ResponseWriter, code int, message string) {
  Json(w, code, map[string]string{"error": message})
}

func Json(w http.ResponseWriter, code int, payload interface{}) {
  response, _ := json.Marshal(payload)

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(code)
  w.Write(response)
}
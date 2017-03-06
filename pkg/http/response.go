package http

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

func JsonWithCookie(w http.ResponseWriter, code int, payload interface{}, cookie http.Cookie) {
  response, _ := json.Marshal(payload)
  http.SetCookie(w, &cookie)

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(code)
  w.Write(response)
}
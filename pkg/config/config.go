package config

import (
  "os"
  "go_rest_api/pkg"
)

func GetConfig() *root.Config {
  return &root.Config {
    Mongo: &root.MongoConfig {
      Ip: envOrDefaultString("go_rest_api:mongo:ip", "127.0.0.1:27017"),
      DbName: envOrDefaultString("go_rest_api:mongo:dbName", "myDb")},
    Server: &root.ServerConfig { Port: envOrDefaultString("go_rest_api:server:port", ":1377")},
    Auth: &root.AuthConfig { Secret: envOrDefaultString("go_rest_api:auth:secret", "mysecret")}}
}

func envOrDefaultString(envVar string, defaultValue string) string {
  value := os.Getenv(envVar)
  if value == "" {
    return defaultValue;
  }
  
  return value
}
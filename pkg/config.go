package root

type MongoConfig struct {
  Ip string `json:"ip"`
  DbName string `json:"dbName"`
}

type ServerConfig struct {
   Port string `json:"port"`
}

type AuthConfig struct {
  Secret string `json:"secret"`
}

type Config struct {
  Mongo  *MongoConfig  `json:"mongo"`
  Server *ServerConfig `json:"server"`
  Auth   *AuthConfig   `json:"auth"`
}
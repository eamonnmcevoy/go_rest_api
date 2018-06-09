package root

type User struct {
  Id           string  `json:"id"`
  Username     string  `json:"username"`
  Password     string  `json:"password"`
}

type UserService interface {
  CreateUser(u *User) error
  GetUserByUsername(username string) (error, User)
  Login(c Credentials) (error, User)
}
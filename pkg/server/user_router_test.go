package server

import (
  "fmt"
  "bytes"
  "errors"
  "context"
  "testing"
  "net/http"
  "net/http/httptest"
  "encoding/json"
  "go_rest_api/pkg"
  "go_rest_api/pkg/mock"
  "github.com/gorilla/mux"
  "github.com/dgrijalva/jwt-go"
)

//createUserHandler tests
func Test_UserRouter_createUserHandler(t *testing.T) {
  t.Run("happy path", createUserHandler_should_pass_User_object_to_UserService_CreateUser)
  t.Run("invalid payload", createUserHandler_should_return_StatusBadRequest_if_payload_is_invalid)
  t.Run("internal error", createUserHandler_should_return_StatusInternalServerError_if_UserService_returns_error)
}

func createUserHandler_should_pass_User_object_to_UserService_CreateUser(t *testing.T) {
  // Arrange
  us := mock.UserService{}
  test_mux := NewUserRouter(&us, mux.NewRouter())
  var result *root.User 
  us.CreateUserFn = func(u *root.User) error {
    result = u
    return nil
  }

  testUsername := "test_username"
  testPassword := "test_password"

  values := map[string]string{"username": testUsername, "password": testPassword}
  jsonValue, _ := json.Marshal(values)
  payload := bytes.NewBuffer(jsonValue)

  // Act
  w := httptest.NewRecorder()
  r, _ := http.NewRequest("PUT", "/", payload)
  r.Header.Set("Content-Type", "application/json")
  test_mux.ServeHTTP(w,r)
  
  // Assert
  if !us.CreateUserInvoked {
    t.Fatal("expected CreateUser() to be invoked")
  }
  if result.Username != testUsername {
    t.Fatalf("expected username to be: `%s`, got: `%s`", testUsername, result.Username)    
  }
  if result.Password != testPassword {
    t.Fatalf("expected username to be: `%s`, got: `%s`", testPassword, result.Password)     
  }  
}

func createUserHandler_should_return_StatusBadRequest_if_payload_is_invalid(t *testing.T) {
  //Arrange
  us := mock.UserService{}
  test_mux := NewUserRouter(&us, mux.NewRouter())
  us.CreateUserFn = func(u *root.User) error {
    return nil
  }

  //Act
  w := httptest.NewRecorder()
  r, _ := http.NewRequest("PUT", "/", nil)
  r.Header.Set("Content-Type", "application/json")
  test_mux.ServeHTTP(w,r)

  //Assert
  if w.Code != http.StatusBadRequest {
    t.Fatal("expected: http.StatusBadRequest, got: %i",w.Code)
  }
}

func createUserHandler_should_return_StatusInternalServerError_if_UserService_returns_error(t *testing.T) {
  //Arrange
  us := mock.UserService{}
  test_mux := NewUserRouter(&us, mux.NewRouter())
  us.CreateUserFn = func(u *root.User) error {
    return errors.New("user service error")
  }

  values := map[string]string{"username": "", "password": ""}
  jsonValue, _ := json.Marshal(values)
  payload := bytes.NewBuffer(jsonValue)

  //Act
  w := httptest.NewRecorder()
  r, _ := http.NewRequest("PUT", "/", payload)
  r.Header.Set("Content-Type", "application/json")
  test_mux.ServeHTTP(w,r)

  //Assert
  if w.Code != http.StatusInternalServerError {
    t.Fatal("expected: http.StatusInternalServerError, got: %i",w.Code)
  }
}

//profileHandler tests
func Test_UserRouter_profileHandler(t *testing.T) {
  t.Run("happy path", profileHandler_should_return_User_from_context)
  t.Run("no context", profileHandler_should_return_StatusBadRequest_if_no_auth_context)
  t.Run("user not found", profileHandler_should_return_StatusNotFound_if_no_user_found)
}

func profileHandler_should_return_User_from_context(t *testing.T) {
  // Arrange
  us := mock.UserService{}
  test_mux := NewUserRouter(&us, mux.NewRouter())
  var result string 
  us.GetUserByUsernameFn = func(username string) (error, root.User) {
    result = username
    return nil, root.User{}
  }

  testUsername := "test_username"
  testUser := root.User{Username:testUsername}

  // Act
  w := httptest.NewRecorder()
  r, _ := http.NewRequest("GET", "/profile", nil)
  testCookie := newAuthCookie(testUser)
  r.AddCookie(&testCookie)
  ctx := context.WithValue(r.Context(), contextKeyAuthtoken, claims { testUsername, jwt.StandardClaims{} })
  test_mux.ServeHTTP(w,r.WithContext(ctx))
  
  // Assert
  if !us.GetUserByUsernameInvoked {
    t.Fatal("expected GetUserByUsername() to be invoked")
  }
  if result != testUsername {
    t.Fatalf("expected username to be: `%s`, got: `%s`", testUsername, result)    
  }
}

func profileHandler_should_return_StatusBadRequest_if_no_auth_context(t *testing.T) {
  //Arrange
  us := mock.UserService{}
  test_mux := NewUserRouter(&us, mux.NewRouter())

  //Act
  w := httptest.NewRecorder()
  r, _ := http.NewRequest("GET", "/profile", nil)
  test_mux.ServeHTTP(w,r)

  //Assert
  if w.Code != http.StatusUnauthorized{
    t.Fatalf("expected StatusUnauthorized, got: %s",w.Code)
  }
}

func profileHandler_should_return_StatusNotFound_if_no_user_found(t *testing.T) {
  //Arrange
  us := mock.UserService{}
  test_mux := NewUserRouter(&us, mux.NewRouter())
  var result string 
  us.GetUserByUsernameFn = func(username string) (error, root.User) {
    result = username
    return errors.New("user service error"), root.User{}
  }
  testUsername := "test_username"
  testUser := root.User{Username:testUsername}

  //Act
  w := httptest.NewRecorder()
  r, _ := http.NewRequest("GET", "/profile", nil)
  testCookie := newAuthCookie(testUser)
  r.AddCookie(&testCookie)
  ctx := context.WithValue(r.Context(), contextKeyAuthtoken, claims { testUsername, jwt.StandardClaims{} })
  test_mux.ServeHTTP(w,r.WithContext(ctx))

  //Assert
  if !us.GetUserByUsernameInvoked {
    t.Fatal("expected GetUserByUsername() to be invoked")
  }
  if w.Code != http.StatusNotFound {
    t.Fatalf("expected: StatusNotFound, got: %s", w.Code)
  }
}

//getUserHandler tests
func Test_UserRouter_getUserHandler(t *testing.T) {
  t.Run("happy path", getUserHandler_should_call_GetUserByUsername_with_username_from_querystring)
  t.Run("no user found", getUserHandler_should_return_StatusNotFound_if_no_user_found)
}

func getUserHandler_should_call_GetUserByUsername_with_username_from_querystring(t *testing.T) {
  // Arrange
  us := mock.UserService{}
  test_mux := NewUserRouter(&us, mux.NewRouter())
  var result string 
  us.GetUserByUsernameFn = func(username string) (error, root.User) {
    result = username
    return nil, root.User{}
  }

  testUsername := "test_username"

  // Act
  w := httptest.NewRecorder()
  r, _ := http.NewRequest("GET", "/"+testUsername, nil)
  test_mux.ServeHTTP(w,r)
  
  // Assert
  if !us.GetUserByUsernameInvoked {
    t.Fatal("expected GetUserByUsername() to be invoked")
  }
  if result != testUsername {
    t.Fatalf("expected username to be: `%s`, got: `%s`", testUsername, result)    
  }
}

func getUserHandler_should_return_StatusNotFound_if_no_user_found(t *testing.T) {
  // Arrange
  us := mock.UserService{}
  test_mux := NewUserRouter(&us, mux.NewRouter())
  var result string 
  us.GetUserByUsernameFn = func(username string) (error, root.User) {
    result = username
    return errors.New("user service error"), root.User{}
  }

  testUsername := "test_username"

  // Act
  w := httptest.NewRecorder()
  r, _ := http.NewRequest("GET", "/"+testUsername, nil)
  test_mux.ServeHTTP(w,r)
  
  // Assert
  if !us.GetUserByUsernameInvoked {
    t.Fatal("expected GetUserByUsername() to be invoked")
  }
  if w.Code != http.StatusNotFound {
    t.Fatalf("expected: StatusNotFound, got: %s", w.Code)
  }
}

//gHandler tests
func Test_UserRouter_loginHandler(t *testing.T) {
  fmt.Println("loginHandler tests")
  t.Run("happy path", loginHandler_should_provide_new_auth_cookie_if_userService_returns_a_user)
  //t.Run("no user found", getUserHandler_should_return_StatusNotFound_if_no_user_found)
}

func loginHandler_should_provide_new_auth_cookie_if_userService_returns_a_user(t *testing.T) {
  // Arrange
  us := mock.UserService{}
  test_mux := NewUserRouter(&us, mux.NewRouter())
  var result string 
  us.LoginFn = func(credentials root.Credentials) (error, root.User) {
    result = credentials.Username
    return nil, root.User{}
  }

  testUsername := "test_username"
  testPassword := "test_password"

  values := map[string]string{"username": testUsername, "password": testPassword}
  jsonValue, _ := json.Marshal(values)
  payload := bytes.NewBuffer(jsonValue)

  // Act
  w := httptest.NewRecorder()
  r, _ := http.NewRequest("POST", "/login", payload)
  test_mux.ServeHTTP(w,r)
  
  // Assert
  if !us.LoginInvoked {
    t.Fatal("expected Login() to be invoked")
  }

  request := &http.Request{Header: http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}}
  cookie, err := request.Cookie("Auth")
  if err != nil || cookie == nil {
    panic("Expected Cookie named 'Auth'")
  }
}
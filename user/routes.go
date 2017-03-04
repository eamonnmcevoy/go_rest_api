package user

import (
  "fmt"
  "time"
  "strconv"
  "context"
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"
  "gopkg.in/mgo.v2"
  "go_rest_api/util/response"
  "github.com/dgrijalva/jwt-go"
)

type userRouter struct {
  service IuserService
}

func NewUserRouter(session *mgo.Session, router *mux.Router) *mux.Router {
  u := userRouter{NewUserService(session)}

  router.HandleFunc("/create", u.createUserHandler)
  router.HandleFunc("/get", validate(u.getUserHandler))
  router.HandleFunc("/authenticate", u.authenticateHandler)
  return router
}

func(u* userRouter) createUserHandler(w http.ResponseWriter, r *http.Request) {
  err, credentials := decodeCredentials(r)
  if err != nil {
    response.Error(w, http.StatusBadRequest, "Invalid request payload")
    return
  }

  err = u.service.createUser(credentials)
  if err != nil {
    response.Error(w, http.StatusInternalServerError, err.Error())
    return
  }

  response.Json(w, http.StatusOK, err)
}

func(u* userRouter) getUserHandler(w http.ResponseWriter, r *http.Request) {
  err, findUser := decodeUser(r)
  if err != nil {
    response.Error(w, http.StatusBadRequest, "Invalid request payload")
    return
  }

  var user User
  err, user = u.service.getUserByUsername(findUser.Username)
  if err != nil {
    response.Error(w, http.StatusInternalServerError, err.Error())
    return
  }

  response.Json(w, http.StatusOK, user)
}

func(u *userRouter) authenticateHandler(w http.ResponseWriter, r *http.Request) {
  err, credentials := decodeCredentials(r)
  if err != nil {
    response.Error(w, http.StatusBadRequest, "Invalid request payload")
    return
  }

  var user User
  err, user = u.service.authenticate(credentials)
  if err == nil {
    token := u.getToken(user)
    expireCookie := time.Now().Add(time.Hour * 1)
    cookie := http.Cookie {
      Name: "Auth",
      Value: token,
      Expires: expireCookie,
      HttpOnly: true }
    fmt.Println("routes")
    fmt.Println(contextKeyAuthtoken.String())
    fmt.Println(cookie)

    response.JsonWithCookie(w, http.StatusOK, user, cookie)
  } else {
    response.Error(w, http.StatusInternalServerError, "Incorrect password")
  }
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////
//AUTH

type Claims struct {
    Username string `json:"username"`
    // recommended having
    jwt.StandardClaims
}

type contextKey string
func (c contextKey) String() string {
    return "mypackage context key " + string(c)
}

var (
    contextKeyAuthtoken = contextKey("auth-token")
)

func(u *userRouter) getToken(user User) string {
  // Expires the token and cookie in 1 hour
  expireToken := time.Now().Add(time.Hour * 1).Unix()

  // We'll manually assign the claims but in production you'd insert values from a database 
  claims := Claims {
    user.Username,
    jwt.StandardClaims {
      ExpiresAt: expireToken,
      Issuer: "localhost !",
    },
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

  /* Sign the token with our secret */
  tokenString, _ := token.SignedString([]byte("secret"))
    
  /* Finally, write the token to the browser window */
   return tokenString
}

func validate(next http.HandlerFunc) http.HandlerFunc {
  return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
    //If no Auth cookie is set then return a 404 not found
    cookie, err := req.Cookie("Auth")
    if err != nil {
      response.Error(res, http.StatusUnauthorized, "No authorization cookie")
      return
    }
       
    // Return a Token using the cookie
    token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(token *jwt.Token) (interface{}, error){
      // Make sure token's signature wasn't changed
      if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("Unexpected siging method")    
      }    
      return []byte("secret"), nil
    })

    if err != nil {
      response.Error(res, http.StatusUnauthorized, "Invalid token")
      return
    }
       
    // Grab the tokens claims and pass it into the original request
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
      fmt.Println("user in cookie: " + claims.Username)
      fmt.Println("cookie expires: " + strconv.FormatInt(claims.StandardClaims.ExpiresAt, 10))
      fmt.Println("cookie issuer: " + claims.StandardClaims.Issuer)
      ctx := context.WithValue(req.Context(), contextKeyAuthtoken, *claims)
      next(res, req.WithContext(ctx))
    } else {
      response.Error(res, http.StatusUnauthorized, "Unauthorized")
      return
    }
  })
}

func decodeUser(r *http.Request) (error,User) {
  var u User
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&u)
  return err, u
}

func decodeCredentials(r *http.Request) (error,Credentials) {
  var c Credentials
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&c)
  return err, c
}  
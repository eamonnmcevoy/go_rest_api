# go_rest_api
Go web server with mongoDb, gorilla toolkit, and jwt authentication

This project follows the package layout recommended by Ben Johnson in his Medium article - Standard package layout
https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1

Structure
~~~
go_rest_api
  /cmd
    /app
      -app.go
      -main.go
    /pkg
      /mongo
        -session.go
        -user_service.go
        -user_model.go
      /server
        -server.go
        -user_router.go
        -response.go
        -auth.go
      /mock
        -user_service_mock.go
      -user.go
      -credentials.go
~~~

Dependancies
  - go get gopkg.in/mgo.v2
  - go get github.com/gorilla/mux
  - go get github.com/gorilla/handlers
  - go get golang.org/x/crypto/bcrypt
  - go get github.com/google/uuid

References
  - https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1#.sdcvblyts
  - https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql
  - https://dinosaurscode.xyz/go/2016/06/17/golang-jwt-authentication/
  - https://medium.com/@matryer/context-keys-in-go-5312346a868d#.hb4spbx1a

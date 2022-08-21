# Authentication in Golang with JWTs. Using auth0.com

[Link](https://auth0.com/blog/authentication-in-golang/)
Beego, Gin Gionix, Echo, Revel

Go Modules, GOPATH

`go mod init github.com/auth-go`
`go get`
`go env -w GO111MODULE=off`

`Negroni, Alice`
middlewares
here `net/http`

in front auth0 is work

There is a problem in main.go. autho0 change api and there is no jwtmiddleware.Options. Tried to use "github.com/auth0/go-jwt-middleware/v2"
stop on error $GOROOT, $GOPATHm

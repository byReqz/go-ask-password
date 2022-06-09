# go-ask-password
systemd-ask-password like password prompt for go

# example usage
```go
pw, err := AskPassword.AskPassword("Password: ")
if err != nil {
	log.Fatal(err)
}
fmt.Println(pw)
```
also in <a href=./demo/password.go>demo/password.go</a> and <a href=./demo/substitute.go>demo/substitute.go</a>
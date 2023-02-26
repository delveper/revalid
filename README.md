# ReValid

> Validator for *Go* structs

### Usage example

```go
type User struct {
	FirstName string `regexp:"[A-Za-z]{2,255}"`
	LastName string `regexp:"[A-Za-z]{2,255}"`	
	Password string `regexp:".{8,255}"`
}

func main() {
	usr := User {
		FirstName: "Jim",
		LastName: "Don",
		Password: "qwerty"
    }
	
    if err := revalid.ValidateStruct(usr); err != nil {
		log.Println(err) // "User has to have valid Password according to pattern: `.{8,255}`"
    }
}
```

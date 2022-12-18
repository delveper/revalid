# ReValid

## Validator for *GO* structs
> Validation is aware of zero values and nested structs.
### Usage example

```go
type User struct {
	FirstName string `regex:"[A-Za-z]{2,255}"`
	LastName string `regex:"[A-Za-z]{2,255}"`	
	Password string `regex:".{8,255}"`
}

func main() {
	usr := User {
		FirstName: "Jim",
		LastName: "Don",
		Password: "qwerty"
    }
	
    if err := ValidateStruct(usr); err != nil {
		log.Pringln(err) // "User has to have valid Password according to pattern: `.{8,255}`"
    }
}
```

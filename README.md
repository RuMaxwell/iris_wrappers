# iris_wrappers

A wrapper library for Iris ([](https://github.com/kataras/iris)) web framework.

iris_wrappers contains several (currently only 1) packages to simplify you work with Iris.

## Package: api

api is a wrapper of Iris' application providing request handlers that take more intuitive handler functions in the case of API call-and-response with JSON data.

Usally, we write such a handler function like this:

```go
type Login struct {
	Name string `json:name`
	Password string `json:pswd`
}

// in func main
app.Post("/api/login", func(ctx iris.Context) {
	login := Login {}
	if err := ctx.ReadJSON(&login); err != nil {
		// ...
	}

	// generate UID and token...

	ctx.JSON(map[string]interface{}{
		"id": uid,
		"token": token,
	})
})
```

There are some things we must write every time. Such as `ctx iris.Context`, `if err := ctx.ReadJSON(&login)...`, and `ctx.JSON`. Our wrapper wraps this behavior and now all you have to do is:

```go
import "github.com/RuMaxwell/iris_wrappers/api"

type Login struct {
	Name string `json:name`
	Password string `json:pswd`
}

type LoginResponse struct {
	UID uint32 `json:id`
	Token []byte `json:token`
}

// in func main
app := iris.New()
wapp := api.New(app)

wapp.Post("/api/login", func(login *Login) LoginResponse {
	// generate UID and token...

	return LoginResponse {
		UID: uid,
		Token: token,
	}
})
```

The api wrapper makes this scenario easier to write.

> Note: As in [this website](https://iris-go.com/start/#dependency-injection), they say the new release of Iris provides this feature in hero package. But I can't find them and use them.

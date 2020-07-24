package api

import (
	"fmt"
	"reflect"

	"github.com/kataras/iris/v12"
)

func wrapHandler(handler interface{}) iris.Handler {
    return func(ctx iris.Context) {
        // user
        inType := reflect.TypeOf(handler).In(0).Elem()

        // *user
        in := reflect.New(inType)

        data := make(map[string]interface{})

        if err := ctx.ReadJSON(&data); err != nil {
            fmt.Println(err.Error())
            ctx.StatusCode(400)
            ctx.WriteString("invalid request form")
            return
        }
        
        // brush `data` onto `in`
        for i := 0; i < in.Elem().NumField(); i++ {
            // field = user.UID uint32 `json:"id"`
            field := inType.Field(i)
            // jsonName = "id"
            jsonName, ok := field.Tag.Lookup("json")
            if !ok {
                jsonName = field.Name
            }

            // inString = data["id"]
            inString, ok := data[jsonName]
            if (!ok) {
                fmt.Printf("WARN: incoming value has a field '%s' that receiving struct not expected\n", jsonName)
            }

            // inStringRV = (data["id"]).convertTo(uint32)
            inStringRV := reflect.ValueOf(inString)
            inRV := inStringRV.Convert(field.Type)
            in.Elem().FieldByName(field.Name).Set(inRV)
        }

        ins := []reflect.Value{ in }


        outType := reflect.TypeOf(handler).Out(0)

        out := reflect.ValueOf(handler).Call(ins)

        if outType.Name() == "" {
            ctx.JSON(out[0])
            return
        }

        outData := make(map[string]interface{})

        for i := 0; i < out[0].NumField(); i++ {
            outData[outType.Field(i).Name] = out[0].Field(i).Interface()
        }

        ctx.JSON(outData)
    }
}

func wrapHandlers(handlers []interface{}) []iris.Handler {
    wrappedHandlers := make([]iris.Handler, len(handlers))

    for i, h := range handlers {
        wrappedHandlers[i] = wrapHandler(h)
    }

    return wrappedHandlers
}

// Application wraps `iris.Application` for use of simpler API handler functions.
// API handler functions must be like: func(in InStruct) OutStruct. After wrapping it seems like
// func(ctx iris.Context) { ctx.ReadJSON(&in); ... ctx.JSON(&out); }
type Application struct {
    app *iris.Application
}

// New creates a wrapper application from a real Iris application.
func New(app *iris.Application) Application {
    return Application{
        app: app,
    }
}

// Get registers a route for the GET HTTP method.
func (app *Application) Get(url string, handlers ...interface{}) {
    wrappedHandlers := wrapHandlers(handlers)

    app.app.Get(url, wrappedHandlers...)
}

// Post registers a route for the POST HTTP method.
func (app *Application) Post(url string, handlers ...interface{}) {
    wrappedHandlers := wrapHandlers(handlers)

    app.app.Post(url, wrappedHandlers...)
}

// Put registers a route for the PUT HTTP method.
func (app *Application) Put(url string, handlers ...interface{}) {
    wrappedHandlers := wrapHandlers(handlers)

    app.app.Put(url, wrappedHandlers...)
}

// Delete registers a route for the DELETE HTTP method.
func (app *Application) Delete(url string, handlers ...interface{}) {
    wrappedHandlers := wrapHandlers(handlers)

    app.app.Delete(url, wrappedHandlers...)
}

// Patch registers a route for the PATCH HTTP method.
func (app *Application) Patch(url string, handlers ...interface{}) {
    wrappedHandlers := wrapHandlers(handlers)

    app.app.Patch(url, wrappedHandlers...)
}

// Head registers a route for the HEAD HTTP method.
func (app *Application) Head(url string, handlers ...interface{}) {
    wrappedHandlers := wrapHandlers(handlers)

    app.app.Head(url, wrappedHandlers...)
}

// Options registers a route for the OPTIONS HTTP method.
func (app *Application) Options(url string, handlers ...interface{}) {
    wrappedHandlers := wrapHandlers(handlers)

    app.app.Options(url, wrappedHandlers...)
}

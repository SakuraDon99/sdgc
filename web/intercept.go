package web

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
)

func Intercept(f any) gin.HandlerFunc {
	ft := reflect.TypeOf(f)
	fv := reflect.ValueOf(f)

	fName := runtime.FuncForPC(fv.Pointer()).Name()
	if !validateHandlerFunc(ft) {
		panic(fmt.Errorf("wrong handler func: %s", fName))
	}

	ftNumIn := ft.NumIn()

	successHTTPCode := http.StatusOK
	var binders []binder
	var err error
	if ft.NumIn() == 2 {
		p := reflect.New(ft.In(1))
		pe := p.Elem()
		for i := 0; i < pe.NumField(); i++ {
			tag := pe.Type().Field(i).Tag
			if path := tag.Get("path"); path != "" {
				binders = append(binders, &pathBinder{
					field: i,
					path:  path,
				})
			} else if query := tag.Get("query"); query != "" {
				binders = append(binders, &queryBinder{
					field: i,
					query: query,
				})
			} else if request := tag.Get("request"); request != "" {
				if request == "json" {
					binders = append(binders, &requestJSONBinder{
						field: i,
						body:  pe.Field(i).Type(),
					})
				}
			} else if code := tag.Get("success"); code != "" {
				successHTTPCode, err = strconv.Atoi(code)
				if err != nil || successHTTPCode < 100 || successHTTPCode > 999 {
					panic(fmt.Errorf("invalid http error code: %s in func %s", code, fName))
				}
			}
		}
	}

	return func(c *gin.Context) {
		incomes := make([]reflect.Value, ftNumIn)
		incomes[0] = reflect.ValueOf(c)

		if ftNumIn == 2 {
			param := reflect.New(ft.In(1))
			paramElem := param.Elem()
			for _, b := range binders {
				err := b.Bind(c, paramElem)
				if err != nil {
					handleError(c, err)
					return
				}
			}

			incomes[1] = paramElem
		}

		outcomes := fv.Call(incomes)

		// if handler does not return parameters, return client success
		if len(outcomes) == 0 {
			c.Status(successHTTPCode)
			return
		}

		response := outcomes[0]
		errResponse := outcomes[len(outcomes)-1]

		if err, ok := errResponse.Interface().(error); ok {
			handleError(c, err)
			return
		}

		handleResponse(c, response, successHTTPCode)
	}
}

func validateHandlerFunc(ft reflect.Type) bool {
	if ft.Kind() != reflect.Func || ft.NumIn() < 1 || ft.NumIn() > 2 || ft.NumOut() > 2 {
		return false
	}
	// the first income parameter must be context
	if ft.In(0).Name() != "Context" {
		return false
	}
	// if handler return 2 parameters, the second parameter must be error
	if ft.NumOut() == 2 && ft.Out(1).Name() != "error" {
		return false
	}

	return true
}

func setVal(fieldValue reflect.Value, val string) error {
	switch fieldValue.Interface().(type) {
	case string:
		fieldValue.SetString(val)
	case uint64:
		uint64Val, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return ErrInvalidParams
		}
		fieldValue.SetUint(uint64Val)
	}
	return nil
}

func handleError(c *gin.Context, err error) {
	c.String(httpErrorCode(err), err.Error())
	c.Abort()
}

func httpErrorCode(err error) int {
	if errors.Is(err, ErrInvalidParams) {
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}

func handleResponse(c *gin.Context, response reflect.Value, statusCode int) {
	if (response.Kind() == reflect.Ptr || response.Kind() == reflect.Interface) && response.IsNil() {
		c.Status(statusCode)
		return
	}
	if response.Kind() == reflect.Slice && response.IsNil() {
		c.JSON(statusCode, make([]any, 0))
		return
	}

	c.JSON(statusCode, response.Interface())
}

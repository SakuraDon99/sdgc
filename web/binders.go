package core

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

type binder interface {
	Bind(c *gin.Context, param reflect.Value) error
}

type pathBinder struct {
	field int
	path  string
}

func (b *pathBinder) Bind(c *gin.Context, param reflect.Value) error {
	val := c.Param(b.path)
	err := setVal(param.Field(b.field), val)
	return err
}

type queryBinder struct {
	field int
	query string
}

func (b *queryBinder) Bind(c *gin.Context, param reflect.Value) error {
	val := c.Query(b.query)
	err := setVal(param.Field(b.field), val)
	return err
}

type bodyBinder struct {
	field int
	body  reflect.Type
}

func (b *bodyBinder) Bind(c *gin.Context, param reflect.Value) error {
	body := reflect.New(b.body).Interface()
	err := c.ShouldBindJSON(body)
	if err != nil {
		return err
	}
	if validator, ok := body.(Validator); ok {
		err = validator.Validate()
		if err != nil {
			return err
		}
	}
	param.Field(b.field).Set(reflect.ValueOf(body).Elem())
	return nil
}

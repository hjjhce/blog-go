package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testHandler = func(ctx *Context) { ctx.Next() }

func TestRouterGroup(t *testing.T) {
	core := New()
	a := core.Group("/a")
	b := core.Group("/b")

	a.GET("/", func(ctx *Context) {
		ctx.String(200, "a")
	})

	a.GET("/hello", func(ctx *Context) {
		ctx.String(200, "a/hello")
	})

	b.GET("/", func(ctx *Context) {
		ctx.String(200, "b")
	})

	r, _ := http.NewRequest("GET", "/a/", nil)
	w := httptest.NewRecorder()
	core.ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code, "should be equal")
	assert.Equal(t, "a", w.Body.String(), "should be equal")

	r, _ = http.NewRequest("GET", "/a/hello/", nil)
	w = httptest.NewRecorder()
	core.ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code, "should be equal")
	assert.Equal(t, "a/hello", w.Body.String(), "should be equal")

	r, _ = http.NewRequest("GET", "/b/", nil)
	w = httptest.NewRecorder()
	core.ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "b", w.Body.String(), "should be equal")

}

func BenchmarkRouterGroup(t *testing.B) {
	core := New()
	a := core.Group("/a")
	b := core.Group("/b")

	a.GET("/", func(ctx *Context) {
		ctx.String(200, "a")
	})

	a.GET("/hello", func(ctx *Context) {
		ctx.String(200, "a/hello")
	})

	b.GET("/", func(ctx *Context) {
		ctx.String(200, "b")
	})

	r, _ := http.NewRequest("GET", "/a/", nil)
	w := httptest.NewRecorder()
	core.ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code, "should be equal")
	assert.Equal(t, "a", w.Body.String(), "should be equal")

	r, _ = http.NewRequest("GET", "/a/hello/", nil)
	w = httptest.NewRecorder()
	core.ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code, "should be equal")
	assert.Equal(t, "a/hello", w.Body.String(), "should be equal")

	r, _ = http.NewRequest("GET", "/b/", nil)
	w = httptest.NewRecorder()
	core.ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "b", w.Body.String(), "should be equal")
}

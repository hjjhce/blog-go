package main

type RouterGroup struct {
}

type HandlerFunc func(*Context)

func (rg *RouterGroup) GET(path string, hanlders ...HandlerFunc) {}

package controllers

import "github.com/robfig/revel"

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Register() revel.Result {
	return c.Render()
}

func (c App) Upload() revel.Result {
	return c.Render()
}
func (c App) Download() revel.Result {
	return c.Render()
}

func (c App) List() revel.Result {
	return c.Render()
}

func (c App) Login() revel.Result {
	return c.Render()
}

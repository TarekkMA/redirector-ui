package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/TarekkMA/redirector-ui/datasource"

	"github.com/TarekkMA/redirector-ui/data"
)

type HomeController struct {
	DS datasource.RedirectsDataSource
}

type HomeData struct {
	Redirects    []*data.Redirect
	RedirectsErr string
}

func (c *HomeController) HomePage(w http.ResponseWriter, r *http.Request) {
	if CheckUserAndRespond(w, r) == false {
		return
	}
	red, err := c.DS.GetAll()
	homepage(&HomeData{
		Redirects:    red,
		RedirectsErr: fmt.Sprintf("+%v", err),
	}, w, r)
}

func (c *HomeController) HomeAction(w http.ResponseWriter, r *http.Request) {
	if CheckUserAndRespond(w, r) == false {
		return
	}

	r.ParseForm()
	requesttype := r.Form["type"]
	from := r.Form["from"]
	to := r.Form["to"]

	if requesttype[0] == "delete" {
		c.DS.RemoveItem(&data.Redirect{From: from[0], To: to[0]})
		c.HomePage(w, r)
		return
	}

	if requesttype[0] == "add" {
		c.DS.AddItem(&data.Redirect{From: from[0], To: to[0]})
		c.HomePage(w, r)
		return
	}

}

func homepage(data *HomeData, w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("tmpls/index.html")
	if err != nil {
		log.Println(err)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

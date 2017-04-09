package controllers

import (
	"rei-site/app/models"

	"github.com/revel/revel"
)

type Api struct {
	*revel.Controller
	GorpController
}

func (c Api) Test() revel.Result {
	results, err := c.Txn.Select(models.User{},
		`select * from User`)
	if err != nil {
		panic(err)
	}

	var users []*models.User
	for _, r := range results {
		u := r.(*models.User)
		users = append(users, u)
	}

	return c.RenderJSON(users)
}

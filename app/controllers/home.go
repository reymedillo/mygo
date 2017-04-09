package controllers

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/revel/revel"

	"rei-site/app/routes"
)

// Home - main struct in this controller
type Home struct {
	GorpController
}

// User - use user controllers functions
type User struct {
	UserController
}

// AddUser - route
func (c User) AddUser() revel.Result {
	if user := c.connected(); user != nil {
		c.ViewArgs["user"] = user
	}
	return nil
}

// Index - route
func (c User) Index() revel.Result {
	if c.connected() != nil {
		return c.Redirect(routes.App.Login())
	}
	c.Flash.Error("Please log in first")
	return c.Render()
}

// Register - route
func (c Home) Register() revel.Result {
	return c.Render()
}

// Login - route
func (c User) Login(username, password string, remember bool) revel.Result {
	user := c.getUser(username)
	if user != nil {
		err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
		if err == nil {
			c.Session["user"] = username
			if remember {
				c.Session.SetDefaultExpiration()
			} else {
				c.Session.SetNoExpiration()
			}
			c.Flash.Success("Welcome, " + username)
			return c.Redirect(routes.Home.Index())
		}
	}

	c.Flash.Out["username"] = username
	c.Flash.Error("Login failed")
	return c.Redirect(routes.App.Login())
}

// Logout - route
func (c Home) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.App.Login())
}

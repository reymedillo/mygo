package controllers

import (
	"rei-site/app/models"
	"rei-site/app/routes"

	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

// UserController - import/use Gorp
type UserController struct {
	*revel.Controller
	GorpController
}

// AddUser - to be called to create new user
func (c UserController) AddUser() revel.Result {
	if user := c.connected(); user != nil {
		c.ViewArgs["user"] = user
	}
	return nil
}

func (c UserController) connected() *models.User {
	if c.ViewArgs["user"] != nil {
		return c.ViewArgs["user"].(*models.User)
	}
	if username, ok := c.Session["user"]; ok {
		return c.getUser(username)
	}
	return nil
}

func (c UserController) getUser(username string) *models.User {
	users, err := c.Txn.Select(models.User{}, `select * from User where Username = ?`, username)
	if err != nil {
		panic(err)
	}
	if len(users) == 0 {
		return nil
	}
	return users[0].(*models.User)
}

// SaveUser - save user to database
func (c UserController) SaveUser(user models.User, verifyPassword string) revel.Result {
	c.Validation.Required(verifyPassword)
	c.Validation.Required(verifyPassword == user.Password).Message("Passwords do not match.")
	user.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.App.Register())
	}

	user.HashedPassword, _ = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	err := c.Txn.Insert(&user)
	if err != nil {
		panic(err)
	}

	c.Session["user"] = user.Username
	c.Flash.Success("Welcome, " + user.Name)
	return c.Redirect(routes.Home.Index())
}

func (c UserController) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please login first.")
		return c.Redirect(routes.App.Login())
	}
	return nil
}

package routes

import (
	"github.com/create-go-app/fiber-go-template/app/controllers"
	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	// Routes for GET method:
	// route.Get("/shortUrl", controllers.CreateShortUrl) // get list of all books

	// Routes for POST method:
	route.Post("/short-url", controllers.CreateShortUrl) // get list of all books

	// route.Post("/user/sign/in", controllers.UserSignIn) // auth, return Access & Refresh tokens
}

package user_routes

import (
	"net/http"

	usercontroller "github.com/http-crud/api/controllers"
	user_middleware "github.com/http-crud/api/middlewares"
)

func UserRoutes(mux *http.ServeMux) {
	// This line of code is registering a route for the "/user/register" endpoint on the provided `mux`
	// ServeMux. It is also adding middleware to the route using `user_middleware.RegisterUserMiddleware`
	// and specifying the handler function for the route as `usercontroller.RegisterUserHandler`. This
	// means that when a request is made to the "/user/register" endpoint, it will first go through the
	// middleware before being handled by the `RegisterUserHandler` function.
	// POST
	mux.Handle("/user/register", user_middleware.RegisterUserMiddleware(http.HandlerFunc(usercontroller.RegisterUserHandler)))

	// This line of code is registering a route for the "/user/login" endpoint on the provided `mux`
	// ServeMux. It is also adding middleware to the route using `user_middleware.LoginUserMiddleware` and
	// specifying the handler function for the route as `usercontroller.LoginUserHandler`. This means that
	// when a request is made to the "/user/login" endpoint, it will first go through the middleware before
	// being handled by the `LoginUserHandler` function. The middleware is responsible for performing any
	// necessary checks or operations before the request is handled by the handler function.
	// POST
	mux.Handle("/user/login", user_middleware.LoginUserMiddleware(http.HandlerFunc(usercontroller.LoginUserHandler)))

	// This line of code is registering a route for the "/user/" endpoint on the provided `mux` ServeMux.
	// It is also adding middleware to the route using `user_middleware.GetUserMiddleware` and specifying
	// the handler function for the route as `usercontroller.GetUserHandler`. This means that when a
	// request is made to the "/user/" endpoint, it will first go through the middleware before being
	// handled by the `GetUserHandler` function. The middleware is responsible for performing any necessary
	// checks or operations before the request is handled by the handler function.
	// GET
	mux.Handle("/user/", user_middleware.GetUserMiddleware(http.HandlerFunc(usercontroller.GetUserHandler)))

	// This line of code is registering a route for the "/user/update" endpoint on the provided `mux`
	// ServeMux. It is also adding middleware to the route using `user_middleware.GetUserMiddleware` and
	// specifying the handler function for the route as `usercontroller.UpdateUserHandler`. This means
	// that when a request is made to the "/user/update" endpoint, it will first go through the middleware
	// before being handled by the `UpdateUserHandler` function. The middleware is responsible for
	// performing any necessary checks or operations before the request is handled by the handler
	// function.
	// PATCH
	mux.Handle("/user/update", user_middleware.GetUserMiddleware(http.HandlerFunc(usercontroller.UpdateUserHandler)))

	// This line of code is registering a route for the "/user/delete" endpoint on the provided `mux`
	// ServeMux. It is also adding middleware to the route using `user_middleware.GetUserMiddleware` and
	// specifying the handler function for the route as `usercontroller.DeletUserHandler`. This means that
	// when a request is made to the "/user/delete" endpoint, it will first go through the middleware
	// before being handled by the `DeletUserHandler` function. The middleware is responsible for
	// performing any necessary checks or operations before the request is handled by the handler function.
	// DELETE
	mux.Handle("/user/delete", user_middleware.GetUserMiddleware(http.HandlerFunc(usercontroller.DeletUserHandler)))
}

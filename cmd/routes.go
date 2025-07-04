package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
	// httpSwagger "github.com/swaggo/http-swagger"
	// _ "naimuBack/docs"
)

func (app *application) JWTMiddlewareWithRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return app.JWTMiddleware(next, requiredRole)
	}
}

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders, makeResponseJSON)
	//authMiddleware := standardMiddleware.Append(app.JWTMiddlewareWithRole("user"))
	adminAuthMiddleware := standardMiddleware.Append(app.JWTMiddlewareWithRole("admin"))
	trainerAuthMiddleware := standardMiddleware.Append(app.JWTMiddlewareWithRole("trainer"))
	clientAuthMiddleware := standardMiddleware.Append(app.JWTMiddlewareWithRole("client"))

	mux := pat.New()

	// Users
	mux.Post("/user", adminAuthMiddleware.ThenFunc(app.userHandler.CreateUser))
	mux.Post("/user/sign_up", standardMiddleware.ThenFunc(app.userHandler.SignUp))
	mux.Post("/user/sign_in", standardMiddleware.ThenFunc(app.userHandler.SignIn))
	mux.Post("/user/upgrade", clientAuthMiddleware.ThenFunc(app.userHandler.UpgradeToTrainer))

	// Programs
	mux.Post("/program", trainerAuthMiddleware.ThenFunc(app.programHandler.CreateProgram))
	mux.Get("/programs", standardMiddleware.ThenFunc(app.programHandler.ProgramsByTrainer))

	// Exercises and Food
	mux.Post("/exercise", trainerAuthMiddleware.ThenFunc(app.exerciseHandler.CreateExercise))
	mux.Post("/food", trainerAuthMiddleware.ThenFunc(app.foodHandler.CreateFood))

	// Days
	mux.Get("/program/day", standardMiddleware.ThenFunc(app.dayHandler.DayDetails))
	mux.Post("/program/day/complete", standardMiddleware.ThenFunc(app.dayHandler.CompleteDay))
	mux.Post("/program/day", trainerAuthMiddleware.ThenFunc(app.dayHandler.CreateDay))
	mux.Post("/program/day/complete", standardMiddleware.ThenFunc(app.dayHandler.CompleteDay))

	// mux.Get("/swagger/", httpSwagger.WrapHandler)

	return standardMiddleware.Then(mux)
}

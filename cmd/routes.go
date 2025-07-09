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
	mux.Get("/programs", trainerAuthMiddleware.ThenFunc(app.programHandler.ProgramsByTrainer))
	mux.Get("/program/:id", trainerAuthMiddleware.ThenFunc(app.programHandler.GetProgram))
	mux.Put("/program/:id", trainerAuthMiddleware.ThenFunc(app.programHandler.UpdateProgram))
	mux.Del("/program/:id", trainerAuthMiddleware.ThenFunc(app.programHandler.DeleteProgram))

	// Clients
	mux.Get("/clients", trainerAuthMiddleware.ThenFunc(app.userHandler.GetAllClients))
	mux.Get("/program/:program_id/clients", trainerAuthMiddleware.ThenFunc(app.userHandler.GetClientsByProgramID))
	mux.Del("/program/:program_id/client/:client_id", trainerAuthMiddleware.ThenFunc(app.userHandler.DeleteClientFromProgram))
	mux.Get("/client/:client_id/programs", trainerAuthMiddleware.ThenFunc(app.userHandler.GetProgramsByClientID))

	// Exercises and Food
	mux.Post("/exercise", trainerAuthMiddleware.ThenFunc(app.exerciseHandler.CreateExercise))
	mux.Put("/exercise/:id", trainerAuthMiddleware.ThenFunc(app.exerciseHandler.UpdateExercise))
	mux.Del("/exercise/:id", trainerAuthMiddleware.ThenFunc(app.exerciseHandler.DeleteExercise))
	mux.Post("/food", trainerAuthMiddleware.ThenFunc(app.foodHandler.CreateFood))
	mux.Put("/food/:id", trainerAuthMiddleware.ThenFunc(app.foodHandler.UpdateFood))
	mux.Del("/food/:id", trainerAuthMiddleware.ThenFunc(app.foodHandler.DeleteFood))

	// Days
	mux.Get("/program/:program_id/days", trainerAuthMiddleware.ThenFunc(app.dayHandler.DaysByProgram))
	mux.Get("/program/:program_id/day/:day", trainerAuthMiddleware.ThenFunc(app.dayHandler.DayDetails))
	mux.Post("/program/day/complete", standardMiddleware.ThenFunc(app.dayHandler.CompleteDay))
	mux.Post("/program/day/food", standardMiddleware.ThenFunc(app.dayHandler.CompleteFood))
	mux.Post("/program/day/exercise", standardMiddleware.ThenFunc(app.dayHandler.CompleteExercise))
	mux.Get("/program/day/progress", standardMiddleware.ThenFunc(app.dayHandler.ProgressStatus))
	mux.Get("/program/:program_id/progress", standardMiddleware.ThenFunc(app.dayHandler.ProgramProgress))
	mux.Post("/program/day", trainerAuthMiddleware.ThenFunc(app.dayHandler.CreateDay))

	mux.Put("/program/day/:id", trainerAuthMiddleware.ThenFunc(app.dayHandler.UpdateDay))
	mux.Del("/program/day/:id", trainerAuthMiddleware.ThenFunc(app.dayHandler.DeleteDay))

	mux.Put("/program/day/:id", trainerAuthMiddleware.ThenFunc(app.dayHandler.UpdateDay))

	// mux.Get("/swagger/", httpSwagger.WrapHandler)

	return standardMiddleware.Then(mux)
}

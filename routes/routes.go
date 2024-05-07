package routes

import (
	"backend-pedika-fiber/handlers"
	"backend-pedika-fiber/middleware"

	"github.com/gofiber/fiber/v2"
)

/*========= || Endpoint yang hanya bisa diakses oleh admin || ====================*/
func SetAdminRoutes(app *fiber.App) {
	adminGroup := app.Group("/api/admin")
	adminGroup.Use(middleware.AdminMiddleware)

	adminGroup.Get("/profile", handlers.GetUserProfile)
	adminGroup.Put("/edit-profile", handlers.UpdateUserProfile)

	adminGroup.Get("/violence-categories", handlers.GetAllViolenceCategories)
	adminGroup.Post("/create-violence-category", handlers.CreateViolenceCategory)
	adminGroup.Put("/violence-category/:id", handlers.UpdateViolenceCategory)
	adminGroup.Delete("/delete-violence-category/:id", handlers.DeleteViolenceCategory)

	adminGroup.Get("/emergency-contact", handlers.GetEmergencyContact)
	adminGroup.Put("/emergency-contact-edit", handlers.UpdateEmergencyContact)

	adminGroup.Get("/contents", handlers.GetAllContents)
	adminGroup.Post("/create-content", handlers.CreateContent)
	adminGroup.Put("/edit-content/:id", handlers.UpdateContent)
	adminGroup.Delete("/delete-content/:id", handlers.DeleteContent)

	adminGroup.Get("/event", handlers.GetAllEvent)
	adminGroup.Post("/create-event", handlers.CreateEvent)
	adminGroup.Put("/edit-event/:id", handlers.UpdateEvent)
	adminGroup.Delete("/delete-event/:id", handlers.DeleteEvent)

	adminGroup.Get("/laporan", handlers.GetLatestReports)
	adminGroup.Get("/detail-laporan/:no_registrasi", handlers.GetLaporanByNoRegistrasi)
}

/*========= ||  Endpoint yang hanya bisa diakses oleh masyarakat || ====================*/
func SetMasyarakatRoutes(app *fiber.App) {
	masyarakatGroup := app.Group("/api/masyarakat")
	masyarakatGroup.Use(middleware.MasyarakatMiddleware)

	masyarakatGroup.Get("/profile", handlers.GetUserProfile)
	masyarakatGroup.Put("/edit-profile", handlers.UpdateUserProfile)

	masyarakatGroup.Get("/kategori-kekerasan", handlers.GetAllViolenceCategories)
	masyarakatGroup.Get("/kategori-kekerasan/:id", handlers.GetViolenceCategoryByID)

	masyarakatGroup.Post("/buat-laporan", handlers.CreateLaporan)
	masyarakatGroup.Get("/laporan", handlers.GetUserReports)
	masyarakatGroup.Get("/detail-laporan/:no_registrasi", handlers.GetReportByNoRegistrasi)

	masyarakatGroup.Post("/pelaku-kekerasan", handlers.CreatePelaku)
	masyarakatGroup.Put("/pelaku-kekerasan/:id", handlers.UpdatePelaku)
	masyarakatGroup.Delete("/pelaku-kekerasan/:id", handlers.DeletePelaku)

	masyarakatGroup.Get("/content", handlers.GetAllContents)
}

/*========= ||  Endpoint bisa di akses tanpa login || ====================*/
func RoutesWithOutLogin(app *fiber.App) {
	app.Get("/api/emergency-contact", handlers.EmergencyContact)
	app.Get("/hello", handlers.HelloMasyarakat)
}

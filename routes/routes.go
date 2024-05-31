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
	adminGroup.Post("/change-password", handlers.ChangePassword)

	adminGroup.Get("/emergency-contact", handlers.GetEmergencyContact)
	adminGroup.Put("/emergency-contact-edit", handlers.UpdateEmergencyContact)
	adminGroup.Get("/emergency-contact/:id", handlers.ShowEmergencyContactByID)

	adminGroup.Get("/laporans", handlers.GetLatestReports)
	adminGroup.Get("/detail-laporan/:no_registrasi", handlers.GetLaporanByNoRegistrasi)
	adminGroup.Put("/lihat-laporan/:no_registrasi", handlers.AdminLihatLaporan)
	adminGroup.Put("/proses-laporan/:no_registrasi", handlers.AdminProsesLaporan)

	adminGroup.Post("/create-tracking-laporan", handlers.CreateTrackingLaporan)
	adminGroup.Delete("/delete-tracking-laporan/:id", handlers.DeleteTrackingLaporan)
	adminGroup.Put("/edit-tracking-laporan/:id", handlers.UpdateTrackingLaporan)

	adminGroup.Post("/create-pelaku-kekerasan", handlers.CreatePelaku)
	adminGroup.Put("/edit-pelaku-kekerasan/:id", handlers.UpdatePelaku)

	adminGroup.Post("/create-korban-kekerasan", handlers.CreateKorban)
	adminGroup.Put("/edit-korban-kekerasan/:id", handlers.UpdateKorban)

	adminGroup.Get("/violence-categories", handlers.GetAllViolenceCategories)
	adminGroup.Get("/detail-violence-category/:id", handlers.GetViolenceCategoryByID)
	adminGroup.Post("/create-violence-category", handlers.CreateViolenceCategory)
	adminGroup.Put("/edit-violence-category/:id", handlers.UpdateViolenceCategory)
	adminGroup.Delete("/delete-violence-category/:id", handlers.DeleteViolenceCategory)

	adminGroup.Get("/contents", handlers.GetAllContents)
	adminGroup.Get("/detail-content/:id", handlers.GetContentByID)
	adminGroup.Post("/create-content", handlers.CreateContent)
	adminGroup.Put("/edit-content/:id", handlers.UpdateContent)
	adminGroup.Delete("/delete-content/:id", handlers.DeleteContent)

	adminGroup.Get("/event", handlers.GetAllEvent)
	adminGroup.Get("/detail-event/:id", handlers.GetEventByID)
	adminGroup.Post("/create-event", handlers.CreateEvent)
	adminGroup.Put("/edit-event/:id", handlers.UpdateEvent)
	adminGroup.Delete("/delete-event/:id", handlers.DeleteEvent)

	adminGroup.Get("/janjitemus", handlers.AdminGetAllJanjiTemu)
	adminGroup.Get("/detail-janjitemu/:id", handlers.AdminJanjiTemuByID)
	adminGroup.Put("/approve-janjitemu/:id", handlers.AdminApproveJanjiTemu)
	adminGroup.Put("/cancel-janjitemu/:id", handlers.AdminCancelJanjiTemu)

}

/*========= ||  Endpoint yang hanya bisa diakses oleh masyarakat || ====================*/
func SetMasyarakatRoutes(app *fiber.App) {
	masyarakatGroup := app.Group("/api/masyarakat")
	masyarakatGroup.Use(middleware.MasyarakatMiddleware)

	masyarakatGroup.Get("/profile", handlers.GetUserProfile)
	masyarakatGroup.Put("/edit-profile", handlers.UpdateUserProfile)
	masyarakatGroup.Put("/change-password", handlers.ChangePassword)

	masyarakatGroup.Get("/kategori-kekerasan", handlers.GetAllViolenceCategories)
	masyarakatGroup.Get("/kategori-kekerasan/:id", handlers.GetViolenceCategoryByID)

	masyarakatGroup.Get("/laporans", handlers.GetUserReports)
	masyarakatGroup.Post("/buat-laporan", handlers.CreateLaporan)
	masyarakatGroup.Put("/edit-laporan/:no_registrasi", handlers.EditLaporan)
	masyarakatGroup.Get("/detail-laporan/:no_registrasi", handlers.GetReportByNoRegistrasi)
	masyarakatGroup.Put("batalkan-laporan/:no_registrasi", handlers.BatalkanLaporan)

	masyarakatGroup.Post("/create-korban-kekerasan", handlers.CreateKorban)
	masyarakatGroup.Put("/edit-korban-kekerasan/:id", handlers.UpdateKorban)

	masyarakatGroup.Post("/create-pelaku-kekerasan", handlers.CreateKorban)
	masyarakatGroup.Put("/edit-pelaku-kekerasan/:id", handlers.UpdateKorban)

	masyarakatGroup.Post("/create-korban-kekerasan", handlers.CreateKorban)
	masyarakatGroup.Put("/edit-korban-kekerasan/:id", handlers.UpdateKorban)

	masyarakatGroup.Get("/janjitemus", handlers.GetUserJanjiTemus)
	masyarakatGroup.Get("/detail-janjitemu/:id", handlers.GetJanjiTemuByID)
	masyarakatGroup.Post("/create-janjitemu", handlers.MasyarakatCreateJanjiTemu)
	masyarakatGroup.Put("/edit-janjitemu/:id", handlers.MasyarakatEditJanjiTemu)
	masyarakatGroup.Put("/batal-janjitemu/:id", handlers.MasyarakatCancelJanjiTemu)

	masyarakatGroup.Get("/content", handlers.GetAllContents)
	masyarakatGroup.Get("/detail-content/:id", handlers.GetContentByID)
}

/*========= ||  Endpoint bisa di akses tanpa login || ====================*/
func RoutesWithOutLogin(app *fiber.App) {
	app.Get("/api/emergency-contact", handlers.EmergencyContact)
	app.Get("/api/publik-content", handlers.GetAllContents)
	app.Get("/api/detail-content/:id", handlers.GetContentByID)
	app.Get("/hello", handlers.HelloMasyarakat)

}

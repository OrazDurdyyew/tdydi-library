package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"library/internal/auth"
	"library/internal/handlers"
	"library/internal/middleware"
	"library/internal/services"
)

func Register(app *fiber.App, db *gorm.DB, tm *auth.TokenManager, files *services.FileStorage, storagePath string) {
	perms := services.NewPermissionsService(db)

	healthH := handlers.NewHealthHandler(db)
	authH := handlers.NewAuthHandler(db, tm)
	menuH := handlers.NewMenuHandler(db)
	sectionsH := handlers.NewSectionsHandler(db)
	settingsH := handlers.NewSiteSettingsHandler(db, files)
	booksH := handlers.NewBooksHandler(db, perms, files)
	filesH := handlers.NewFilesHandler(storagePath)
	periodicalsH := handlers.NewPeriodicalsHandler(db, files)
	collectionsH := handlers.NewCollectionsHandler(db, files)
	postsH := handlers.NewPostsHandler(db, files)
	pagesH := handlers.NewPagesHandler(db)
	faqH := handlers.NewFAQHandler(db)
	homeH := handlers.NewHomeHandler(db, files)
	adminsH := handlers.NewAdminsHandler(db)
	statsH := handlers.NewStatsHandler(db, perms)
	searchH := handlers.NewSearchHandler(db)
	contentUploadsH := handlers.NewContentUploadsHandler(files)
	departmentsH := handlers.NewDepartmentsHandler(db)
	teachersH := handlers.NewTeachersHandler(db, files)
	pubCategoriesH := handlers.NewPublicationCategoriesHandler(db)
	publicationsH := handlers.NewPublicationsHandler(db, files)
	pubStatsH := handlers.NewPublicationStatsHandler(db)

	app.Get("/api/health", healthH.Get)

	// OptionalAuth — на всей группе /api: публичные GET'ы видят залогиненного
	// админа (черновики постов/книг доступны при редактировании в админке).
	api := app.Group("/api", middleware.OptionalAuth(tm))

	// ── Auth ──────────────────────────────────────────────────────────────
	api.Post("/auth/login", authH.Login)

	// ── Site settings (public) ────────────────────────────────────────────
	api.Get("/site-settings", settingsH.Get)

	// ── Menu ──────────────────────────────────────────────────────────────
	api.Get("/menu", menuH.Get)

	// ── Home ──────────────────────────────────────────────────────────────
	api.Get("/home", homeH.GetHomeData)

	// ── Sections ──────────────────────────────────────────────────────────
	api.Get("/sections", sectionsH.List)
	api.Get("/sections/slug/:slug", sectionsH.GetBySlug)
	api.Get("/sections/:id", sectionsH.Get)

	// ── Books ─────────────────────────────────────────────────────────────
	api.Get("/books", booksH.List)
	api.Get("/books/slug/:slug", booksH.GetBySlug)
	api.Get("/books/:id", booksH.Get)
	api.Get("/books/:id/similar", booksH.Similar)
	api.Get("/books/:id/file", booksH.ServeFile)
	api.Get("/books/:id/download", booksH.Download)

	// ── Periodicals ───────────────────────────────────────────────────────
	api.Get("/periodicals", periodicalsH.List)
	api.Get("/periodicals/slug/:slug", periodicalsH.GetBySlug)
	api.Get("/periodicals/issues", periodicalsH.ListAllIssues) // before /:id
	api.Get("/periodicals/:id", periodicalsH.Get)
	api.Get("/periodicals/:id/issues", periodicalsH.ListIssues)
	api.Get("/periodicals/:id/issue-years", periodicalsH.IssueYears)
	api.Get("/periodicals/issues/:id", periodicalsH.GetIssue)
	api.Get("/periodicals/issues/:id/file", periodicalsH.ServeIssueFile)
	api.Get("/periodicals/issues/:id/download", periodicalsH.DownloadIssue)

	// ── Collections ───────────────────────────────────────────────────────
	api.Get("/collections", collectionsH.List)
	api.Get("/collections/slug/:slug", collectionsH.GetBySlug)
	api.Get("/collections/:id", collectionsH.Get)

	// ── Posts ─────────────────────────────────────────────────────────────
	api.Get("/posts", postsH.List)
	api.Get("/posts/slug/:slug", postsH.GetBySlug)
	api.Get("/posts/:id", postsH.Get)

	// ── Pages ─────────────────────────────────────────────────────────────
	api.Get("/pages", pagesH.List)
	api.Get("/pages/slug/:slug", pagesH.GetBySlug)
	api.Get("/pages/:id", pagesH.Get)

	// ── FAQ ───────────────────────────────────────────────────────────────
	api.Get("/faq", faqH.ListCategories)

	// ── Departments ───────────────────────────────────────────────────────
	api.Get("/departments", departmentsH.List)
	api.Get("/departments/slug/:slug", departmentsH.GetBySlug)
	api.Get("/departments/:id", departmentsH.Get)

	// ── Teachers ──────────────────────────────────────────────────────────
	api.Get("/teachers", teachersH.List)
	api.Get("/teachers/slug/:slug", teachersH.GetBySlug)
	api.Get("/teachers/:id", teachersH.Get)

	// ── Publication categories ────────────────────────────────────────────
	api.Get("/publication-categories", pubCategoriesH.List)
	api.Get("/publication-categories/:id", pubCategoriesH.Get)

	// ── Publications ──────────────────────────────────────────────────────
	api.Get("/publications", publicationsH.List)
	api.Get("/publications/slug/:slug", publicationsH.GetBySlug)
	api.Get("/publications/teacher/:slug", publicationsH.GetByTeacher) // before /:id
	api.Get("/publications/:id", publicationsH.Get)
	api.Get("/publications/:id/file", publicationsH.ServeFile)
	api.Get("/publications/:id/download", publicationsH.Download)

	// ── Search ────────────────────────────────────────────────────────────
	api.Get("/search", searchH.Get)

	// ── Static file serving (covers, logos, etc.) ─────────────────────────
	api.Get("/files/*", filesH.ServePublic)

	// ══════════════════════════════════════════════════════════════════════
	// Protected: any authenticated admin
	// ══════════════════════════════════════════════════════════════════════
	admin := api.Group("/admin", middleware.RequireAuth(tm))
	admin.Get("/me", authH.Me)

	// ── Stats ─────────────────────────────────────────────────────────────
	admin.Get("/stats", statsH.Overview)

	// ── Books (role-based section permissions checked inside handlers) ─────
	admin.Post("/books", booksH.Create)
	admin.Put("/books/:id", booksH.Update)
	admin.Delete("/books/:id", booksH.Delete)
	admin.Post("/books/:id/file", booksH.UploadFile)
	admin.Post("/books/:id/cover", booksH.UploadCover)

	// ══════════════════════════════════════════════════════════════════════
	// Superadmin only
	// ══════════════════════════════════════════════════════════════════════
	super := middleware.RequireSuperadmin()
	pubMgr := middleware.RequirePublicationManager()

	// ── Site settings ─────────────────────────────────────────────────────
	admin.Put("/site-settings", super, settingsH.Update)
	admin.Post("/site-settings/asset", super, settingsH.UploadAsset)
	admin.Post("/site-settings/font", super, settingsH.UploadFont)
	admin.Delete("/site-settings/font", super, settingsH.DeleteFont)

	// ── Sections ──────────────────────────────────────────────────────────
	admin.Post("/sections", super, sectionsH.Create)
	admin.Put("/sections/reorder", super, sectionsH.Reorder) // before /:id
	admin.Put("/sections/:id", super, sectionsH.Update)
	admin.Delete("/sections/:id", super, sectionsH.Delete)

	// ── Admins management ─────────────────────────────────────────────────
	admin.Get("/admins", super, adminsH.List)
	admin.Get("/admins/:id", super, adminsH.Get)
	admin.Post("/admins", super, adminsH.Create)
	admin.Put("/admins/:id", super, adminsH.Update)
	admin.Delete("/admins/:id", super, adminsH.Delete)
	admin.Put("/admins/:id/permissions", super, adminsH.SetPermissions)

	// ── Periodicals ───────────────────────────────────────────────────────
	admin.Post("/periodicals", super, periodicalsH.Create)
	admin.Put("/periodicals/:id", super, periodicalsH.Update)
	admin.Delete("/periodicals/:id", super, periodicalsH.Delete)
	admin.Post("/periodicals/:id/cover", super, periodicalsH.UploadCover)
	admin.Post("/periodicals/:id/issues", super, periodicalsH.CreateIssue)
	admin.Put("/periodicals/issues/:id", super, periodicalsH.UpdateIssue)
	admin.Delete("/periodicals/issues/:id", super, periodicalsH.DeleteIssue)
	admin.Post("/periodicals/issues/:id/file", super, periodicalsH.UploadIssueFile)
	admin.Post("/periodicals/issues/:id/cover", super, periodicalsH.UploadIssueCover)

	// ── Departments (publication_manager or superadmin) ───────────────────
	admin.Post("/departments", pubMgr, departmentsH.Create)
	admin.Put("/departments/:id", pubMgr, departmentsH.Update)
	admin.Delete("/departments/:id", pubMgr, departmentsH.Delete)

	// ── Teachers ──────────────────────────────────────────────────────────
	admin.Post("/teachers", pubMgr, teachersH.Create)
	admin.Put("/teachers/:id", pubMgr, teachersH.Update)
	admin.Delete("/teachers/:id", pubMgr, teachersH.Delete)
	admin.Post("/teachers/:id/photo", pubMgr, teachersH.UploadPhoto)

	// ── Publication categories ────────────────────────────────────────────
	admin.Post("/publication-categories", pubMgr, pubCategoriesH.Create)
	admin.Put("/publication-categories/:id", pubMgr, pubCategoriesH.Update)
	admin.Delete("/publication-categories/:id", pubMgr, pubCategoriesH.Delete)

	// ── Publications ──────────────────────────────────────────────────────
	admin.Post("/publications", pubMgr, publicationsH.Create)
	admin.Put("/publications/:id", pubMgr, publicationsH.Update)
	admin.Delete("/publications/:id", pubMgr, publicationsH.Delete)
	admin.Post("/publications/:id/file", pubMgr, publicationsH.UploadFile)

	// ── Publication stats ─────────────────────────────────────────────────
	admin.Get("/stats/publications", pubMgr, pubStatsH.Overview)

	// ── Collections ───────────────────────────────────────────────────────
	admin.Post("/collections", super, collectionsH.Create)
	admin.Put("/collections/:id", super, collectionsH.Update)
	admin.Delete("/collections/:id", super, collectionsH.Delete)
	admin.Post("/collections/:id/cover", super, collectionsH.UploadCover)
	admin.Put("/collections/:id/books/order", super, collectionsH.ReorderBooks) // before /:id
	admin.Post("/collections/:id/books", super, collectionsH.AddBook)
	admin.Delete("/collections/:id/books/:book_id", super, collectionsH.RemoveBook)

	// ── Posts / News ──────────────────────────────────────────────────────
	admin.Post("/posts", super, postsH.Create)
	admin.Put("/posts/:id", super, postsH.Update)
	admin.Delete("/posts/:id", super, postsH.Delete)
	admin.Post("/posts/:id/cover", super, postsH.UploadCover)

	// ── Static pages ──────────────────────────────────────────────────────
	admin.Post("/pages", super, pagesH.Create)
	admin.Put("/pages/:id", super, pagesH.Update)
	admin.Delete("/pages/:id", super, pagesH.Delete)

	// ── Inline content uploads (WYSIWYG images for posts/pages) ───────────
	admin.Post("/content/image", contentUploadsH.UploadImage)

	// ── FAQ ───────────────────────────────────────────────────────────────
	// Order routes registered before :id routes to avoid precedence conflicts
	admin.Post("/faq/categories", super, faqH.CreateCategory)
	admin.Put("/faq/categories/order", super, faqH.ReorderCategories)
	admin.Put("/faq/categories/:id", super, faqH.UpdateCategory)
	admin.Delete("/faq/categories/:id", super, faqH.DeleteCategory)
	admin.Post("/faq/items", super, faqH.CreateItem)
	admin.Put("/faq/items/order", super, faqH.ReorderItems)
	admin.Put("/faq/items/:id", super, faqH.UpdateItem)
	admin.Delete("/faq/items/:id", super, faqH.DeleteItem)

	// ── Home ──────────────────────────────────────────────────────────────
	admin.Get("/home/slides", super, homeH.ListSlides)
	admin.Post("/home/slides", super, homeH.CreateSlide)
	admin.Put("/home/slides/order", super, homeH.ReorderSlides) // before /:id
	admin.Put("/home/slides/:id", super, homeH.UpdateSlide)
	admin.Delete("/home/slides/:id", super, homeH.DeleteSlide)
	admin.Post("/home/slides/:id/image", super, homeH.UploadSlideImage)
	admin.Get("/home/sections", super, homeH.ListHomeSections)
	admin.Post("/home/sections", super, homeH.AddHomeSection)
	admin.Put("/home/sections/:id", super, homeH.UpdateHomeSection)
	admin.Delete("/home/sections/:id", super, homeH.RemoveHomeSection)
	admin.Get("/home/collections", super, homeH.ListHomeCollections)
	admin.Post("/home/collections", super, homeH.AddHomeCollection)
	admin.Put("/home/collections/:id", super, homeH.UpdateHomeCollection)
	admin.Delete("/home/collections/:id", super, homeH.RemoveHomeCollection)
	admin.Get("/home/promos", super, homeH.ListPromos)
	admin.Post("/home/promos", super, homeH.CreatePromo)
	admin.Put("/home/promos/order", super, homeH.ReorderPromos) // before /:id
	admin.Put("/home/promos/:id", super, homeH.UpdatePromo)
	admin.Delete("/home/promos/:id", super, homeH.DeletePromo)
	admin.Post("/home/promos/:id/image", super, homeH.UploadPromoImage)

	// ── Stats export ──────────────────────────────────────────────────────
	admin.Get("/stats/export.csv", super, statsH.ExportCSV)
}

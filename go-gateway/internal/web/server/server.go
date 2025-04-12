package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mpGustavo06/go-gateway-api/go-gateway/internal/service"
	"github.com/mpGustavo06/go-gateway-api/go-gateway/internal/web/handlers"
	"github.com/mpGustavo06/go-gateway-api/go-gateway/internal/web/middleware"
)

type Server struct {
	router *chi.Mux
	server *http.Server
	accountService *service.AccountService
	invoiceService *service.InvoiceService
	port string
}

func NewServer(accountService *service.AccountService, invoiceService *service.InvoiceService, port string) *Server {
	return &Server{
		router: chi.NewRouter(), 
		accountService: accountService, 
		invoiceService: invoiceService,
		port: port,
	}
}

func (s *Server) ConfigureRoutes() {
	accountHandler := handlers.NewAccountHandler(s.accountService)
	invoiceHandler := handlers.NewInvoiceHandler(s.invoiceService)
	authMiddleware := middleware.NewAuthMiddleware(s.accountService)

	s.router.Post("/accounts", accountHandler.Create)
	s.router.Get("/accounts", accountHandler.Get)

	s.router.Group(func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)

		s.router.Post("/invoices", invoiceHandler.Create)
		s.router.Get("/invoices", invoiceHandler.ListByAccount)
		s.router.Get("/invoices/{id}", invoiceHandler.GetByID)
	})

	
}

func (s *Server) Start() error {
	s.server = &http.Server {
		Addr: ":" + s.port, 
		Handler: s.router,
	}

	return s.server.ListenAndServe()
}
package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ic3network/mccs-alpha/internal/app/http/controller"
	"github.com/ic3network/mccs-alpha/internal/app/http/middleware"
)

func RegisterRoutes(r *mux.Router) {
	public := r.PathPrefix("/").Subrouter()
	public.Use(
		middleware.Recover(),
		middleware.NoCache(),
		middleware.Logging(),
		middleware.GetLoggedInUser(),
	)
	private := r.PathPrefix("/").Subrouter()
	private.Use(
		middleware.Recover(),
		middleware.NoCache(),
		middleware.Logging(),
		middleware.GetLoggedInUser(),
		middleware.RequireUser(),
	)
	adminPublic := r.PathPrefix("/admin").Subrouter()
	adminPublic.Use(
		middleware.Recover(),
		middleware.NoCache(),
		middleware.Logging(),
		middleware.GetLoggedInUser(),
	)
	adminPrivate := r.PathPrefix("/admin").Subrouter()
	adminPrivate.Use(
		middleware.Recover(),
		middleware.NoCache(),
		middleware.Logging(),
		middleware.GetLoggedInUser(),
		middleware.RequireAdmin(),
	)

	// Serving static files.
	fs := http.FileServer(http.Dir("web/static"))
	public.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	controller.ServiceDiscovery.RegisterRoutes(
		public,
		private,
		adminPublic,
		adminPrivate,
	)
	controller.DashBoardHandler.RegisterRoutes(
		public,
		private,
		adminPublic,
		adminPrivate,
	)
	controller.BusinessHandler.RegisterRoutes(
		public,
		private,
		adminPublic,
		adminPrivate,
	)
	controller.UserHandler.RegisterRoutes(
		public,
		private,
		adminPublic,
		adminPrivate,
	)
	controller.TransactionHandler.RegisterRoutes(
		public,
		private,
		adminPublic,
		adminPrivate,
	)
	controller.HistoryHandler.RegisterRoutes(
		public,
		private,
		adminPublic,
		adminPrivate,
	)
	controller.TradingHandler.RegisterRoutes(
		public,
		private,
		adminPublic,
		adminPrivate,
	)

	controller.AdminBusinessHandler.RegisterRoutes(
		public,
		private,
		adminPublic,
		adminPrivate,
	)
	controller.AdminUserHandler.RegisterRoutes(
		public,
		private,
		adminPublic,
		adminPrivate,
	)
	controller.AdminHistoryHandler.RegisterRoutes(
		public,
		private,
		adminPublic,
		adminPrivate,
	)
	controller.AdminTransactionHandler.RegisterRoutes(
		public,
		private,
		adminPublic,
		adminPrivate,
	)
	controller.AdminTagHandler.RegisterRoutes(
		public,
		private,
		adminPublic,
		adminPrivate,
	)
	controller.LogHandler.RegisterRoutes(
		public,
		private,
		adminPublic,
		adminPrivate,
	)

	controller.AccountHandler.RegisterRoutes(
		public,
		private,
		adminPublic,
		adminPrivate,
	)
	controller.TagHandler.RegisterRoutes(
		public,
		private,
		adminPublic,
		adminPrivate,
	)
}

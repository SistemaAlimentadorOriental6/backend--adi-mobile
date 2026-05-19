package main

import (
	"log"
	"net/http"
	"os"
	"turnovia-backend/internal/infrastructure/database"
	"turnovia-backend/internal/infrastructure/handler"
	"turnovia-backend/internal/infrastructure/repository"
	"turnovia-backend/internal/usecase"
)

func middlewareCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {

	dbMySQL := database.InitDB()
	defer dbMySQL.Close()

	dbSQLServer := database.InitSQLServer()
	defer dbSQLServer.Close()
	userMySQLRepo := repository.NewMysqlUserRepository(dbMySQL)
	userSQLServerRepo := repository.NewSqlServerUserRepository(dbSQLServer)
	reporteRepo := repository.NewMysqlReporteRepository(dbMySQL)
	trackingRepo := repository.NewMysqlTrackingRepository(dbMySQL)

	authUC := usecase.NewAuthUseCase(userMySQLRepo, userSQLServerRepo)
	reporteUC := usecase.NewReporteUseCase(reporteRepo)
	trackingUC := usecase.NewTrackingUseCase(trackingRepo)

	authHandler := handler.NewAuthHandler(authUC)
	reporteHandler := handler.NewReporteHandler(reporteUC)
	trackingHandler := handler.NewTrackingHandler(trackingUC)
	proxyRepo := repository.NewProxyEmpleadoRepository()
	proxyHandler := handler.NewProxyHandler(proxyRepo, userSQLServerRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/validate-login", authHandler.Login)
	mux.HandleFunc("/api/registrar-actividad", reporteHandler.Registrar)
	mux.HandleFunc("/api/listar-actividad", reporteHandler.Listar)
	mux.HandleFunc("/api/registro-hoy", reporteHandler.ObtenerRegistroHoy)
	mux.HandleFunc("/api/tracking", trackingHandler.GuardarUbicacion)
	mux.HandleFunc("/api/tracking/batch", trackingHandler.GuardarUbicacionBatch)
	mux.HandleFunc("/api/proxy/empleados", proxyHandler.GetEmpleados)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3036"
	}

	addr := ":" + port
	log.Printf("🚀 Servidor Turnovia corriendo en el puerto %s", addr)

	if err := http.ListenAndServe(addr, middlewareCORS(mux)); err != nil {
		log.Fatal("Error iniciando el servidor: ", err)
	}
}

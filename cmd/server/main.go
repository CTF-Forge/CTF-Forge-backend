package server

import docs "github.com/CTF-Forge/CTF-Forge-backend/docs"

func init() {
	docs.SwaggerInfo.Title = "CTFForge API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
}

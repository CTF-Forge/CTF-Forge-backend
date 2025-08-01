package server

import docs "github.com/Saku0512/CTFLab/ctflab/docs"

func init() {
	docs.SwaggerInfo.Title = "CTFLab API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
}

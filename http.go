package sweb

import (
	"sweb/pages"
	"log"
	"net/http"
	"path/filepath"
)

func ServeStatic(config ServerConfig, dataDir string, appName string) {
	staticDir := filepath.Join(dataDir, "web", appName, "static")
	http.Handle( "/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
}

func ServeHandler(config ServerConfig, handler pages.Handler) {
	log.Println("[HTTP]", "registered handler", "/")
	http.Handle("/", handler)
}

func Listen(config ServerConfig) {
	if config.Ssl {
		log.Println("[HTTP]", "listening to", config.Address, "https")
		err := http.ListenAndServeTLS(config.Address, config.CertFile, config.KeyFile, nil)

		if err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Println("[HTTP]", "listening to", config.Address, "http")
		err := http.ListenAndServe(config.Address, nil)

		if err != nil {
			log.Fatalln(err)
		}
	}
}

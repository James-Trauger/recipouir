package main

import (
	"flag"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var (
	addr = flag.String("listen", "127.0.0.1:8443", "listen address")
	cert = flag.String("cert", "", "certificate")
	pkey = flag.String("key", "", "private key")
)

func main() {
	flag.Parse()
	// load the config file
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("couldn't load .env -> ", err)
	}
	certPath := os.Getenv("LOCAL_CERT")
	if certPath != "" {
		certPath = "../" + certPath
		certFile := os.Getenv("CERT_NAME")
		if *cert == "" && certFile != "" {
			*cert = certPath + certFile
		} else {
			log.Fatal("no certificate provided")
		}
		pkeyFile := os.Getenv("KEY_NAME")
		if *pkey == "" && pkeyFile != "" {
			*pkey = certPath + pkeyFile
		}
	}

	//; charset=utf-8
	mime.AddExtensionType(".js", "text/javascript; charset=utf-8")
	mime.AddExtensionType(".css", "text/css; charset=utf-8")
	mux := http.NewServeMux()

	routes := []string{
		"/", // home page
		"/about",
	}

	// serve assets built with react
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	// serve a single static page
	for _, r := range routes {
		mux.Handle(r, rootHandler())
	}

	fmt.Println(http.ListenAndServeTLS(*addr, *cert, *pkey, mux))
}

// serves index.tmpl
func rootHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		http.ServeFile(w, r, "./views/index.html")
	})
}

func serverErrorHandler(err error, from string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error from ", from, " -> ", err.Error())
		// write the error message
		_, writeErr := w.Write([]byte("internal server error"))
		if writeErr != nil {
			log.Println("couldn't write error message -> ", writeErr)
		}
	})
}

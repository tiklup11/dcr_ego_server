package service

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func RunEgoServer() {
	fmt.Print("ego-server started on port 8080 \n")

	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/run", RunHandler)
	http.HandleFunc("/download", HandleDownload)
	
	http.ListenAndServe(":8080", nil)
	// Use the default ServeMux.
	// enableHttps()
}

func enableHttps() {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatal(err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	server := &http.Server{
		Addr:      ":8443",
		Handler:   nil,
		TLSConfig: config,
	}

	log.Printf("Server listening on https://localhost%s\n", server.Addr)
	log.Fatal(server.ListenAndServeTLS("", ""))
}

package main

import(
	"net/http"
	"os"
	"time"
)

func startServer(){
	http.HandleFunc("/stt",stt)
	var server=&http.server{
		port: ":7891",
		ReadTimeout: 60 * time.second,
		WriteTimeout: 60 * time.second,
		MaxHeaderBytes: 1 << 20,
	}

	err := server.ListenAndServe()
	if nil != err{
		Error("Fail to listen the server :", err)
		os.Exit(-1)
	}
}
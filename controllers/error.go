package controllers

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func NotFound(w http.ResponseWriter, r *http.Request, p httprouter.Params, err error) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 page not found"))
	log.Println(err)

}

func InternalServerError(w http.ResponseWriter, r *http.Request, p httprouter.Params, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 internal server error"))
	log.Println(err)
}

func BadRequest(w http.ResponseWriter, r *http.Request, p httprouter.Params, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("400 bad request"))
	log.Println(err)
}

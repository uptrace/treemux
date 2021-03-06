package main

import (
	"errors"
	"log"
	"math/rand"
	"net/http"

	"github.com/uptrace/treemux"
	"github.com/uptrace/treemux/extra/reqlog"
)

var (
	err1 = errors.New("error1")
	err2 = errors.New("error2")
)

func main() {
	router := treemux.New(
		treemux.WithMiddleware(reqlog.NewMiddleware()),
		treemux.WithMiddleware(errorHandler),
	)

	router.GET("/", indexHandler)

	log.Println("listening on http://localhost:8888")
	log.Println(http.ListenAndServe(":8888", router))
}

func indexHandler(w http.ResponseWriter, req treemux.Request) error {
	if rand.Float64() > 0.5 {
		return err1
	}
	return err2
}

func errorHandler(next treemux.HandlerFunc) treemux.HandlerFunc {
	return func(w http.ResponseWriter, req treemux.Request) error {
		err := next(w, req)

		switch err {
		case nil:
			// ok
		case err1:
			w.WriteHeader(http.StatusBadRequest)
			_ = treemux.JSON(w, treemux.H{
				"message": "bad request",
				"hint":    "reload to see how error message is changed",
			})
		default:
			w.WriteHeader(http.StatusInternalServerError)
			_ = treemux.JSON(w, treemux.H{
				"message": err.Error(),
				"hint":    "reload to see how error message is changed",
			})
		}

		return err
	}
}

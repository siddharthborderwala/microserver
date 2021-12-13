package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type GoodBye struct {
	l *log.Logger
}

func NewGoodBye(l *log.Logger) *GoodBye {
	return &GoodBye{l}
}

// a handlerFunc need to implement the ServeHTTP function
func (h *GoodBye) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Oops", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Goodbye %s", d)
}

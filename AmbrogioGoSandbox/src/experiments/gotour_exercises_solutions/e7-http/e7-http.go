package main

import (
	"fmt"
	"net/http"
)

/* Types */

type String string

type Responder struct {
	Greeting string
	Punct    string
	Who      string
}

/* Handlers */

func (self String) ServeHTTP( w  http.ResponseWriter,
                              r *http.Request         ) {
	fmt.Fprint(w, self)
}

func (self Responder) ServeHTTP( w  http.ResponseWriter,
                                 r *http.Request         ) {
	fmt.Fprint(w,self.Greeting+" "+self.Who)
}

/* main */

func main() {
	http.Handle("/string", String("I'm a frayed knot."))
	http.Handle("/answer", &Responder{"Hello", ":", "Gophers!"})
	http.ListenAndServe("localhost:4000", nil)
}

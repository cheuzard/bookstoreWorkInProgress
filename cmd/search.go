package bookstore

import "net/http"

type searchArguments struct {
	advanced  bool
	booksID   int
	isbn      int
	title     string
	author    string
	publisher string
	category  int
}

func (s searchArguments) lookFor(w http.ResponseWriter) {

}

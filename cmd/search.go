package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	client       = "client"
	admin        = "admin"
	baseQuery    = "select isbn, title, author, publisher, publication_date, description, price, rental_price_per_day"
	StaffQuery   = ", quantity, created_at, updated_at "
	queryEndISBN = " from books where isbn = ?"
	queryEndTerm = " from books WHERE title LIKE ? OR description LIKE ? ORDER BY CASE WHEN title LIKE ? THEN 1 WHEN description LIKE ? THEN 2 ELSE 3 END, title LIMIT 50"
	queryEndAdv  = " from books where title = ? and author = ? and publisher = ? and category = ?"
)

// enum to handle the type of
type queryType int

const (
	isbnQuery queryType = iota
	termQuery
	advancedQuery
	e
)

// struct to manage the search arguments sent by the client
type searchArguments struct {
	advanced  bool
	source    string
	isbn      int
	title     string
	author    string
	publisher string
	category  int
}

type BookStaff struct {
	ISBN              string    `json:"isbn" db:"isbn"`
	Title             string    `json:"title" db:"title"`
	Author            string    `json:"author" db:"author"`
	Publisher         string    `json:"publisher" db:"publisher"`
	PublicationDate   time.Time `json:"publication_date" db:"publication_date"`
	Description       string    `json:"description" db:"description"`
	Price             float64   `json:"price" db:"price"`
	RentalPricePerDay float64   `json:"rental_price_per_day" db:"rental_price_per_day"`
	Quantity          int       `json:"quantity" db:"quantity"`
	//Category          string    `json:"category" db:"category"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (b BookStaff) GetTitle() string { return b.Title }

type BookClient struct {
	ISBN              string    `json:"isbn" db:"isbn"`
	Title             string    `json:"title" db:"title"`
	Author            string    `json:"author" db:"author"`
	Publisher         string    `json:"publisher" db:"publisher"`
	PublicationDate   time.Time `json:"publication_date" db:"publication_date"`
	Description       string    `json:"description" db:"description"`
	Price             float64   `json:"price" db:"price"`
	RentalPricePerDay float64   `json:"rental_price_per_day" db:"rental_price_per_day"`
	//Category          string    `json:"category" db:"category"`
}

func (b BookClient) GetTitle() string { return b.Title }

type Book interface {
	GetTitle() string
}

func (app App) OpenConnectionHandler() http.Handler {
	log.Printf("Open connection Handler started \n")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		result, err := searchParser(r).dbCall(app.db)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("results received %v\n", result)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			return
		}
		fmt.Printf("finished encoding results\n")
		println()
		println()
		println()
	})
}

func searchParser(r *http.Request) *searchArguments {
	log.Printf("search parser started \n")

	params := r.URL.Query()
	return &searchArguments{
		advanced: func() bool {
			if adv, err := strconv.ParseBool(params.Get("advanced")); err == nil {
				return adv
			}
			return false
		}(),
		source: params.Get("source"),
		isbn: func() int {
			if n, err := strconv.Atoi(params.Get("isbn")); err == nil {
				return n
			} else {
				return 0
			}
		}(),

		title:     params.Get("title"),
		author:    params.Get("author"),
		publisher: params.Get("publisher"),
		category: func() int {
			if n, err := strconv.Atoi(params.Get("category")); err == nil {
				return n
			} else {
				return 0
			}
		}(),
	}
}

func (s searchArguments) queryMaker() (string, queryType, error) {
	query := baseQuery
	QType := e

	if s.source == admin {
		query = query + StaffQuery
	} else if s.source == client {
		//nothing special
	} else {
		return "", QType, errors.New("query maker e: invalid source")
	}

	if s.isbn != 0 {
		query = query + queryEndISBN
		QType = isbnQuery
	} else if s.advanced {
		query = query + queryEndAdv
		QType = advancedQuery
	} else {
		query = query + queryEndTerm
		QType = termQuery
	}
	return query, QType, nil
}

func (s searchArguments) dbCall(db *sql.DB) ([]Book, error) {
	log.Printf("argument parssing finished, started dbCall\nargs:\n%v\n", s)
	var rows *sql.Rows

	query, QType, err := s.queryMaker()
	if err != nil {
		return nil, err
	}
	log.Printf("query string generated: %s\n", query)
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	s.title = fmt.Sprintf("%%%s%%", s.title)

	switch QType {
	case isbnQuery:
		rows, err = stmt.Query(s.isbn)
	case termQuery:
		rows, err = stmt.Query(s.title, s.title, s.title, s.title)
	case advancedQuery:
		rows, err = stmt.Query(s.title, s.author, s.publisher, s.category)
	default:
		return nil, fmt.Errorf("unhandled query type")
	}
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var books []Book
	if s.source == "admin" {
		var book BookStaff
		for rows.Next() {
			err := rows.Scan(
				&book.ISBN,
				&book.Title,
				&book.Author,
				&book.Publisher,
				&book.PublicationDate,
				&book.Description,
				&book.Price,
				&book.RentalPricePerDay,
				&book.Quantity,
				&book.CreatedAt,
				&book.UpdatedAt,
			)
			if err != nil {
				return nil, err
			}
			log.Printf("row received: %v\n", book)
			books = append(books, book)
		}
	} else {
		var book BookClient
		for rows.Next() {
			err := rows.Scan(
				&book.ISBN,
				&book.Title,
				&book.Author,
				&book.Publisher,
				&book.PublicationDate,
				&book.Description,
				&book.Price,
				&book.RentalPricePerDay,
			)
			if err != nil {
				return nil, err
			}
			log.Printf("row received: %v\n", book)
			books = append(books, book)
		}
	}
	return books, nil
}

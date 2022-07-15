// Copyright (c) 2022, Geert JM Vanderkelen

package kolekto

import (
	"testing"

	"github.com/geertjanvdk/xkit/xt"
	"github.com/golistic/kolekto/kolektor"
)

type Book struct {
	kolektor.Model
	ISBN13  string   `json:"isbn13"`
	Title   string   `json:"title"`
	Authors []string `json:"authors,omitempty"`
}

var _ kolektor.Modeler = &Book{}

func (b Book) CollectionName() string {
	return "books"
}

func (b Book) Indexes(kind kolektor.StoreKind) []kolektor.Index {
	m := map[kolektor.StoreKind][]kolektor.Index{
		kolektor.MySQL: {
			{
				Name:       "uq_books_isbn13",
				Unique:     true,
				Expression: "((CAST(data->>'$.isbn13' AS CHAR(20))))",
			},
			{
				Name:       "ix_books_title",
				Expression: "((CAST(data->>'$.title' AS CHAR(200))))",
			},
		},
		kolektor.PgSQL: {
			{
				Name:       "uq_books_isbn13",
				Unique:     true,
				Expression: "((data->>'isbn13'))",
			},
			{
				Name:       "ix_books_title",
				Expression: "((data->>'title'))",
			},
		},
	}

	if res, have := m[kind]; have {
		return res
	}

	return nil
}

func testCollection_Store(t *testing.T, session *Session) {
	t.Run("store object and retrieve it", func(t *testing.T) {
		book := &Book{
			ISBN13: "978-0135800911",
			Title:  "Presentation Zen",
		}
		books, err := session.Collection(book)

		t.Run(testDBPrefix(session)+" store book using", func(t *testing.T) {
			xt.OK(t, err)
			xt.OK(t, books.Store(book))

			t.Run("retrieve the book & update", func(t *testing.T) {
				xt.OK(t, err)
				b := &Book{}
				xt.OK(t, books.Get(b, book.Meta.ID))
				xt.Eq(t, book.Meta.UID, b.Meta.UID)
				xt.Eq(t, book.ISBN13, b.ISBN13)
				xt.Eq(t, nil, book.Meta.Updated)

				authors := []string{"Garr Reynolds"}
				b.Authors = authors
				xt.OK(t, books.Store(b))
				xt.Eq(t, book.Meta.ID, b.Meta.ID)

				t.Run("verify update", func(t *testing.T) {
					xt.OK(t, err)
					b := &Book{}
					xt.OK(t, books.Get(b, book.Meta.ID))
					xt.OK(t, err)
					xt.Eq(t, book.ISBN13, b.ISBN13)
					xt.Eq(t, authors, b.Authors)
				})
			})
		})
	})
}

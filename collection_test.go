// Copyright (c) 2022, Geert JM Vanderkelen

package kolekto

import (
	"testing"

	"github.com/geertjanvdk/xkit/xt"
	"github.com/golistic/kolekto/kolektor"
	"github.com/golistic/kolekto/stores"
)

func TestCollection_ByFields(t *testing.T) {
	for storeKind, storeFn := range stores.Registered() {
		session, err := newSession(testAllDSN[storeKind], storeFn)
		xt.OK(t, err)

		t.Run(storeKind.String(), func(t *testing.T) {
			booksData := []*Book{
				{
					Model: kolektor.Model{
						Meta: &kolektor.Meta{UID: "book1.foo"},
					},
					ISBN13:    "978-3-1194-1744-0",
					Title:     "Book1",
					Publisher: "Foo",
				},
				{
					ISBN13:    "978-5-5861-6011-9",
					Title:     "Book2",
					Publisher: "Foo",
				},
				{
					ISBN13:    "978-9-6557-3995-4",
					Title:     "Book3",
					Publisher: "Bar",
				},
			}
			books, err := session.Collection(&Book{})

			for _, b := range booksData {
				xt.OK(t, err)
				xt.OK(t, books.Store(b))
			}

			t.Run("get object using one field", func(t *testing.T) {
				book := &Book{}
				expISBN13 := "978-3-1194-1744-0"
				xt.OK(t, books.GetByFields(book, kolektor.FieldMap{"isbn13": expISBN13}))
				xt.Eq(t, expISBN13, book.ISBN13)
			})

			t.Run("get object using reserved field", func(t *testing.T) {
				book := &Book{}
				expISNB13 := "978-3-1194-1744-0"
				xt.OK(t, books.GetByFields(book, kolektor.FieldMap{"uid": "book1.foo"}))
				xt.Eq(t, expISNB13, book.ISBN13)
			})

			t.Run("get object using reserved field and data field", func(t *testing.T) {
				book := &Book{}
				expISNB13 := "978-3-1194-1744-0"
				xt.OK(t, books.GetByFields(book, kolektor.FieldMap{
					"uid": "book1.foo", "isbn13": expISNB13}))
				xt.Eq(t, expISNB13, book.ISBN13)
			})
		})

	}
}

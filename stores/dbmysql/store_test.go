// Copyright (c) 2022, Geert JM Vanderkelen

//go:build !nomysql

package dbmysql

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

	fuIndex          func() map[kolektor.StoreKind][]kolektor.Index
	fuCollectionName func() string
}

var _ kolektor.Modeler = &Book{}

func (b Book) CollectionName() string {
	return b.fuCollectionName()
}

func (b Book) Indexes(kind kolektor.StoreKind) []kolektor.Index {
	m := b.fuIndex()

	if res, have := m[kind]; have {
		return res
	}

	return nil
}

func TestStore_InitCollection(t *testing.T) {
	s, err := New(testDSN)
	xt.OK(t, err)
	store := s.(*Store)

	expIndex1 := []string{"uq_books_isbn13", "((CAST(data->>'$.isbn13' AS CHAR(20))))"}
	expIndex2 := []string{"ix_books_title", "((CAST(data->>'$.title' AS CHAR(200))))"}

	t.Run("add indexes", func(t *testing.T) {
		fuCollectionName := func() string {
			return "books_239d9k3d"
		}

		book := &Book{}
		book.fuCollectionName = fuCollectionName

		book.fuIndex = func() map[kolektor.StoreKind][]kolektor.Index {
			return map[kolektor.StoreKind][]kolektor.Index{
				kolektor.MySQL: {
					{
						Name:       expIndex1[0],
						Unique:     true,
						Expression: expIndex1[1],
					},
					{
						Name:       expIndex2[0],
						Expression: expIndex2[1],
					},
				},
			}
		}

		xt.OK(t, store.InitCollection(book))
		indexes, err := getIndexes(store.mustSQLConn(), book.CollectionName())
		xt.OK(t, err)
		xt.Eq(t, 2, len(indexes))
		exprSum, _ := indexes[expIndex1[0]]
		xt.Eq(t, md5sum(expIndex1[1]), exprSum)
		exprSum, _ = indexes[expIndex2[0]]
		xt.Eq(t, md5sum(expIndex2[1]), exprSum)

		t.Run("remove index", func(t *testing.T) {
			book := &Book{}
			book.fuCollectionName = fuCollectionName

			book.fuIndex = func() map[kolektor.StoreKind][]kolektor.Index {
				return map[kolektor.StoreKind][]kolektor.Index{
					kolektor.MySQL: {
						{
							Name:       expIndex1[0],
							Unique:     true,
							Expression: expIndex1[1],
						},
					},
				}
			}

			xt.OK(t, store.InitCollection(book))
			indexes, err := getIndexes(store.mustSQLConn(), book.CollectionName())
			xt.OK(t, err)
			xt.Eq(t, 1, len(indexes))
			exprSum, _ := indexes[expIndex1[0]]
			xt.Eq(t, md5sum(expIndex1[1]), exprSum)
		})

		t.Run("change index", func(t *testing.T) {
			book := &Book{}
			book.fuCollectionName = fuCollectionName

			// change length of char
			expIndex1 := []string{"uq_books_isbn13", "((CAST(data->>'$.isbn13' AS CHAR(40))))"}
			book.fuIndex = func() map[kolektor.StoreKind][]kolektor.Index {
				return map[kolektor.StoreKind][]kolektor.Index{
					kolektor.MySQL: {
						{
							Name:       expIndex1[0],
							Unique:     true,
							Expression: expIndex1[1],
						},
					},
				}
			}

			xt.OK(t, store.InitCollection(book))
			indexes, err := getIndexes(store.mustSQLConn(), book.CollectionName())
			xt.OK(t, err)
			xt.Eq(t, 1, len(indexes))
			exprSum, _ := indexes[expIndex1[0]]
			xt.Eq(t, md5sum(expIndex1[1]), exprSum)
		})
	})
}

// Copyright (c) 2022, Geert JM Vanderkelen

package kolekto_test

import (
	"fmt"
	"os"
	"text/template"

	"github.com/golistic/kolekto"
	"github.com/golistic/kolekto/kolektor"
)

type Band struct {
	kolektor.Model
	Name    string   `json:"name"`
	Members []string `json:"members"`
	Active  bool     `json:"active"`
}

func (b Band) CollectionName() string {
	return "bands"
}

func exampleStoreRetrieveBand(kind kolektor.StoreKind, dsn string) {
	session, err := kolekto.NewSession(kind, dsn)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	if err := session.RemoveCollection(&Band{}); err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	bands, err := session.Collection(&Band{})
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	newBand := &Band{
		Name:    "A Tribe Called Quest", // first band that crossed my mind!
		Members: []string{"Q-Tip", "Phife Dwag", "Ali Shaheed Muhammad", "Jarobi White"},
	}
	// make test deterministic (usually generated)
	newBand.Meta = &kolektor.Meta{
		UID: "f5dea144-caac-4735-a521-34a82b12f20b",
	}

	if err := bands.Store(newBand); err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	band := &Band{}
	if err := bands.Get(band, newBand.Meta.UID); err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	var tmplData = struct {
		Band *Band
	}{
		Band: band,
	}

	tmpl, err := template.New("").Parse(`
UID    : f5dea144-caac-4735-a521-34a82b12f20b
Band   : {{.Band.Name}}
Members:
{{- range $member := .Band.Members }}
 - {{ $member }}
{{- end}}
`)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	if err := tmpl.Execute(os.Stdout, tmplData); err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
}

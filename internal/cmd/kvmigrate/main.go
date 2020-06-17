package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/influxdata/influxdb/v2/kv/migration/all"
)

var usageMsg = "Usage: kvmigrate create <migration name / description>"

func usage() {
	fmt.Println(usageMsg)
	os.Exit(1)
}

const newMigrationFmt = `package all

var %s = &Migration{}
`

func main() {
	if len(os.Args) < 3 {
		usage()
	}

	if os.Args[1] != "create" {
		fmt.Printf("unrecognized command %q\n", os.Args[1])
		usage()
	}

	migrationName := strings.Join(os.Args[2:], " ")

	camelName := strings.Replace(strings.Title(migrationName), " ", "", -1)

	newMigrationNumber := len(all.Migrations) + 1

	newMigrationVariable := fmt.Sprintf("Migration%04d_%s", newMigrationNumber, camelName)

	newMigrationFile := fmt.Sprintf("./kv/migration/all/%04d_%s.go", newMigrationNumber, strings.Join(os.Args[2:], "_"))

	fmt.Println("Creating new migration:", newMigrationFile)

	if err := ioutil.WriteFile(newMigrationFile, []byte(fmt.Sprintf(newMigrationFmt, newMigrationVariable)), 0644); err != nil {
		panic(err)
	}

	fmt.Println("Inserting migration into ./kv/migration/all/all.go")

	tmplData, err := ioutil.ReadFile("./kv/migration/all/all.go")
	if err != nil {
		panic(err)
	}

	type Context struct {
		Name     string
		Variable string
	}

	tmpl := template.Must(
		template.
			New("migrations").
			Funcs(template.FuncMap{"do_not_edit": func(c Context) string {
				return fmt.Sprintf("%s\n%s,\n// {{ do_not_edit . }}", c.Name, c.Variable)
			}}).
			Parse(string(tmplData)),
	)

	buf := new(bytes.Buffer)

	if err := tmpl.Execute(buf, Context{
		Name:     migrationName,
		Variable: newMigrationVariable,
	}); err != nil {
		panic(err)
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile("./kv/migration/all/all.go", src, 0644); err != nil {
		panic(err)
	}
}

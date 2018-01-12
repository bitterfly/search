package documents

import (
	"fmt"
	"strings"
)

type Document struct {
	Title   string
	Classes []string
	Body    string
	Date    string
}

func (d *Document) String() string {
	return fmt.Sprintf(
		"[%s], classes [%s], date %s:\n%s\n",
		d.Title,
		strings.Join(d.Classes, ", "),
		d.Date,
		d.Body,
	)
}

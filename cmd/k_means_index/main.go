package main

import (
	"log"
	"os"

	"github.com/bitterfly/search/indices"
	"github.com/bitterfly/search/kmeans"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "k_means_index"
	app.Usage = "Creates index for kMeans"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "input, i",
			Usage: "File with index",
			Value: "/tmp/index.gob.gz",
		},

		cli.IntFlag{
			Name:  "k",
			Usage: "The parameter of k means",
			Value: 3,
		},
		cli.StringFlag{
			Name:  "s",
			Usage: "File to serialise to",
			Value: "/tmp/kmeans.gob.gz",
		},
	}

	app.Action = mainCommand

	app.Run(os.Args)
}

func mainCommand(c *cli.Context) {
	ti := indices.NewTotalIndex()
	err := ti.DeserialiseFromFile(c.String("input"))
	if err != nil {
		log.Fatal(err)
	}

	kmeans.ProcessArguments(ti, c.Int("k"))

	err = ti.SerialiseToFile(c.String("s"))
	if err != nil {
		log.Printf("Could not write to file %s", c.String("s"))
	}
}

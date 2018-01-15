package main

import (
	"log"
	"os"

	"github.com/DexterLB/search/indices"
	"github.com/DexterLB/search/kmeans"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "kmeans"
	app.Usage = "k-means"
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

	kmeans.KMeans(ti, c.Int("k"))
}

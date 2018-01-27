package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/DexterLB/search/featureselection"
	"github.com/DexterLB/search/indices"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "knn"
	app.Usage = "Perform kNN"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "input, i",
			Usage: "File with index",
			Value: "/tmp/index.gob.gz",
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

	features := featureselection.ChiSquared(ti, 10, runtime.NumCPU())

	for _, termID := range features {
		fmt.Printf("%d\n", termID)
	}
}

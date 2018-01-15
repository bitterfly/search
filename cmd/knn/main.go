package main

import (
	"fmt"
	"log"
	"os"

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

	info := featureselection.ComputeClassInfo(ti)
	for i, count := range info.DocumentsWhichHaveClass {
		fmt.Printf("%s: %d\n", ti.ClassNames.GetInverse(int32(i)), count)
	}
}

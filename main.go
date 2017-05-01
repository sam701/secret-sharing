package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/hashicorp/vault/shamir"
	"github.com/urfave/cli"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	app := cli.NewApp()
	app.Version = "0.1.0"
	app.Commands = []cli.Command{
		{
			Name:      "split",
			Usage:     "split file into parts",
			ArgsUsage: "<file to split>",
			Action:    split,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "parts, p",
					Usage: "number of parts to split into",
				},
				cli.IntFlag{
					Name:  "threshold, t",
					Usage: "threshold",
				},
				cli.StringFlag{
					Name:  "output-dir, o",
					Usage: "output directory that will contain <file-name>.partn files",
					Value: ".",
				},
			},
		},
		{
			Name:      "combine",
			Usage:     "combine parts into a file",
			ArgsUsage: "<part0> [<part1> ...]",
			Action:    combine,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output, o",
					Usage: "output file name",
				},
			},
		},
	}
	app.Run(os.Args)

}

var encoding = base64.RawURLEncoding

func split(ctx *cli.Context) error {
	no := ctx.Int("parts")
	if no == 0 {
		return errors.New("parts flag is missing")
	}
	threshold := ctx.Int("threshold")
	if threshold == 0 {
		return errors.New("threshold flag is missing")
	}

	inputFile := ctx.Args().First()
	if inputFile == "" {
		return errors.New("No input file")
	}
	outputDir := ctx.String("output-dir")

	content, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	out, err := shamir.Split(content, no, threshold)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	for ix, part := range out {
		fn := fmt.Sprintf("%s.part%d", path.Base(inputFile), ix)
		fn = path.Join(outputDir, fn)
		err = ioutil.WriteFile(fn, part, 0600)
		if err != nil {
			log.Fatalln("ERROR", err)
		}
	}

	return nil
}

func combine(ctx *cli.Context) error {
	if !ctx.Args().Present() {
		return errors.New("No parts present")
	}

	input := [][]byte{}
	for _, part := range ctx.Args() {
		bb, err := ioutil.ReadFile(part)
		if err != nil {
			log.Fatalln("ERROR", err)
		}
		input = append(input, bb)
	}

	of := ctx.String("output")
	if of == "" {
		return errors.New("No output present")
	}

	out, err := shamir.Combine(input)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(of, out, 0600)
	if err != nil {
		return err
	}

	fmt.Println("Result was written into", of)

	return nil
}

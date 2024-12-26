package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"path"
	"strings"
	"text/template"
)

type CexPlatform struct {
	Package string
}

const input = "template"

const output = "platforms"

func main() {
	app := &cli.App{
		Name:  "cex generator",
		Usage: "A CLI tool to quickly generate exchange components from templates.",
		Commands: cli.Commands{
			{
				Name:    "gen",
				Aliases: []string{"i"},
				Usage:   "Initialize a new exchange project.",
				Action: func(ctx *cli.Context) error {
					projectName := ctx.Args().First()
					if projectName == "" {
						return fmt.Errorf("project name is required")
					}
					root, err := os.Getwd()
					if err != nil {
						return err
					}

					tmplPath := path.Join(root, input, "cex")
					destPath := path.Join(root, output)
					data := CexPlatform{
						Package: projectName,
					}
					if err = os.MkdirAll(path.Join(destPath, projectName), os.ModePerm); err != nil {
						return err
					}

					templates, err := template.ParseGlob(path.Join(tmplPath, "*.tmpl"))
					if err != nil {
						return err
					}
					for _, tmpl := range templates.Templates() {
						pathName := tmpl.Name()
						if strings.HasSuffix(pathName, "cex.go.tmpl") {
							pathName = strings.ReplaceAll(pathName, "cex.go.tmpl", projectName+".go")
						} else if strings.HasSuffix(pathName, "cex_test.go.tmpl") {
							pathName = strings.ReplaceAll(pathName, "cex_test.go.tmpl", projectName+"_test.go")
						} else {
							pathName = strings.ReplaceAll(pathName, ".tmpl", "")
						}

						fs, err := os.Create(path.Join(destPath, projectName, pathName))
						if err != nil {
							continue
						}

						err = tmpl.Execute(fs, data)
						if err != nil {
							continue
						}
						err = fs.Close()
						if err != nil {
							continue
						}
					}
					return nil
				},
			},
			{
				Name:    "remove",
				Aliases: []string{"rm"},
				Action: func(ctx *cli.Context) error {
					projectName := ctx.Args().First()
					root, err := os.Getwd()
					if err != nil {
						return err
					}
					destPath := path.Join(root, output, projectName)
					err = os.RemoveAll(destPath)
					if err != nil {
						return err
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

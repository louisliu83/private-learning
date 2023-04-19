package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
	"pa.cn/fedlearn/psi/cmd/tools/ctl/useradd"
)

func main() {
	app := cli.NewApp()
	app.Name = "PSIctl"
	app.Description = "psi related tools"
	app.Authors = []cli.Author{
		{
			Name:  "Zhang Dongqi",
			Email: "zhangdongqi217@pingan.com.cn",
		},
	}
	app.Flags = []cli.Flag{}
	app.Commands = []cli.Command{
		useradd.Command(),
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

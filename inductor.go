package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := createApp()
	app.Run(os.Args)
	//testFunc()))
	//app.Command(name)
}

func createApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Inductor"
	app.Usage = "Generate Packer Templates"
	app.Commands = []cli.Command{
		{
			Name:      "generate",
			ShortName: "g",
			Usage:     "Generate a Packer template",
			Action: func(c *cli.Context) {
				println("Generating:", c.String("osfamily")+" "+c.String("osversion")+" "+c.String("osedition"))
			},
			Flags: []cli.Flag{
				cli.StringFlag{"osfamily, os", "windows", "Operating System Family (e.g. windows)", false},
				cli.StringFlag{"osversion, ver", "2012-r2", "Operating System Version (e.g. 2012-r2)", false},
				cli.StringFlag{"osedition, ed", "standard", "Operating System Edition (e.g. standard)", false},
				cli.StringFlag{"communicator, c", "ssh", "Packer Communicator (e.g. ssh or winrm)", false},
				cli.BoolFlag{"skipwindowsupdates, u", "Skip Installation Of Windows Updates"},
				cli.BoolFlag{"skipguesttools", "Skip Installation Of VM Guest Tools"},
				cli.BoolFlag{"chef", "Install Chef"},
				cli.BoolFlag{"puppet", "Install Puppet"},
				cli.BoolFlag{"test", "Test"},
			},
		},
	}

	return app
}

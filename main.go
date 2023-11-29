package main

import (
	"flag"
	"fmt"
	"github.com/SecurityNeo/NexusbackupTool/service/backup"
	"github.com/SecurityNeo/NexusbackupTool/service/recover"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Parameters struct {
	action      string
	data        string
	help        bool
	parallelism int
}

var Flag Parameters

func init() {

	flag.BoolVar(&Flag.help, "h", false, "Print help.")
	flag.StringVar(&Flag.action, "a", "backup", "The action of the tool,must be backup or recover.")
	flag.StringVar(&Flag.data, "d", "config.json", "The configuration file,should specify the "+
		"full path.")
	flag.IntVar(&Flag.parallelism, "parallelism", 5, "Limit the number of parallel resource operations.")

	flag.Usage = usage

}

func usage() {
	fmt.Fprintf(os.Stderr, `Nexus backup and recover tool version: 0.1 auther: neo
Usage: NexusTool [-adh]
 
Options:
`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	if Flag.help {
		flag.Usage()
		os.Exit(0)
	}

	var configDirectory string
	if Flag.data != "" {
		configDirectory = Flag.data
	} else {
		configDirectory = GetCurrentDirectory() + "config.json"
	}

	if Flag.action == "backup" {
		err := backup.StartBackup(configDirectory, Flag.parallelism)
		if err != nil {
			log.Fatalln("Error: ", err)
			os.Exit(20)
		}
	} else if Flag.action == "recover" {
		err := recover.StartRecover(configDirectory, Flag.parallelism)
		if err != nil {
			log.Fatalln("Error: ", err)
			os.Exit(20)
		}
	}

}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

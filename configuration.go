package main

import (
	"flag"
	"github.com/gobwas/glob"
	log "github.com/sirupsen/logrus"
	"time"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type Configuration struct {
	Directories     arrayFlags
	Exclusions      []glob.Glob
	ListenAddress   string
	WalkingInterval time.Duration
}

func (config *Configuration) Parse() {
	flag.Var(&config.Directories, "dir", "The directory to be explored by the exporter to compute metrics.\nCan be provided multiple times.")
	var exclusions arrayFlags
	flag.Var(&exclusions, "exclude", "Glob pattern to be excluded from the directories specified.\nCan be provided multiple times.")
	flag.StringVar(&config.ListenAddress, "listen-address", ":9123", "The address to listen on for HTTP requests.")
	flag.DurationVar(&config.WalkingInterval, "walking-interval", 1*time.Minute, "The interval to walk the directories and compute the metrics.")
	flag.Parse()

	if len(config.Directories) == 0 {
		log.Fatalf("At least one directory must be provided.")
	}

	if len(exclusions) == 0 {
		log.Warnf("No exclusions provided. All files will be considered.")
	}

	for _, exclusion := range exclusions {
		globPattern, err := glob.Compile(exclusion)
		if err != nil {
			log.Errorf("Ignoring %s glob pattern, error compiling the expression: %v", exclusion, err)
		}
		config.Exclusions = append(config.Exclusions, globPattern)
	}
}

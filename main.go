package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jbarone/ghostToHugo/ghosttohugo"

	jww "github.com/spf13/jwalterweatherman"
	flag "github.com/spf13/pflag"
)

// Print usage information
func usage() {
	fmt.Printf("Usage: %s [OPTIONS] <Ghost Export>\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	var (
		path, loc, format              string
		subdirs, force, verbose, debug bool
	)

	flag.Usage = usage

	flag.StringVarP(&path, "hugo", "p", "newhugosite",
		"path to create the new hugo project")
	flag.BoolVarP(&subdirs, "subdirs", "s", false,
		"Export posts and pages as index.md within their own named sub-directories.")
	flag.StringVarP(&loc, "location", "l", "",
		"location to use for time conversions (default: local)")
	flag.StringVarP(&format, "dateformat", "d", "2006-01-02 15:04:05",
		"date format string to use for time conversions")
	flag.BoolVarP(&force, "force", "f", false,
		"allow import into non-empty target directory")
	flag.BoolVarP(&verbose, "verbose", "v", false,
		"print verbose logging output")
	flag.BoolVarP(&debug, "debug", "", false,
		"print verbose logging output")

	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(0)
	}

	opts := []func(*ghosttohugo.Converter){
		ghosttohugo.WithHugoPath(path),
	}
	if loc != "" {
		location, err := time.LoadLocation(loc)
		if err != nil {
			jww.FATAL.Fatalf("Error loading location %s: %v\n", loc, err)
		}
		opts = append(opts, ghosttohugo.WithLocation(location))
	}

	if format != "" {
		opts = append(opts, ghosttohugo.WithDateFormat(format))
	}

	if subdirs {
		opts = append(opts, ghosttohugo.WithSubDirs())
	}

	if force {
		opts = append(opts, ghosttohugo.WithForce())
	}

	c, err := ghosttohugo.New(opts...)
	if err != nil {
		jww.FATAL.Fatalf("Error initializing converter (%v)\n", err)
	}

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		jww.FATAL.Fatalf("Error opening export: %v\n", err)
	}
	defer file.Close()

	// setup logging
	lvl := jww.LevelWarn
	if verbose {
		lvl = jww.LevelInfo
	}
	if debug {
		lvl = jww.LevelDebug
	}
	jww.SetStdoutThreshold(lvl)

	jww.FEEDBACK.Println("Importing...")

	count, err := c.Convert(file)
	if err != nil {
		jww.FATAL.Fatalf("Error opening export: %v\n", err)
	}

	jww.FEEDBACK.Printf("Congratulations! %d post(s) imported!\n", count)
	jww.FEEDBACK.Printf("Now, start Hugo by yourself:\n"+
		"$ git clone https://github.com/spf13/herring-cove.git "+
		"%s/themes/herring-cove\n", path)
	jww.FEEDBACK.Printf("$ cd %s\n$ hugo server --theme=herring-cove\n", path)
}

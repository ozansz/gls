package main

import (
	"flag"
	"os"
	"strings"
	"sync"

	"go.sazak.io/gls/gui"
	"go.sazak.io/gls/internal"
	"go.sazak.io/gls/internal/fs"
	"go.sazak.io/gls/internal/local"
	"go.sazak.io/gls/internal/types"
	"go.sazak.io/gls/log"

	"github.com/rivo/tview"
)

const (
	logFile = "gls.log"
)

var (
	path          = flag.String("path", "", "path to run on (required)")
	formatter     = flag.String("fmt", "bytes", "size formatter, one of bytes, pow10 or none")
	noGUI         = flag.Bool("nogui", false, "text-only mode")
	sort          = flag.Bool("sort", true, "sort nodes by size")
	sizeThreshold = flag.String("thresh", "", "size filter threshold, e.g. 10M, 100K, etc.")
	ignoreFiles   = flag.String("ignore", "", "Comma-separated ignore files that specify which files/folders to exclude")
	debug         = flag.Bool("debug", false, "Increase log verbosity")

	formatters = map[string]types.SizeFormatter{
		"bytes": types.SizeFormatterBytes,
		"pow10": types.SizeFormatterPow10,
		"none":  types.NoFormat,
	}
)

func main() {
	flag.Parse()
	if *path == "" {
		flag.Usage()
		return
	}
	if *debug {
		log.SetDebug(1)
	}
	logF, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer func() {
		if err := logF.Close(); err != nil {
			panic(err)
		}
	}()
	formatterFunc, ok := formatters[*formatter]
	if !ok {
		log.Errorf("Unknown formatter: %s", *formatter)
		flag.Usage()
		return
	}
	ignoreChecker, err := getIgnoreChecker()
	if err != nil {
		log.Fatalf("Failed to get ignore checker: %v", err)
	}
	var sizeThreshBytes int64 = 0
	if *sizeThreshold != "" {
		byteSize, mult, err := internal.ParseByteSize(*sizeThreshold)
		if err != nil {
			log.Errorf("Failed to parse size threshold: %s", *sizeThreshold)
			return
		}
		if mult <= 0 {
			log.Errorf("Size threshold cannot be less than or equal to zero: %s", *sizeThreshold)
			return
		}
		sizeThreshBytes = int64(byteSize) * mult
	}
	log.Infof("Starting gls with path: %s, log file: %s, formatter: %s, gui: %t", *path, logFile, *formatter, !*noGUI)
	if !*noGUI {
		log.SetOutput(logF)
		log.Infof("Started gls with path: %s, log file: %s, formatter: %s, gui: %t", *path, logFile, *formatter, !*noGUI)
	}
	log.Debugf("Ignore checking rules:\n%s", ignoreChecker.Dump())
	var app *tview.Application
	if !*noGUI {
		app = gui.GetApp(*path, formatterFunc)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		opts := []fs.FileTreeBuilderOption{
			fs.WithSizeFormatter(formatterFunc),
			fs.WithIgnoreChecker(ignoreChecker),
		}
		if *sort {
			opts = append(opts, fs.WithSortingBySize())
		}
		if sizeThreshBytes > 0 {
			opts = append(opts, fs.WithSizeThreshold(sizeThreshBytes))
		}
		b := fs.NewFileTreeBuilder(*path, opts...)
		if err := b.Build(); err != nil {
			log.Fatalf("Failed to build file tree: %v", err)
			return
		}
		log.Info("Finished building file tree")
		if *noGUI {
			if err := b.Print(); err != nil {
				log.Fatalf("Error while printing the file tree: %v\n", err)
			}
			return
		}
		if !*noGUI {
			log.Info("Loading tree view on GUI")
			gui.LoadTreeView(app, b.Root(), *path)
			log.Info("Loaded the tree view on GUI")
		}
	}()
	if !*noGUI {
		if err := app.Run(); err != nil {
			log.Fatalf("Error running GUI app: %v", err)
		}
	}
	if *noGUI {
		wg.Wait()
	}
}

func getIgnoreChecker() (*local.IgnoreChecker, error) {
	ignoreCheckerOpts := []local.IgnoreCheckerOption{}
	if *ignoreFiles != "" {
		for _, path := range strings.Split(*ignoreFiles, ",") {
			ignoreCheckerOpts = append(ignoreCheckerOpts, local.WithRuleFile(path))
		}
	}
	return local.NewIgnoreChecker(ignoreCheckerOpts...)
}

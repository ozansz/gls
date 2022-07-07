package main

import (
	"flag"
	"os"
	"sync"

	"github.com/ozansz/gls/gui"
	"github.com/ozansz/gls/internal"
	"github.com/ozansz/gls/log"

	"github.com/rivo/tview"
)

const (
	logFile = "gls.log"
)

var (
	path          = flag.String("path", "", "path to list")
	formatter     = flag.String("fmt", "bytes", "formatter to use: bytes, pow10 or none")
	noGUI         = flag.Bool("nogui", false, "do not show GUI")
	sort          = flag.Bool("sort", true, "sort nodes by size")
	sizeThreshold = flag.String("thresh", "", "size filter threshold")

	formatters = map[string]internal.SizeFormatter{
		"bytes": internal.SizeFormatterBytes,
		"pow10": internal.SizeFormatterPow10,
		"none":  internal.NoFormat,
	}
)

func main() {
	flag.Parse()
	if *path == "" {
		flag.Usage()
		return
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
	var app *tview.Application
	if !*noGUI {
		app = gui.GetApp(*path, formatterFunc)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		opts := []internal.FileTreeBuilderOption{
			internal.WithSizeFormatter(formatterFunc),
		}
		if *sort {
			opts = append(opts, internal.WithSortingBySize())
		}
		if sizeThreshBytes > 0 {
			opts = append(opts, internal.WithSizeThreshold(sizeThreshBytes))
		}
		b := internal.NewFileTreeBuilder(*path, opts...)
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

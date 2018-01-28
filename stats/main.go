package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"word-statistics/stats/processor"
	"word-statistics/stats/statistics"
)

const (
	wordsFile   = "data.txt"
	defaultPort = 8080
)

// errIsDirectory error returned when the supplied file is a directory.
var errIsDirectory = errors.New("file is directory")

// server which listens for a GET request on `/stats` which reads the raw data stored at
// the file supplied in flags and generates statistics on the contents.
// The file pointer and stats are updated every 10 seconds.
func main() {
	file := flag.String("f", filepath.Join(os.TempDir(), wordsFile), "location to read received words from")
	port := flag.Int("p", defaultPort, "port to listen on")
	flag.Parse()

	f, err := openFile(*file)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// perform an initial process.
	p := newProcessor(ctx, f)
	f.Close()
	go func(f *os.File) {
		ticker := time.NewTicker(10 * time.Second)
		// perform initial process.
		for {
			select {
			case <-ticker.C:
				// re-process every 10 seconds to update stats.
				f, err := openFile(*file)
				if err != nil {
					panic(err)
				}
				p = newProcessor(ctx, f)
				f.Close()
			}
		}
	}(f)

	router := mux.NewRouter()
	// write the process results to the stats endpoint.
	router.HandleFunc("/stats", func(writer http.ResponseWriter, request *http.Request) {
		p.Write(writer)
	})

	http.ListenAndServe(fmt.Sprintf(":%v", *port), router)
}

// openFile helper function to open a file.
func openFile(path string) (*os.File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	if stat, err := f.Stat(); err != nil || stat.IsDir() {
		if err != nil {
			err = errIsDirectory
		}
		return nil, err
	}

	return f, nil
}

// newProcessor initialises and processes the provided reader.
func newProcessor(ctx context.Context, r io.Reader) processor.Processor {
	p := processor.New(ctx)
	p.Register(statistics.NewTopLetters())
	p.Register(statistics.NewTopFive())
	p.Register(&statistics.Count{})
	p.Process(r)

	return p
}

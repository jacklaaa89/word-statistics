package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"io/ioutil"

	"github.com/gorilla/mux"
)

const (
	wordsFile   = "data.txt"
	defaultPort = 5555
)

// errIsDirectory error returned when the supplied file is a directory.
var errIsDirectory = errors.New("file is directory")

// very quick server which reads the raw body of a post request and appends
// it to the file supplied in the flags.
func main() {
	file := flag.String("f", filepath.Join(os.TempDir(), wordsFile), "location to read received words from")
	port := flag.Int("p", defaultPort, "port to listen on")
	flag.Parse()

	router := mux.NewRouter()
	// write the process results to the stats endpoint.
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		data, _ := ioutil.ReadAll(request.Body)
		f, err := openFile(*file)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		fmt.Fprintln(f, string(data))
	}).Methods(http.MethodPost)

	http.ListenAndServe(fmt.Sprintf(":%v", *port), router)
}

// openFile helper function to open a file.
func openFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dw3105/siduri-web/internal/site"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: siduri build | dev")
	}

	switch os.Args[1] {
	case "build":
		build(os.Args[2:], false)
	case "dev":
		dev(os.Args[2:])
	default:
		log.Fatalf("unknown command %q; use build or dev", os.Args[1])
	}
}

func build(args []string, includeDrafts bool) {
	flags := flag.NewFlagSet("build", flag.ExitOnError)
	output := flags.String("output", "dist", "directory for generated files")
	_ = flags.Parse(args)

	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if err := site.Build(root, filepath.Clean(*output), includeDrafts); err != nil {
		log.Fatal(err)
	}
}

func dev(args []string) {
	flags := flag.NewFlagSet("dev", flag.ExitOnError)
	output := flags.String("output", ".dev-dist", "directory for generated preview files")
	addr := flags.String("addr", "127.0.0.1:8080", "address for the development server")
	_ = flags.Parse(args)

	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if err := site.Build(root, filepath.Clean(*output), true); err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.Dir(*output))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Robots-Tag", "noindex")
		fileServer.ServeHTTP(w, r)
	})
	fmt.Printf("Siduri preview: http://%s/journal/hello-siduri/\n", *addr)
	log.Printf("serving %s on %s", *output, *addr)
	if err := http.ListenAndServe(*addr, handler); err != nil {
		log.Fatal(err)
	}
}

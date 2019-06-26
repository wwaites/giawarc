package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	"github.com/wwaites/giawarc"
)

var outdir string
var outform string

func init() {
	flag.StringVar(&outdir, "o", ".", "Output location")
	flag.StringVar(&outform, "f", "bitextor", "Output format")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] WARCFile\nFlags:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(),
`Formats:
  bitextor
        Output format compatible with bitextor (circa June 2019)
  rocks
        Concatenated gzipped pages indexed with RocksDB
  rockslang
        Concatenated gzipped pages split by language and indexed with RocksDB
`)
	}
}

func PreProcessFile(filename string) (proc *giawarc.WARCPreProcessor, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	var tw giawarc.TextWriter
	if outform == "bitextor" {
		tw, err = giawarc.NewBitextorWriter(outdir)
		if err != nil {
			return
		}
	} else if outform == "gzip" {
		tw, err = giawarc.NewZipWriter(outdir)
		if err != nil {
			return
		}
	} else if outform == "gzlang" {
		m := func(o string) (giawarc.TextWriter, error) { return giawarc.NewZipWriter(o) }
		tw, err = giawarc.NewLangWriter(outdir, m)
		if err != nil {
			return
		}
	} else {
		fmt.Fprintf(flag.CommandLine.Output(), "Unknown output format %s\n", outform)
		os.Exit(1)
	}
	defer tw.Close()

	proc, err = giawarc.NewWARCPreProcessor(f, tw)
	if err != nil {
		return
	}

	proc.Process()

	return
}

func main() {

	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	filename := flag.Arg(0)

	start := time.Now()
	proc, err := PreProcessFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	end := time.Now()

	elapsed := end.Sub(start)

	fmt.Printf("total records: %v\n", proc.TotalRecords)
	fmt.Printf("text records: %v\n", proc.TextRecords)
	fmt.Printf("lang records: %v\n", proc.LangRecords)
	fmt.Printf("total bytes: %v\n", proc.TotalBytes)
	fmt.Printf("text bytes: %v\n", proc.TextBytes)
	fmt.Printf("lang bytes: %v\n", proc.LangBytes)
	fmt.Printf("elapsed time: %v\n", elapsed)
	fmt.Printf("content types:\n")

/*
	cts := proc.ContentTypeStats()
	for _, s := range cts {
		fmt.Printf("    %v: %0.08f\n", s.ContentType, s.Prevalence)
	}
*/
}
package main

import (
	"flag"
	"fmt"
	"github.com/dustin/go-humanize"
	"math"
	"math/rand"
	"os"
	"path/filepath"
)

var (
	fileSizes     []float64
	baseDirectory string
	directories   Directories
)

type randByteMaker struct {
	src rand.Source
}

func (r *randByteMaker) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = byte(r.src.Int63() & 0xff) // mask to only the first 255 byte values
	}
	return len(p), nil
}

// Directories is a named type of []string that is a Sortable so that we can
// sort this in the order of length of the string, with the longest strings first
type Directories []string

func (a Directories) Len() int           { return len(a) }
func (a Directories) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Directories) Less(i, j int) bool { return len(a[i]) > len(a[j]) }

func main() {

	flag.Usage = func() {
		flag.PrintDefaults()
	}
	flag.Parse()

	baseDirectory = flag.Arg(0)
	if info, err := os.Stat(baseDirectory); err != nil || !info.IsDir() {
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "%s is not a directory\n\n", baseDirectory)
		}
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println("Directory stats")
	if err := filepath.Walk(baseDirectory, fileWalker); err != nil {
		fmt.Printf("\n%s\n", err)
		os.Exit(1)
	}

	var totalFilesize float64
	biggest := math.SmallestNonzeroFloat64
	smallest := math.MaxFloat64
	for i := 0; i < len(fileSizes); i++ {
		if fileSizes[i] == 0 {
			continue
		}
		if fileSizes[i] < smallest {
			smallest = fileSizes[i]
		}
		if fileSizes[i] > biggest {
			biggest = fileSizes[i]
		}
		totalFilesize += fileSizes[i]
	}

	fmt.Printf("\nfiles: %d\n", len(fileSizes))
	fmt.Printf("directories: %d\n", len(directories))

	fmt.Printf("total filesize: %s\n", humanize.Bytes(uint64(totalFilesize)))

	mean := totalFilesize / float64(len(fileSizes))
	fmt.Printf("mean filesize: %s\n", humanize.Bytes(uint64(mean)))
	fmt.Printf("largest filesize: %s\n", humanize.Bytes(uint64(biggest)))
	fmt.Printf("smallest filesize: %s\n", humanize.Bytes(uint64(smallest)))

	var sum float64
	for i := 0; i < len(fileSizes); i++ {
		tmp := fileSizes[i] - mean
		sum += tmp * tmp
	}

	t := math.Sqrt(sum / float64(len(fileSizes)))
	fmt.Printf("standard deviation: %.2f\n", t)
}

func fileWalker(filePath string, f os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if f.IsDir() {
		directories = append(directories, filePath)
		return nil
	}
	mode := f.Mode()
	if !mode.IsRegular() {
		fmt.Printf("%s skipped because not regular file\n", filePath)
		return nil
	}
	fileSizes = append(fileSizes, float64(f.Size()))
	return nil
}

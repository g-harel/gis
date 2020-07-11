package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/g-harel/gis"
)

func main() {
	// TODO make root and dest cli args.
	root := path.Join(os.Getenv("GOROOT"), "src")
	path := path.Join(os.Getenv("GOPATH"), "src")
	dest := "./.index"
	// TODO move query to different script + server.
	query := "ios reder option"

	fmt.Printf("ROOT=%v\n", root)
	fmt.Printf("PATH=%v\n", path)
	fmt.Printf("DEST=%v\n", dest)
	fmt.Printf("QUERY=%v\n", query)
	fmt.Println("========")

	indexTime := time.Now()
	idx := gis.NewSearchIndex()
	err := idx.Include(root)
	if err != nil {
		panic(err)
	}
	err = idx.Include(path)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed in %s\n", time.Since(indexTime))

	encodeTime := time.Now()
	var buf bytes.Buffer
	err = idx.ToBytes(&buf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Encoded in %s\n", time.Since(encodeTime))

	decodeTime := time.Now()
	idx, err = gis.NewSearchIndexFromBytes(&buf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decoded in %s\n", time.Since(decodeTime))

	writeTime := time.Now()
	f, err := os.Create(dest)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = idx.ToBytes(f)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Written in %s\n", time.Since(writeTime))

	searchStart := time.Now()
	interfaces, err := idx.Search(query)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Queried in %s\n", time.Since(searchStart))
	fmt.Println("========")

	for _, ifc := range interfaces[:16] {
		println(ifc.String())
	}
}

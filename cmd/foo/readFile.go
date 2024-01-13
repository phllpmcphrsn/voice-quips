package main

import (
	"fmt"
	"os"

	"github.com/dhowden/tag"
)

func main() {
	// open input file
	file, err := os.Open("G:\\Users\\phllp\\go\\github.com\\phllpmcphrsn\\voice-uploader\\audio_files\\voice-quip-1.mp3")
	if err != nil {
		panic(err)
	}

	// close file on exit and check for its returned error
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	
	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s %d %+v %s %+v\n", stat.Name(), stat.Size(), stat.Sys(), stat.ModTime().String(), stat.Mode())

	// buffer := make([]byte, 1024)
	// for {
	// 	// read a chunk
	// 	n, err := file.Read(buffer)
	// 	if err != nil && err != io.EOF {
	// 		panic(err)
	// 	}
	// 	if n == 0 {
	// 		break
	// 	}

	// 	// write a chunk
	// 	fmt.Print(n)
	// }

	println()
	m, err := tag.ReadFrom(file)
	if err != nil {
		println("error occurred with tag.ReadFrom() ", err.Error())
		return
	}

	println("title: ", m.Title())
	println("artist: ", m.Artist())
	println("album: ", m.Year())
	print(m.FileType())
	print(m.)
}
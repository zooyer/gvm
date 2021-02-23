package main

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"time"
)

func main() {
	count := 10000
	// create and start new bar
	//bar := pb.StartNew(count)

	// start bar from 'default' template
	//bar := pb.Default.Start(count)

	// start bar from 'simple' template
	//bar := pb.Simple.Start(count)

	// start bar from 'full' template
	bar := pb.Full.Start(count)
	bar.

	fmt.Println("downloading:")
	for i := 0; i < count; i++ {
		bar.Add(1)
		time.Sleep(time.Millisecond)
	}
	bar.Finish()
}
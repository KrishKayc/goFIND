package main

import (
	"flag"
	"fmt"
	"gofind/finder"
)

func main() {

	dir := flag.String("dir", "", "the directory path")
	searchText := flag.String("search", "", "the text to search")

	flag.Parse()

	if *dir == "" || *searchText == "" {
		fmt.Println("Provide a valid folder and search text")
		return
	}

	//set the configuration first...
	finder.Init("config.json")

	//both name and path are same for parent directory..
	root := finder.NewDir(*dir, *dir)

	//finds the search text inside the directory..\
	root.Find(*searchText)
}

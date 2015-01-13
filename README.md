youtube-go
=========
[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/adrg/youtube-go)
[![License: MIT](http://img.shields.io/badge/license-MIT-red.svg?style=flat-square)](http://opensource.org/licenses/MIT)

Small package for searching Youtube videos and getting information about
them. It can read video formats and also does signature descrambling so the
library can be easily used to download videos.

Full documentation can be found at: http://godoc.org/github.com/adrg/youtube-go

## Installation
```
go get github.com/adrg/youtube-go
```

### Usage

#### Get video information
```go
package main

import (
	"fmt"

	youtube "github.com/adrg/youtube-go"
)

func main() {
	video, err := youtube.GetVideo("uFknAPhEcQM")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Title: %s\n", video.Title)
	fmt.Printf("Id: %s\n", video.Id)
	fmt.Printf("Duration: %d\n", video.Duration)
	fmt.Printf("Description: %s\n", video.Description)
	fmt.Printf("View count: %d\n", video.ViewCount)
	fmt.Printf("Rating: %f\n", video.Rating)

	fmt.Println("Formats:")
	for _, f := range video.Formats {
		fmt.Printf("  * %d - %s %s ===\n", f.Itag, f.Type, f.Quality)
		fmt.Printf("    URL: %s\n", f.URL)
	}
}
```

#### Search for videos
```go
package main

import (
	"fmt"

	youtube "github.com/adrg/youtube-go"
)

func main() {
	// The search parameters include the query, the page you want
	// to get results for and the maximum results allowed
	searchParams := &youtube.SearchParams{"amv", 1, 3}
	videos, err := youtube.Search(searchParams)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, video := range videos {
		if err = video.ReadFormats(); err != nil {
			fmt.Println(err)
			// Handle error here
		}

		fmt.Printf("Title: %s\n", video.Title)
		fmt.Printf("Id: %s\n", video.Id)
		fmt.Printf("Duration: %d\n", video.Duration)
		fmt.Printf("Description: %s\n", video.Description)
		fmt.Printf("View count: %d\n", video.ViewCount)
		fmt.Printf("Rating: %f\n", video.Rating)

		fmt.Println("Formats:")
		for _, f := range video.Formats {
			fmt.Printf("  * %d - %s %s ===\n", f.Itag, f.Type, f.Quality)
			fmt.Printf("    URL: %s\n", f.URL)
		}

		fmt.Println()
	}
}
```

## License
Copyright (c) 2015 Adrian-George Bostan.

This project is licensed under the [MIT license](http://opensource.org/licenses/MIT). See LICENSE for more details.

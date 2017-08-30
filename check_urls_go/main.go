package main

import (
	"fmt"
)

type Site struct {
	label    string
	url      string
	expected string
}

func main() {

	/* initialize array (https://stackoverflow.com/questions/26159416/init-array-of-structs-in-go)
	 */

	var sites = []Site{
		Site{
			label:    "repo_file",
			url:      "https://repository.library.brown.edu/storage/bdr:6758/PDF/",
			expected: "BleedBox",
		},
		Site{
			"repo_search", "https://repository.library.brown.edu/studio/search/?q=elliptic", "The sequence of division polynomials"},
	}


	fmt.Println(sites)
	fmt.Println(sites[0])

} // end func main()

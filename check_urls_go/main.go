package main

import (
	"fmt"
    "github.com/davecgh/go-spew/spew"
)

type Site struct {
	label    string
	url      string
	expected string
}

func main() {

	/* initialize array (https://stackoverflow.com/questions/26159416/init-array-of-structs-in-go) */

	var sites = []Site{
		Site{
			label:    "repo_file",
			url:      "https://repository.library.brown.edu/storage/bdr:6758/PDF/",
			expected: "BleedBox", // note: since brace is on following line, this comma is required
		},
		Site{"repo_search",
			"https://repository.library.brown.edu/studio/search/?q=elliptic",
			"The sequence of division polynomials"},
		Site{"bipg_wiki",
			"https://wiki.brown.edu/confluence/display/bipg/Brown+Internet+Programming+Group+Home",
			"The BIPG idea"},
		Site{"booklocator_app",
			"http://library.brown.edu/services/book_locator/?callnumber=GC97+.C46&location=sci&title=Chemistry+and+biochemistry+of+estuaries&status=AVAILABLE&oclc_number=05831908&public=true",
			"GC97 .C46 Level 11, Aisle 2A"},
		Site{"callnumber_app",
			"https://apps.library.brown.edu/callnumber/v2/?callnumber=PS3576",
			"American Literature"},
		Site{"clusters api",
			"https://library.brown.edu/clusters_api/data/",
			"scili-friedman"},
		Site{"easyborrow_feed",
			"http://library.brown.edu/easyborrow/feeds/latest_items/",
			"easyBorrow -- recent requests"},
		Site{"freecite",
			"http://freecite.library.brown.edu/welcome/",
			"About FreeCite"},
		Site{"iip_inscriptions",
			"http://library.brown.edu/cds/projects/iip/viewinscr/abur0001/",
			"Khirbet Abu Rish"},
		Site{"iip_processor",
			"https://apps.library.brown.edu/iip_processor/info/",
			"hi"},
	}

	fmt.Println(sites)
    fmt.Println( "---" )
	// fmt.Println(sites[0])
	fmt.Printf("%+v\n", sites) // adding `%+v` prints the field-names
    fmt.Println( "---" )
    spew.Dump(sites)
    fmt.Println( "---" )

} // end func main()

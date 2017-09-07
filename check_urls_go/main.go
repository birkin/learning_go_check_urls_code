package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"net/http"
    "reflect"
	"strings"
	"time"
)

type Site struct {
	label    string
	url      string
	expected string
}

type Result struct {
	label        string
	check_result string
    time_taken time.Duration
}

var sites []Site     // i think this is declaring a slice, not an array
var results []Result // same as above

func main() {
	initialize_sites() // initialize array (https://stackoverflow.com/questions/26159416/init-array-of-structs-in-go)
	check_sites(sites) // do the work
	/* print stuff */
	// fmt.Println("plain sites -- ", sites)
	// fmt.Println("---")
	// mt.Printf("sites with labels, -- %+v\n", sites) // adding `%+v` prints the field-names
	// fmt.Println("---")
	// fmt.Println("dump...")
	// spew.Dump(sites)
	// fmt.Println("---")

} // end func main()

/* ----------------------------------------------------------------------
   helper functions
   ---------------------------------------------------------------------- */

func initialize_sites() []Site {
	sites = []Site{}
	sites = append(
		sites,
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
		Site{"not_found_test",
			"https://apps.library.brown.edu/iip_processor/info/",
			"foo"},
	)
	return sites
}

func check_sites(sites []Site) {
	total_start := time.Now()
	for _, site_element := range sites {
		start := time.Now()
		fmt.Println("\nsite -- ", site_element.label)
		resp, _ := http.Get(site_element.url)
		fmt.Println("status code -- ", resp.StatusCode)
		body_bytes, _ := ioutil.ReadAll(resp.Body)
		text := string(body_bytes)
		var site_check_result string = "not_found"
		if strings.Contains(text, site_element.expected) {
			site_check_result = "found"
		}
		fmt.Println("site_check_result, ", site_check_result)

		elapsed := time.Since(start)
        fmt.Println("elapsed, ", elapsed)

        fmt.Println("elapsed has TypeOf: ", reflect.TypeOf(elapsed))
        elapsed_k := reflect.ValueOf(elapsed)
        fmt.Println("elapsed has Kind: ", elapsed_k.Kind())

		result_instance := Result{
			label:        site_element.label,
			check_result: site_check_result,
            time_taken: elapsed,
		}
		fmt.Println("result_instance.check_result, ", result_instance.check_result)
		results = append(results, result_instance)

	}
	// resp, _ := http.Get("https://library.brown.edu/bjd/internationalization.html")
	// fmt.Println("response -- ", resp)
	// fmt.Println("status code -- ", resp.StatusCode)
	// fmt.Println("body -- ", resp.Body)
	// fmt.Println("body2 -- ", resp.Body)

	// body_bytes, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("body_bytes -- ", body_bytes)
	// text := string(body_bytes)
	// fmt.Println("body -- ", text)
	total_elapsed := time.Since(total_start)
	fmt.Println("total_elapsed, ", total_elapsed)
	fmt.Println("\n results...")
	spew.Dump(results)
	// above may not handle non-ascii characters: <https://stackoverflow.com/a/38808838> -- update, it appears to handle non-ascii characters fine
}

// defer timeTrack(time.Now(), "lookup-and-check")

// func timeTrack(start time.Time, name string) {
// 	elapsed := time.Since(start)
// 	fmt.Printf("%s took %s\n", name, elapsed)
// }

// if strings.Contains(str, subStr) {}

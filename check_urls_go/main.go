package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/kelseyhightower/envconfig"
)

type Settings struct {
	LOGPATH  string `envconfig:"LOGPATH" default:"oops-default-logpath"`
	APPTITLE string `default:"oops-default-apptitle"`
}

type Site struct {
	label    string
	url      string
	expected string
}

type Result struct {
	label        string
	check_result string
	time_taken   time.Duration
}

var settings Settings
var sites []Site // i think this declares a slice, not an array
// var results []Result // same as above

func main() {
	/* Loads settings, initializes sites array, calls worker function. */

	/// initialize settings
	fmt.Printf("LOGPATH in main() before settings initialized, ```%v```\n", settings.LOGPATH)
	load_settings()
	fmt.Printf("LOGPATH in main() after settings initialized, ```%v```\n", settings.LOGPATH)

	/// initialize sites array
	initialize_sites() // (https://stackoverflow.com/questions/26159416/init-array-of-structs-in-go)

	/// call worker function
	check_sites_just_with_routines(sites)

} // end func main()

/* ----------------------------------------------------------------------
   helper functions
   ---------------------------------------------------------------------- */

func load_settings() Settings {
	/* Loads settings, eventually for logging and database. */
	err := envconfig.Process("url_check_", &settings) // env settings look like `URL_CHECK__THE_SETTING`
	if err != nil {
		fmt.Printf("error, ```%v```", err.Error)
		panic(err)
	}
	fmt.Printf("LOGPATH in load_settings(), ```%v```\n", settings.LOGPATH)
	return Settings{}
}

func initialize_sites() []Site {
	/* Populates sites slice. */
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
	fmt.Println("\n sites, spewed...")
	spew.Dump(sites)
	return sites
}

func check_sites_just_with_routines(sites []Site) {
	/* Creates channel, kicks off go-routines, prints channel output, and closes channel. */

	/// initialize channel
	writer_channel := make(chan string)

	/// start go routines
	for _, site_element := range sites {
		go check_site(site_element, writer_channel)
	}

	/// output channel data
	var counter int
	var channel_output string
	for channel_output = range writer_channel {
		counter++
		time.Sleep(50 * time.Millisecond)
		fmt.Println("channel-value, ", channel_output)
		if counter == len(sites) {
			fmt.Println("about to close")
			close(writer_channel)
		}
	}
}

func check_site(site Site, writer_channel chan string) {
	/* Checks site, stores data to result, & writes info to channel. */

	start := time.Now()

	/// check site
	resp, _ := http.Get(site.url)
	body_bytes, _ := ioutil.ReadAll(resp.Body)
	text := string(body_bytes)
	var site_check_result string = "not_found"
	if strings.Contains(text, site.expected) {
		site_check_result = "found"
	}

	elapsed := time.Since(start)

	/// store result
	result_instance := Result{
		label:        site.label,
		check_result: site_check_result,
		time_taken:   elapsed,
	}

	/// write info to channel
	fmt.Println("result_instance.label, ", result_instance.label)
	writer_channel <- result_instance.label
}

/// EOF

package main

import (
	// "bytes"
	"fmt"
	"io/ioutil"
	// "log"
	"net/http"
	"strings"
	"time"

	// "github.com/davecgh/go-spew/spew"      // easy way to pretty-print structs
	"github.com/kelseyhightower/envconfig" // binds env vars to settings struct
	"github.com/romana/rlog"
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
	rlog.Info("\n\nstarting")

	/// initialize settings
	rlog.Debug(fmt.Sprintf("LOGPATH in main() before settings initialized, ```%v```", settings.LOGPATH))
	load_settings()

	/// initialize sites array
	initialize_sites() // (https://stackoverflow.com/questions/26159416/init-array-of-structs-in-go)

	/// call worker function
	check_sites_with_goroutines(sites)

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
	rlog.Debug(fmt.Sprintf("LOGPATH in load_settings(), ```%v```", settings.LOGPATH))
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
	rlog.Info(fmt.Sprintf("sites to process, ```%#v```", sites)) // prints, eg, `{label:"clusters api", url:"etc...`

	return sites
}

func check_sites_with_goroutines(sites []Site) {
	/* Creates channel, kicks off go-routines, prints channel output, and closes channel. */

	main_start := time.Now()

	/// initialize channel
	writer_channel := make(chan Result)

	/// start go routines
	for _, site_element := range sites {
		go check_site(site_element, writer_channel)
	}

	/// output channel data
	var counter int = 0
	var channel_output Result
	for channel_output = range writer_channel {
		counter++
		time.Sleep(50 * time.Millisecond)
		rlog.Info(fmt.Sprintf("channel-value, ```%v```", channel_output))
		if counter == len(sites) {
			rlog.Info("about to close channel")
			close(writer_channel)
			rlog.Info("channel closed")
			break // shouldn't be needed
		}
	}
	main_elapsed := time.Since(main_start)
	rlog.Info(fmt.Sprintf("main_elapsed, ```%v```", main_elapsed) )

}

func check_site(site Site, writer_channel chan Result) {
	/* Checks site, stores data to result, & writes info to channel. */

	/// check site
	mini_start := time.Now()
	resp, _ := http.Get(site.url)
	body_bytes, _ := ioutil.ReadAll(resp.Body)
	text := string(body_bytes)
	var site_check_result string = "not_found"
	if strings.Contains(text, site.expected) {
		site_check_result = "found"
	}

	/// store result
	mini_elapsed := time.Since(mini_start)
	result_instance := Result{
		label:        site.label,
		check_result: site_check_result,
		time_taken:   mini_elapsed,
	}

	/// write info to channel
	writer_channel <- result_instance
	rlog.Info(fmt.Sprintf("result_instance after write to channel, ```%#v```", result_instance))
}

/// EOF

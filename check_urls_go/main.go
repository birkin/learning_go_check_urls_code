package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/kelseyhightower/envconfig"
)

type Settings struct {
	LOGPATH  string `envconfig:"LOGPATH" default:"./url_check.log"`
	APPTITLE string `default:"Go Url-Checker"`
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
var sites []Site     // i think this is declaring a slice, not an array
var results []Result // same as above

func main() {


	/* initialize settings */
	fmt.Printf("LOGPATH in main() before settings initialized, ```%v```\n", settings.LOGPATH)
	load_settings()
	fmt.Printf("LOGPATH in main() after settings initialized, ```%v```\n", settings.LOGPATH)

	/* initialize sites array */
	initialize_sites() // (https://stackoverflow.com/questions/26159416/init-array-of-structs-in-go)

	/* do the work */
	// check_sites(sites)
	// check_sites_just_with_routines(sites)
	check_sites_just_with_routines2(sites)

} // end func main()

/* ----------------------------------------------------------------------
   helper functions
   ---------------------------------------------------------------------- */

func load_settings() Settings {
	err := envconfig.Process("URL_CHECK", &settings)
	if err != nil {
		fmt.Printf("error, ```%v```", err.Error)
		panic(err)
	}
	fmt.Printf("LOGPATH in load_settings(), ```%v```\n", settings.LOGPATH)
	return Settings{}
}

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
	fmt.Println("\n sites...")
	spew.Dump(sites)
	return sites
}

/* ----------------------------------------------------------------------
   current experimentation
   ---------------------------------------------------------------------- */

func check_sites_just_with_routines2(sites []Site) {
	writer_channel := make(chan string)

	for _, site_element := range sites {
		// defer timeTrack(time.Now(), "check_sites_just_with_routines")
		go check_site2(site_element, writer_channel)
	}


	time.Sleep(5 * time.Second)

	x := <-writer_channel
	fmt.Println("value-x, ", x)

	y := <-writer_channel
	fmt.Println("value-y, ", y)

	close(writer_channel)

	// time.Sleep(100 * time.Millisecond)
}

func check_site2(site Site, writer_channel chan string) {
	start := time.Now()

	resp, _ := http.Get(site.url)
	body_bytes, _ := ioutil.ReadAll(resp.Body)
	text := string(body_bytes)
	var site_check_result string = "not_found"
	if strings.Contains(text, site.expected) {
		site_check_result = "found"
	}

	elapsed := time.Since(start)

	result_instance := Result{
		label:        site.label,
		check_result: site_check_result,
		time_taken:   elapsed,
	}
	fmt.Println("result_instance.label, ", result_instance.label)
	writer_channel <- result_instance.label
}

/* ---------------------------------------------------------------------- */

// func timeTrack(start time.Time, name string) {
//     elapsed := time.Since(start)
//     fmt.Printf("%s took %dms", name, elapsed.Nanoseconds()/1000)
// }

func check_sites_just_with_routines(sites []Site) {
	total_start := time.Now()
	for _, site_element := range sites {
		// defer timeTrack(time.Now(), "check_sites_just_with_routines")
		go check_site(site_element)
	}
	// time.Sleep(100 * time.Millisecond)
	var input string
	fmt.Scanln(&input)
	fmt.Println("done")
	total_elapsed := time.Since(total_start)
	fmt.Println("total_elapsed, ", total_elapsed)

	// fmt.Println("\n results...")
	// spew.Dump(results)

	fmt.Println("\n results before sort...")
	spew.Dump(results)

	sort.Slice(results, func(i, j int) bool { return results[i].label < results[j].label })
	fmt.Println("\n results after sorting by label...")
	spew.Dump(results)

	sort.Slice(results, func(i, j int) bool { return results[i].time_taken < results[j].time_taken })
	fmt.Println("\n results after sorting by time_taken...")
	spew.Dump(results)
}

func check_site(site Site) {
	start := time.Now()
	// fmt.Println( "start, ", start )
	// fmt.Println("\nsite -- ", site.label)
	// fmt.Println( "start for site %s, ```%v```", site.label, start  )
	// fmt.Printf("start for site %v, ```%v```\n", site.label, start)

	resp, _ := http.Get(site.url)
	// fmt.Println("status code -- ", resp.StatusCode)
	body_bytes, _ := ioutil.ReadAll(resp.Body)
	text := string(body_bytes)
	var site_check_result string = "not_found"
	if strings.Contains(text, site.expected) {
		site_check_result = "found"
	}
	// fmt.Println("site_check_result, ", site_check_result)

	elapsed := time.Since(start)
	// var elapsed time.Duration
	// elapsed = time.Since(start)

	// end := time.Now()
	// // fmt.Println( "end.String, ", end.String() )
	// fmt.Printf("end for site %v, ```%v```\n", site.label, end)

	// elapsed := end.Sub(start)
	// fmt.Println("elapsed, ", elapsed)

	// fmt.Println("elapsed has TypeOf: ", reflect.TypeOf(elapsed))
	// elapsed_k := reflect.ValueOf(elapsed)
	// fmt.Println("elapsed has Kind: ", elapsed_k.Kind())

	result_instance := Result{
		label:        site.label,
		check_result: site_check_result,
		time_taken:   elapsed,
	}
	fmt.Println("result_instance, ", result_instance)
	results = append(results, result_instance)

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
			time_taken:   elapsed,
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

	// fmt.Println("\n results before sort...")
	// spew.Dump(results)

	// sort.Slice(results, func(i, j int) bool { return results[i].time_taken < results[j].time_taken })

	// fmt.Println("\n results after sort...")
	// spew.Dump(results)

	// above may not handle non-ascii characters: <https://stackoverflow.com/a/38808838> -- update, it appears to handle non-ascii characters fine
}

// defer timeTrack(time.Now(), "lookup-and-check")

// func timeTrack(start time.Time, name string) {
//  elapsed := time.Since(start)
//  fmt.Printf("%s took %s\n", name, elapsed)
// }

// if strings.Contains(str, subStr) {}

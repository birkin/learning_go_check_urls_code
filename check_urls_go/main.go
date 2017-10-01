package main

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	// "reflect"
	"strings"
	"time"

	// "github.com/davecgh/go-spew/spew"      // easy way to pretty-print structs
	"github.com/kelseyhightower/envconfig" // binds env vars to settings struct
	"github.com/romana/rlog"
)

type Settings struct {
	DB_USERNAME string `envconfig:"DB_USERNAME" default:"default_db_username"`
	DB_PASSWORD string `envconfig:"DB_PASSWORD" default:"default_db_password"`
	DB_HOST     string `envconfig:"DB_HOST" default:"default_db_host"`
	DB_PORT     string `envconfig:"DB_PORT" default:"default_db_port"`
	DB_NAME     string `envconfig:"DB_NAME" default:"default_db_name"`
}

type Site struct {
	id                          int
	name                        string
	url                         string
	text_expected               string
	email_addresses             string
	email_message               string
	recent_checked_time         time.Time
	recent_checked_result       string
	previous_checked_result     string
	pre_previous_checked_result string
	next_check_time             time.Time
	custom_time_taken           time.Duration
}

var settings Settings
var sites []Site // i think this declares a slice, not an array
var db *sql.DB
var now_string string

func main() {
	/* Loads settings, initializes sites array, calls worker function. */

	rlog.Info("\n\nstarting")

	/// initialize settings
	rlog.Debug(fmt.Sprintf("settings before settings initialized, ```%#v```", settings))
	load_settings()

	/// access db
	db = setup_db()

	/// prepare current-time
	t := time.Now()
	now_string = fmt.Sprintf("%v", t.Format("2006-01-02 15:04:05"))
	rlog.Debug(fmt.Sprintf("now_string, ```%v```", now_string))
	// rlog.Debug(now_string)

	/// initialize sites
	initialize_sites_from_db()
	rlog.Debug("sites from db initialized")
	defer db.Close()

	/// call worker function
	check_sites_with_goroutines(sites)

} // end func main()

/* ----------------------------------------------------------------------
   helper functions
   ---------------------------------------------------------------------- */

func load_settings() Settings {
	/* Loads settings, eventually for logging and database.
	   Called by main() */
	err := envconfig.Process("url_check_", &settings) // env settings look like `URL_CHECK__THE_SETTING`
	if err != nil {
		raiseErr(err)
	}
	rlog.Debug(fmt.Sprintf("settings after settings initialized, ```%#v```", settings))
	return Settings{}
}

func setup_db() *sql.DB {
	/* Initializes db object and confirms connection.
	   Called by main() */
	var connect_str string = fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?parseTime=true",
		settings.DB_USERNAME, settings.DB_PASSWORD, settings.DB_HOST, settings.DB_PORT, settings.DB_NAME) // user:password@tcp(host:port)/dbname
	// rlog.Debug(fmt.Sprintf("connect_str, ```%v```", connect_str))
	db, err := sql.Open("mysql", connect_str)
	if err != nil {
		raiseErr(err)
	}
	// rlog.Debug(fmt.Sprintf("db after open, ```%v```", db))

	/// sql.Open doesn't open a connection, so validate DSN (data source name) data
	err = db.Ping()
	if err != nil {
		raiseErr(err)
	}
	// rlog.Debug(fmt.Sprintf("db after ping, ```%v```", db))
	// fmt.Println("db has TypeOf: ", reflect.TypeOf(db))
	// db_k := reflect.ValueOf(db)
	// fmt.Println("db has Kind: ", db_k.Kind())
	return db
}

func initialize_sites_from_db() []Site {
	/* Loads sites from db data
	   (https://stackoverflow.com/questions/26159416/init-array-of-structs-in-go)
	   Called by main() */
	sites = []Site{}
	querystring := fmt.Sprintf("SELECT `id`, `name`, `url`, `text_expected`, `email_addresses`, `email_message`, `previous_checked_result`, `pre_previous_checked_result`, `next_check_time` FROM `site_check_app_checksite`")
	// querystring := fmt.Sprintf("SELECT `id`, `name`, `url`, `text_expected`, `email_addresses`, `email_message`, `previous_checked_result`, `pre_previous_checked_result`, `next_check_time` FROM `site_check_app_checksite` WHERE `next_check_time` <= '%v' ORDER BY `next_check_time` ASC", now_string)
	rlog.Debug(fmt.Sprintf("querystring, ```%v```", querystring))
	rows, err := db.Query(querystring)
	if err != nil {
		raiseErr(err)
	}
	for rows.Next() {
		var id int
		var name string
		var url string
		var text_expected string
		var email_addresses string
		var email_message string
		var previous_checked_result string
		var pre_previous_checked_result string
		var next_check_time time.Time
		err = rows.Scan(&id, &name, &url, &text_expected, &email_addresses, &email_message, &previous_checked_result, &pre_previous_checked_result, &next_check_time)
		if err != nil {
			raiseErr(err)
		}
		sites = append(
			sites,
			Site{id, name, url, text_expected, email_addresses, email_message, time.Now(), "insert_check_result_here", previous_checked_result, pre_previous_checked_result, next_check_time, 0}, // name, url-to-check, text_expected, email_addresses, email_message, recent_checked_time, recent_checked_result, previous_checked_result, pre_previous_checked_result, next_check_time, custom_time_taken
		)
	}
	// rlog.Debug(fmt.Sprintf("rows, ```%v```", rows))
	db.Close()

	/// temp -- to just take a subset of the above during testing
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	site1 := sites[rand.Intn(len(sites))]
	site2 := sites[rand.Intn(len(sites))]
	sites = []Site{}
	sites = append(sites, site1, site2)
	/// end temp

	rlog.Info(fmt.Sprintf("sites to process, ```%#v```", sites)) // prints, eg, `{name:"clusters api", url:"etc...`
	return sites
}

//
func check_sites_with_goroutines(sites []Site) {
	/* Flow:
	   - creates channel,
	   - kicks off go-routines to run the web-checks,
	   - channel writes each check-result to db,
	   - closes channel.
	   Called by main() */

	/*
		TODO flow...
		- initialize db-writer channel, which will get a Site, not a Result
		- for each site:
			- start `check_site()` go-routine
		- have the channel write the results of each updated site to the db
		- for each updated site, start `check_email_need()` go-routine, which will:
			- get email_flag
			- if email_flag is `send_failure_email` or `send_success_email`:
				- send email
	*/

	rlog.Debug(fmt.Sprintf("starting check_sites"))
	main_start := time.Now()

	/// initialize channel
	dbwriter_channel := make(chan Site)

	/// start go routines
	for _, site_element := range sites {
		// rlog.Debug(fmt.Sprintf("here"))
		go check_site(site_element, dbwriter_channel)
	}

	// rlog.Info(fmt.Sprintf("len(dbwriter_channel), ```%v```", len(dbwriter_channel)))

	/// output channel data
	var counter int = 0
	var channel_output Site
	for channel_output = range dbwriter_channel {
		// rlog.Debug(fmt.Sprintf("counter, ```%v```", counter))
		// rlog.Debug(fmt.Sprintf("len(sites), ```%v```", len(sites)))
		counter++
		time.Sleep(50 * time.Millisecond)
		rlog.Info(fmt.Sprintf("channel-value, ```%#v```", channel_output))
		if counter == len(sites) {
			// rlog.Info("about to close channel")
			close(dbwriter_channel)
			rlog.Info("channel closed")
			break // shouldn't be needed
		}
	}
	main_elapsed := time.Since(main_start)
	rlog.Info(fmt.Sprintf("main_elapsed, ```%v```", main_elapsed))

}

func check_site(site Site, dbwriter_channel chan Site) {
	/* Checks site, stores data to updated-site, & writes updated-site to channel. */
	rlog.Debug(fmt.Sprintf("go routine started for site, ```%v```", site.name))

	/// check site
	var site_check_result string = "text_not_found"
	mini_start := time.Now()
	resp, err := http.Get(site.url)
	if err != nil {
		rlog.Info(fmt.Sprintf("error accessing site, `%v`; error, ```%v```", site.name, err))
		site_check_result = "site_not_reachable"
	} else {
		body_bytes, _ := ioutil.ReadAll(resp.Body)
		text := string(body_bytes)
		if strings.Contains(text, site.text_expected) {
			site_check_result = "passed"
		}
	}

	/// store result
	mini_elapsed := time.Since(mini_start)
	site.recent_checked_result = site_check_result
	site.custom_time_taken = mini_elapsed

	/// write info to channel
	dbwriter_channel <- site
	rlog.Info(fmt.Sprintf("site-info after write to channel, ```%#v```", site))

}

func raiseErr(err error) {
	/* logs and raises error */
	if err != nil {
		rlog.Error(fmt.Sprintf("err, ```%v```", err))
		panic(err)
	}
}

// func initialize_sites() []Site {
// 	/* Populates sites slice.
// 	   (https://stackoverflow.com/questions/26159416/init-array-of-structs-in-go) */
// 	sites = []Site{}
// 	sites = append(
// 		sites,
// 		Site{
// 			name:    "repo_file",
// 			url:      "https://repository.library.brown.edu/storage/bdr:6758/PDF/",
// 			text_expected: "BleedBox", // note: since brace is on following line, this comma is required
// 		},
// 		Site{"repo_search",
// 			"https://repository.library.brown.edu/studio/search/?q=elliptic",
// 			"The sequence of division polynomials"},
// 		Site{"bipg_wiki",
// 			"https://wiki.brown.edu/confluence/display/bipg/Brown+Internet+Programming+Group+Home",
// 			"The BIPG idea"},
// 		Site{"booklocator_app",
// 			"http://library.brown.edu/services/book_locator/?callnumber=GC97+.C46&location=sci&title=Chemistry+and+biochemistry+of+estuaries&status=AVAILABLE&oclc_number=05831908&public=true",
// 			"GC97 .C46 Level 11, Aisle 2A"},
// 		Site{"callnumber_app",
// 			"https://apps.library.brown.edu/callnumber/v2/?callnumber=PS3576",
// 			"American Literature"},
// 		Site{"clusters api",
// 			"https://library.brown.edu/clusters_api/data/",
// 			"scili-friedman"},
// 		Site{"easyborrow_feed",
// 			"http://library.brown.edu/easyborrow/feeds/latest_items/",
// 			"easyBorrow -- recent requests"},
// 		Site{"freecite",
// 			"http://freecite.library.brown.edu/welcome/",
// 			"About FreeCite"},
// 		Site{"iip_inscriptions",
// 			"http://library.brown.edu/cds/projects/iip/viewinscr/abur0001/",
// 			"Khirbet Abu Rish"},
// 		Site{"iip_processor",
// 			"https://apps.library.brown.edu/iip_processor/info/",
// 			"hi"},
// 		Site{"not_found_test",
// 			"https://apps.library.brown.edu/iip_processor/info/",
// 			"foo"},
// 	)

// 	/// temp -- to just take a subset of the above during testing
// 	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
// 	site1 := sites[rand.Intn(len(sites))]
// 	site2 := sites[rand.Intn(len(sites))]
// 	sites = []Site{}
// 	sites = append(sites, site1, site2)
// 	/// end temp

// 	rlog.Info(fmt.Sprintf("sites to process, ```%#v```", sites)) // prints, eg, `{name:"clusters api", url:"etc...`
// 	return sites
// }

/// EOF

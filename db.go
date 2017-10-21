package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // package is imported only for its `side-effects`; it gets registered as the driver for the regular database/sql package
	"github.com/romana/rlog"
	"math/rand"
	// "reflect"
	"time"
)

func setup_db(user string, pass string, host string, port string, name string) *sql.DB {
	/* Initializes db object and confirms connection.
	   Called by main() */
	var connect_str string = fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?parseTime=true",
		user, pass, host, port, name) // user:password@tcp(host:port)/dbname
	db, err := sql.Open("mysql", connect_str)
	if err != nil {
		msg := fmt.Sprintf("error connecting to db, ```%v```", err)
		rlog.Error(msg)
		panic(msg)
	}
	/// sql.Open doesn't open a connection, so validate DSN (data source name) data
	err = db.Ping()
	if err != nil {
		msg := fmt.Sprintf("error accessing db, ```%v```", err)
		rlog.Error(msg)
		panic(msg)
	}
	return db
} // end func setup_db()

func initialize_sites_from_db(db *sql.DB) []Site {
	/* Loads sites from db data
	   (https://stackoverflow.com/questions/26159416/init-array-of-structs-in-go)
	   Called by main() */
	sites = []Site{}
	querystring := fmt.Sprintf("SELECT `id`, `name`, `url`, `text_expected`, `email_addresses`, `email_message`, `recent_checked_result`, `previous_checked_result`, `pre_previous_checked_result`, `calculated_seconds`, `next_check_time` FROM `site_check_app_checksite`")
	// querystring := fmt.Sprintf("SELECT `id`, `name`, `url`, `text_expected`, `email_addresses`, `email_message`, `previous_checked_result`, `pre_previous_checked_result`, `next_check_time` FROM `site_check_app_checksite`")
	// querystring := fmt.Sprintf("SELECT `id`, `name`, `url`, `text_expected`, `email_addresses`, `email_message`, `previous_checked_result`, `pre_previous_checked_result`, `calculated_seconds`, `next_check_time` FROM `site_check_app_checksite` WHERE `next_check_time` <= '%v' ORDER BY `next_check_time` ASC", now_string)
	rlog.Debug(fmt.Sprintf("querystring, ```%v```", querystring))
	rows, err := db.Query(querystring)
	if err != nil {
		msg := fmt.Sprintf("error querying db, ```%v```", err)
		rlog.Error(msg)
		panic(msg)
	}
	for rows.Next() {
		var id int
		var name string
		var url string
		var text_expected string
		var email_addresses string
		var email_message string
		var recent_checked_result string
		var previous_checked_result string
		var pre_previous_checked_result string
		var calculated_seconds int
		var next_check_time time.Time
		err = rows.Scan(&id, &name, &url, &text_expected, &email_addresses, &email_message, &recent_checked_result, &previous_checked_result, &pre_previous_checked_result, &calculated_seconds, &next_check_time)
		if err != nil {
			msg := fmt.Sprintf("error scanning db rows, ```%v```", err)
			rlog.Error(msg)
			panic(msg)
		}
		sites = append(
			sites,
			Site{id, name, url, text_expected, email_addresses, email_message, time.Now(), recent_checked_result, previous_checked_result, pre_previous_checked_result, calculated_seconds, next_check_time, 0}, // name, url-to-check, text_expected, email_addresses, email_message, recent_checked_time, recent_checked_result, previous_checked_result, pre_previous_checked_result, next_check_time, custom_time_taken
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

} // end func initialize_sites_from_db()

func save_check_result(site Site) {
	/* 	Saves data to db.
	Called by: check.go/check_sites_with_goroutines() */
	rlog.Info(fmt.Sprintf("will save check-result to db for site, ```%#v```", site))

	settings := load_settings() // settings.go
	rlog.Debug(fmt.Sprintf("settings, ```%#v```", settings))
	db = setup_db(settings.DB_USERNAME, settings.DB_PASSWORD, settings.DB_HOST, settings.DB_PORT, settings.DB_NAME) // db.go
	defer db.Close()

	var sql_save_string string = fmt.Sprintf(
		"UPDATE `site_check_app_checksite` "+
			"SET `pre_previous_checked_result`='%v', `previous_checked_result`='%v', `next_check_time`='%v' "+
			"WHERE `id`=%v;",
		site.previous_checked_result, site.recent_checked_result, site.next_check_time, site.id)
	rlog.Debug(fmt.Sprintf("sql_save_string, ```%v```", sql_save_string))

	// var sql_save_string string = fmt.Sprintf(
	// 	"UPDATE `site_check_app_checksite` " +
	// 	"SET `pre_previous_checked_result`=site.previous_checked_result, `previous_checked_result`=site.recent_checked_result, `next_check_time`=site.next_check_time " +
	// 	"WHERE `id`=site.id;"
	// 	)

	result, err := db.Exec(sql_save_string)
	rlog.Debug(fmt.Sprintf("result, ```%v```", result))
	rlog.Debug(fmt.Sprintf("err, ```%v```", err))

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

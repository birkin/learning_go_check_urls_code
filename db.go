package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql" // package is imported only for its `side-effects`; it gets registered as the driver for the regular database/sql package
	"github.com/romana/rlog"
)

var db *sql.DB

func initialize_sites_from_db(DB_USERNAME string, DB_PASSWORD string, DB_HOST string, DB_PORT string, DB_NAME string, DB_TABLE string) []Site {
	/* Loads sites from db data
	   (https://stackoverflow.com/questions/26159416/init-array-of-structs-in-go)
	   Called by main()
	   TODO: once testing is complete, add the clause to the querystring: ```WHERE `next_check_time` <= '%v' ORDER BY `next_check_time` ASC", DB_TABLE, now_string)``` */
	db = setup_db(DB_USERNAME, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	defer db.Close()
	sites = []Site{}

	// querystring := fmt.Sprintf("SELECT `id`, `name`, `url`, `text_expected`, `email_addresses`, `email_message`, `recent_checked_result`, `previous_checked_result`, `pre_previous_checked_result`, `calculated_seconds`, `next_check_time` FROM `%v`", DB_TABLE)
	querystring := fmt.Sprintf("SELECT `id`, `name`, `url`, `text_expected`, `email_addresses`, `email_message`, `recent_checked_result`, `previous_checked_result`, `pre_previous_checked_result`, `calculated_seconds`, `next_check_time`, `check_frequency_number`, `check_frequency_unit` FROM `%v`", DB_TABLE)

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
		var check_frequency_number int
		var check_frequency_unit string
		err = rows.Scan(&id, &name, &url, &text_expected, &email_addresses, &email_message, &recent_checked_result, &previous_checked_result, &pre_previous_checked_result, &calculated_seconds, &next_check_time, &check_frequency_number, &check_frequency_unit)
		if err != nil {
			msg := fmt.Sprintf("error scanning db rows, ```%v```", err)
			rlog.Error(msg)
			panic(msg)
		}
		sites = append(
			sites,
			Site{id, name, url, text_expected, email_addresses, email_message, time.Now(), recent_checked_result, previous_checked_result, pre_previous_checked_result, calculated_seconds, next_check_time, check_frequency_number, check_frequency_unit, 0}, // name, url-to-check, text_expected, email_addresses, email_message, recent_checked_time, recent_checked_result, previous_checked_result, pre_previous_checked_result, next_check_time, check_frequency_number, custom_time_taken
		)

	}
	// rlog.Debug(fmt.Sprintf("rows, ```%v```", rows))

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

func setup_db(user string, pass string, host string, port string, name string) *sql.DB {
	/* Initializes db object and confirms connection.
	   Called by initialize_sites_from_db() and save_check_result() */
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

func save_check_result(site Site) {
	/*  Saves data to db.
	    Called by: check.go/check_sites_with_goroutines() */
	rlog.Info(fmt.Sprintf("will save check-result to db for site, ```%#v```", site))

	settings := load_settings() // settings.go
	rlog.Debug(fmt.Sprintf("settings, ```%#v```", settings))
	db = setup_db(settings.DB_USERNAME, settings.DB_PASSWORD, settings.DB_HOST, settings.DB_PORT, settings.DB_NAME) // db.go
	defer db.Close()

	var next_check_time_string string = site.next_check_time.Format("2006-01-02 15:04:05")
	var sql_save_string string = fmt.Sprintf(
		"UPDATE `%v` "+
			"SET `pre_previous_checked_result`='%v', `previous_checked_result`='%v', `recent_checked_result`='%v', `next_check_time`='%v' "+
			"WHERE `id`=%v;",
		settings.DB_TABLE, site.pre_previous_checked_result, site.previous_checked_result, site.recent_checked_result, next_check_time_string, site.id)
	rlog.Debug(fmt.Sprintf("sql_save_string, ```%v```", sql_save_string))
	result, err := db.Exec(sql_save_string)
	rlog.Debug(fmt.Sprintf("result, ```%v```", result))
	rlog.Debug(fmt.Sprintf("err, ```%v```", err))
}

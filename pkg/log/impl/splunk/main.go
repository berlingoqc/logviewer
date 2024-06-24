package splunk

import (
	"time"

	"github.com/berlingoqc/logviewer/pkg/http"
	"github.com/berlingoqc/logviewer/pkg/log/splunk/restapi"
	"github.com/berlingoqc/logviewer/pkg/ty"
)

func lol() {

	target := restapi.SplunkTarget{
		Endpoint: "https://splunk.nonprod.tktm.io/en-US/splunkd/__raw/servicesNS/william.quintal/search",
		Auth: http.CookieAuth{
			Cookie: `splunkweb_csrf_token_8000=13207931756135681617; token_key=13207931756135681617; experience_id=f59d1d7a-168d-c722-b996-4b3cbd113309; session_id_8000=6f7cd086145b9d7b14d28eb0f7874994663e5e1c; splunkd_8000=jKj32tFfKiGGQ96hjnx8KsGHjLbZcoBwzDPg1Aldv33Btcn6HH1C3aVdP3jxsEV0p5TM7IvAsUjP2iLT_kR^7IPskzFPJbnzvEprLSR9qNoLbNbH0qyYNlMK5mw7VGtl58u_Bfbcpo3RwpWmZIQDQ1NZ`,
		},
	}

	client, err := restapi.GetSplunkRestClient(target)
	if err != nil {
		panic(err)
	}

	headersCreate := ty.MS{
		"Host":              "splunk.nonprod.tktm.io",
		"User-Agent":        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:127.0) Gecko/20100101 Firefox/127.0",
		"Accept":            "text/javascript, text/html, application/xml, text/xml",
		"Accept-Language":   "en-US,en;q=0.5",
		"Accept-Encoding":   "gzip, deflate, br, zstd",
		"X-Splunk-Form-Key": "13207931756135681617",
		"Content-Type":      "application/x-www-form-urlencoded; charset=UTF-8",
		"X-Requested-With":  "XMLHttpRequest",
		"Origin":            "https://splunk.nonprod.tktm.io",
		"Sec-Fetch-Dest":    "empty",
		"Sec-Fetch-Mode":    "cors",
		"Sec-Fetch-Site":    "same-origin",
		"Sec-GPC":           "1",
	}

	dataCreate := ty.MS{
		"rf":                              "*",
		"auto_cancel":                     "62",
		"status_buckets":                  "300",
		"output_mode":                     "json",
		"custom.display.page.search.mode": "smart",
		"custom.dispatch.sample_ratio":    "1",
		"custom.display.page.search.tab":  "events",
		"custom.display.general.type":     "events",
		"ui_dispatch_app":                 "search",
		"preview":                         "1",
		"adhoc_search_level":              "smart",
		"sample_ratio":                    "1",
		"check_risky_command":             "false",
		"provenance":                      "UI%3ASearch",
	}

	job, err := client.CreateSearchJob(`index=prd3392+application_name="tmc.services.checkout"s`, "-24h@h", "now", headersCreate, dataCreate)

	if err != nil {
		panic(err)
	}

	sid := job.Sid

	isDone := false

	for !isDone {
		time.Sleep(1 * time.Second)
		status, err := client.GetSearchStatus(sid)

		if err != nil {
			panic(err)
		}

		isDone = status.Entry[0].Content.IsDone
	}

	results, err := client.GetSearchResult(sid, 0, 20)

	if err != nil {
		panic(err)
	}

	println(len(results.Results))

}

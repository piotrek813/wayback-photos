package wayback

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func handleErrorMsg(err error, msg string) {
	if err != nil {
		print(msg)
		handleError(err)
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}

}

func MockGetUrls(website string, mimetype []string, limit int, resumeKey string) (*WaybackResponse, bool) {
	return &WaybackResponse{
		Urls:      []WaybackUrl{*NewWaybackUrl("http://staff.edu.pl/Akademia Lokalna Cisco/album/thumbs/DSC02866.JPG", "20120823163741", "image/jpeg"), *NewWaybackUrl("http://staff.edu.pl/Akademia Lokalna Cisco/album/thumbs/DSC02866.JPG", "20120823163741", "image/jpeg"), *NewWaybackUrl("http://staff.edu.pl/Akademia Lokalna Cisco/album/thumbs/DSC02866.JPG", "20120823163741", "image/jpeg")},
		ResumeKey: "csdcd",
	}, true
}

func GetUrls(website string, mimetype []string, limit int, resumeKey string) (*WaybackResponse, bool) {
	u, _ := url.Parse("http://web.archive.org/cdx/search/cdx")

	q := u.Query()

	q.Set("output", "json")
	q.Set("matchType", "domain")
	q.Set("fl", "original,timestamp,mimetype")

	q.Set("url", website)
	q.Set("filter", "mimetype:("+strings.Join(mimetype, "|")+")")
	q.Set("limit", fmt.Sprint(limit))
	q.Set("showResumeKey", "true")

	if resumeKey != "" {
		q.Set("resumeKey", resumeKey)
	}

	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())

	if err != nil {
		log.Println(err)
		return nil, false
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, false
	}

	var res WaybackResponse
	err = json.Unmarshal(body, &res)
	handleError(err)

	return &res, true
}

type WaybackResponse struct {
	Urls      []WaybackUrl
	ResumeKey string
}

func (w *WaybackResponse) UnmarshalJSON(data []byte) error {
	var s [][]string

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	for i, row := range s {
		// first row contains names of columns
		if i == 0 {
			continue
		}

		// second to last row is empty
		if len(row) == 0 {
			continue
		}

		if i == len(s)-1 {
			w.ResumeKey = row[0]
			continue
		}
		url, err := WaybackUrlFromRow(row)

		if err != nil {
			continue
		}

		w.Urls = append(w.Urls, *url)
	}

	return nil
}

type WaybackUrl struct {
	Timestamp string
	Original  string
	Mimetype  string
	Direct    string
}

func NewWaybackUrl(original string, timestamp string, mimetype string) *WaybackUrl {
	var w WaybackUrl
	w.Original = original
	w.Timestamp = timestamp
	w.Mimetype = mimetype

	d := []string{"http://web.archive.org/web", timestamp + "if_", original}
	w.Direct = strings.Join(d, "/")

	return &w
}

func WaybackUrlFromRow(row []string) (*WaybackUrl, error) {
	if len(row) >= 3 {
		return NewWaybackUrl(row[0], row[1], row[2]), nil
	}

	msg := fmt.Sprintf("row: %v couldn't be transformed into WaybackUrl, length is to short", row)
	return nil, errors.New(msg)

}

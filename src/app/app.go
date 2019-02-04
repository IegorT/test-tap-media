package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
	maxminddb "github.com/oschwald/maxminddb-golang"
	"github.com/ua-parser/uap-go/uaparser"
)

// The Output struct
type Output struct {
	ID             string `json:"request_id"`
	OS             string `json:"os"`
	Device         string `json:"device"`
	Browser        string `json:"browser"`
	Domain         string `json:"domain"`
	CountryISOCode string `json:"country_code"`
}

// The Location struct
type Location struct {
	Country struct {
		ISOCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
}

// The BidRequest struct
// TODO: validation
type BidRequest struct {
	ID     string `validate:"required"`
	Device struct {
		UA string
		IP string
	}
	Site struct {
		Page string
	}
	App struct {
		Domain string
	}

	UAParser   *uaparser.Parser
	LocationDB *maxminddb.Reader
}

// URLParse func
func (br *BidRequest) URLParse() (*url.URL, error) {

	u, err := url.Parse(br.Site.Page)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// BodyParse func
func (br *BidRequest) BodyParse(i io.Reader) error {
	if err := json.NewDecoder(i).Decode(br); err != nil {
		return err
	}
	return nil
}

// HTTPHandle func
func (br BidRequest) HTTPHandle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	br.BodyParse(r.Body)
	out, err := br.NewOutput()
	if err != nil {
		fmt.Println("error in request: ", &out.ID)
	}

	json, err := json.Marshal(out)
	if err != nil {
		fmt.Println("error parse to json: ", &out)
	}

	fmt.Println(string(json))
	w.Write([]byte("{\"status\": \"Ok\"}"))
}

// UserAgentParse func
func (br *BidRequest) UserAgentParse() *uaparser.Client {

	return br.UAParser.Parse(br.Device.UA)
}

// IPLocation func
func (br *BidRequest) IPLocation() (*Location, error) {
	var loc Location

	ip := net.ParseIP(br.Device.IP)
	if ip == nil {
		return &loc, nil
	}

	err := br.LocationDB.Lookup(ip, &loc)
	if err != nil {
		return &loc, err
	}

	return &loc, nil
}

// NewOutput func
// TODO: error handling
func (br *BidRequest) NewOutput() (*Output, error) {
	o := &Output{
		ID: br.ID,
	}

	if br.App.Domain != "" {
		o.Domain = br.App.Domain
	} else {
		url, err := br.URLParse()
		if err == nil {
			o.Domain = url.Hostname()
		}
	}

	ua := br.UserAgentParse()
	o.OS = ua.Os.ToString()
	o.Device = ua.Device.ToString()
	o.Browser = ua.UserAgent.ToString()

	loc, err := br.IPLocation()
	if err == nil {
		o.CountryISOCode = loc.Country.ISOCode
	}

	return o, nil
}

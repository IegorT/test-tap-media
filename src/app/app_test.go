package app

import (
	"log"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	maxminddb "github.com/oschwald/maxminddb-golang"
	"github.com/ua-parser/uap-go/uaparser"
)

func TestApp(t *testing.T) {

	db, err := maxminddb.Open("../../db/GeoLite2-Country.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	parser, err := uaparser.New("../../db/regexes.yaml")
	if err != nil {
		log.Fatal(err)
	}

	tables := []struct {
		device string
		output *Output
	}{
		{
			"android",
			&Output{
				ID:             "FxlYjsHZs1",
				OS:             "Android 6.0",
				Device:         "Huawei Nexus 6P",
				Browser:        "Chrome Mobile 45.0.2454",
				CountryISOCode: "AR",
				Domain:         "play.google.com",
			},
		},
		{
			"app",
			&Output{
				ID:             "c1ba8e13-62da-46e3-884e-376f901e28f9",
				OS:             "Mac OS X 10.10.5",
				Device:         "Other",
				Browser:        "Chrome 52.0.2743",
				CountryISOCode: "US",
				Domain:         "yourapp.com",
			},
		},
		{
			"ios",
			&Output{
				ID:             "FxU0032U8a",
				OS:             "iOS 6.1.4",
				Device:         "iPhone",
				Browser:        "Mobile Safari 6.0",
				CountryISOCode: "AR",
				Domain:         "play.google.com",
			},
		},
		{
			"web",
			&Output{
				ID:             "fb7b0979-2d9f-47d0-a854-3d3a275b471e",
				OS:             "Windows 10",
				Device:         "Other",
				Browser:        "Chrome 56.0.2924",
				Domain:         "www.yoursite.com",
				CountryISOCode: "US",
			},
		},
	}

	for _, table := range tables {
		br := &BidRequest{
			LocationDB: db,
			UAParser:   parser,
		}
		filePath := "../../mock/request-" + table.device + ".json"
		file, err := func(s string) (*os.File, error) {
			file, err := os.Open(s)
			if err != nil {
				return nil, err
			}

			return file, nil
		}(filePath)
		if err != nil {
			panic(err)
		}

		br.BodyParse(file)
		output, _ := br.NewOutput()
		if !cmp.Equal(output, table.output) {
			t.Errorf("Output of bid request was incorrect, got: %v, want: %v.", output, table.output)
		}
	}
}

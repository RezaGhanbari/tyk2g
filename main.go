// Tyk to Graylog
// Author: Reza Ghanbari
//

package main

import (
	"labix.org/v2/mgo"
	b64 "encoding/base64"
	"fmt"
	"encoding/json"
	"labix.org/v2/mgo/bson"
	"time"
	"strings"
	"github.com/robertkowalski/graylog-golang"
	"strconv"
	"os"
	//"net/http"
	//"github.com/gorilla/mux"
)

var (
	mongodb_host = os.Getenv("TYK_MONGO_HOST")
	mongodb_db = os.Getenv("TYK_MONGO_DB")
	mongodb_collection = os.Getenv("TYK_MONGO_COL")
	mongodb_port = os.Getenv("TYK_MONGO_PORT")
	tyk_time = os.Getenv("TYK_TIME")
	graylog_host = os.Getenv("TYK_GRAYLOG_HOST")
	graylog_port = os.Getenv("TYK_GRAYLOG_PORT")
)
type Response struct {
	ID            bson.ObjectId `bson:"_id,omitempty"`
	Method        string        `bson:"method"`
	Path          string        `bson:"path"`
	Rawpath       string        `bson:"rawpath"`
	Contentlength string        `bson:"contentlength"`
	Useragent     string        `bson:"useragent"`
	Day           string        `bson:"day"`
	Month         int16         `bson:"month"`
	Year          int32         `bson:"year"`
	Hour          int16         `bson:"hour"`
	Responsecode  int32         `bson:"responsecode"`
	Apikey        string        `bson:"apikey"`
	Timestamp     time.Time     `bson:"timestamp"`
	Apiversion    string        `bson:"apiversion"`
	Apiname       string        `bson:"apiname"`
	Apiid         string        `bson:"apiid"`
	Orgid         int64         `bson:"orgid"`
	Oauthid       string        `bson:"oauthid"`
	Requesttime   string        `bson:"requesttime"`
	Ipaddress     string        `bson:"ipaddress"`
	Geo           string        `bson:"geo"`
	Tags          string        `bson:"tags"`
	Alias         string        `bson:"alias"`
	Trackpath     string        `bson:"trackpath"`
	ExpireAt      string        `bson:"expireAt"`
	Rawrequest    string        `bson:"rawrequest"`
	Rawresponse   string        `bson:"rawresponse"`
}

type Message struct {
	TYK_MongoID       bson.ObjectId `bson:"_id,omitempty"`
	TYK_Method        string        `bson:"method"`
	TYK_Path          string        `bson:"path"`
	TYK_Rawpath       string        `bson:"rawpath"`
	TYK_Contentlength string        `bson:"contentlength"`
	TYK_Useragent     string        `bson:"useragent"`
	TYK_Day           string        `bson:"day"`
	TYK_Month         int16         `bson:"month"`
	TYK_Year          int32         `bson:"year"`
	TYK_Hour          int16         `bson:"hour"`
	TYK_Responsecode  int32         `bson:"responsecode"`
	TYK_Apikey        string        `bson:"apikey"`
	TYK_Timestamp     time.Time     `bson:"timestamp"`
	TYK_Apiversion    string        `bson:"apiversion"`
	TYK_Apiname       string        `bson:"apiname"`
	TYK_Apiid         string        `bson:"apiid"`
	TYK_Orgid         int64         `bson:"orgid"`
	TYK_Oauthid       string        `bson:"oauthid"`
	TYK_Requesttime   string        `bson:"requesttime"`
	TYK_Ipaddress     string        `bson:"ipaddress"`
	TYK_Geo           string        `bson:"geo"`
	TYK_Tags          string        `bson:"tags"`
	TYK_Alias         string        `bson:"alias"`
	TYK_Trackpath     string        `bson:"trackpath"`
	TYK_ExpireAt      string        `bson:"expireAt"`
	TYK_Rawrequest    string        `bson:"rawrequest"`
	TYK_Rawresponse   string        `bson:"rawresponse"`
}

// function to run every x milliseconds.
func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

// TYK_LOGGER function
func tykLogger(t time.Time) {
	session, err := mgo.Dial(mongodb_db+":"+mongodb_port)

	if err != nil {
		panic(err)
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB(mongodb_db).C(mongodb_collection)

	// Set time
	//fromDate := time.Date(2018, time.April, 16, 0, 0, 0, 0, time.UTC)
	j, err := strconv.ParseInt(tyk_time, 10, 32)
	if err != nil {
		panic(err)
	}
	fromDate := time.Now().Add(time.Duration(j) * time.Millisecond * -1)
	//toDate := time.Date(2018, time.April, 16, 24, 0, 0, 0, time.UTC)
	toDate := time.Now().UTC()

	var results []Response

	err = c.Find(
		bson.M{
			"timestamp": bson.M{
				"$gt": fromDate,
				"$lt": toDate,
			},
		}).Select(bson.M{
		"_id":           1,
		"method":        1,
		"path":          1,
		"rawpath":       1,
		"contentlength": 1,
		"useragent":     1,
		"day":           1,
		"month":         1,
		"year":          1,
		"hour":          1,
		"responsecode":  1,
		"apikey":        1,
		"timestamp":     1,
		"apiversion":    1,
		"apiname":       1,
		"apiid":         1,
		"orgid":         1,
		"oauthid":       1,
		"requesttime":   1,
		"ipaddress":     1,
		"geo":           1,
		"tags":          1,
		"alias":         1,
		"trackpath":     1,
		"expireAt":      1,
		"rawrequest":    1,
		"rawresponse":   1,
	}).All(&results)

	for _, v := range results {
		reqDec, _ := b64.StdEncoding.DecodeString(v.Rawrequest)
		resDec, _ := b64.StdEncoding.DecodeString(v.Rawresponse)

		var replacer = strings.NewReplacer("\n", " ", "   ", "", "\r", "")
		var replacer2 = strings.NewReplacer(",", " ", "\\", "", "\"", "", " n ", "")
		res1D := &Message{
			TYK_MongoID:       v.ID,
			TYK_Method:        v.Method,
			TYK_Path:          v.Path,
			TYK_Rawpath:       v.Rawpath,
			TYK_Contentlength: v.Contentlength,
			TYK_Useragent:     v.Useragent,
			TYK_Day:           v.Day,
			TYK_Month:         v.Month,
			TYK_Year:          v.Year,
			TYK_Hour:          v.Hour,
			TYK_Responsecode:  v.Responsecode,
			TYK_Apikey:        v.Apikey,
			TYK_Timestamp:     v.Timestamp,
			TYK_Apiversion:    v.Apiversion,
			TYK_Apiname:       v.Apiname,
			TYK_Apiid:         v.Apiid,
			TYK_Orgid:         v.Orgid,
			TYK_Oauthid:       v.Oauthid,
			TYK_Requesttime:   v.Requesttime,
			TYK_Ipaddress:     v.Ipaddress,
			TYK_Geo:           v.Geo,
			TYK_Tags:          v.Tags,
			TYK_Alias:         v.Alias,
			TYK_Trackpath:     v.Trackpath,
			TYK_ExpireAt:      v.ExpireAt,
			TYK_Rawrequest:    strings.Replace(replacer.Replace(string(reqDec)), " ", "", -1),
			TYK_Rawresponse:   strings.Replace(replacer.Replace(string(resDec)), " ", "", -1),
		}

		port , err := strconv.ParseInt(graylog_port, 10, 64)
		if err != nil {
			panic(err)
		}
		g := gelf.New(gelf.Config{
			GraylogPort:     int(port),
			GraylogHostname: graylog_host,
			Connection:      "udp",
			MaxChunkSizeWan: 42,
			MaxChunkSizeLan: 1337,
		})
		res1B, _ := json.Marshal(res1D)
		f := replacer2.Replace(string(res1B))
		//f := strings.Replace(v, " ","", -1)

		x := fmt.Sprintf(`{
		"version": "1.0",
		"host": "%s",
		"short_message": "%s"}`, string(graylog_host+":"+graylog_port) , string(f))

		g.Log(x)

		//fmt.Println(string(f))
		//fmt.Println(string(x))
	}
}

func main() {
	i, err := strconv.ParseInt(tyk_time, 10, 32)
	if err != nil {
		panic(err)
	}
	doEvery(time.Duration(i)*time.Millisecond, tykLogger)
}

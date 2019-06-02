package main

import (
	"fmt"
	"github.com/aerospike/aerospike-client-go"
	"github.com/caarlos0/env"
	"github.com/goji/httpauth"
	"html/template"
	"log"
	"net/http"
	"os"
)

var (
	cfg      config
	asClient *aerospike.Client
	tmplCamp *template.Template
	tmplDsp  *template.Template
)

type AdvCabinetRow struct {
	ID          string
	UserID      int
	Icon        string
	Image       string
	ClickAction string
	Country     string
	Price       float64
}

type PageAerospikeData struct {
	PageTitle      string
	AdvCabinetRows []AdvCabinetRow
}

type DspProxyTrackerRow struct {
	PageTitle   string
	PushID      string
	SSP         string
	DSP         string
	IconURL     string
	ClickAction string
	Country     string
	Node        string
	PriceSell   string
	PriceBuy    string
}

type config struct {
	AerospikeHosts []string `env:"AEROSPIKE_HOSTS" envDefault:"127.0.0.1" envSeparator:","`
	AerospikePort  int      `env:"AEROSPIKE_PORT" envDefault:"3000"`
	AerospikeNS    string   `env:"AEROSPIKE_NS" envDefault:"test"`
	AerospikeSet   string   `env:"AEROSPIKE_SET" envDefault:"campaigns"`
	Debug          bool     `env:"DEBUG" envDefault:"false"`
}

func init() {
	initEnv()
	initAerospike()
}

func initEnv() {
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal("failed to parse env: " + err.Error())
	}
	log.Printf("config: %#v", cfg)
}

func initAerospike() {
	var err error
	hosts := []*aerospike.Host{}

	for _, host := range cfg.AerospikeHosts {
		hosts = append(hosts, aerospike.NewHost(host, cfg.AerospikePort))
	}

	asClient, err = aerospike.NewClientWithPolicyAndHost(nil, hosts...)
	if err != nil {
		log.Println("failed to create aerospike client: " + err.Error())
		os.Exit(1)
	}

	log.Println("aerospike connection established")
}

func CampaignsHandler(w http.ResponseWriter, r *http.Request) {
	data := PageAerospikeData{
		PageTitle:      "AerospikeSet: " + cfg.AerospikeNS,
		AdvCabinetRows: []AdvCabinetRow{},
	}

	recs, _ := asClient.ScanAll(nil, cfg.AerospikeNS, cfg.AerospikeSet)
	for rec := range recs.Records {
		data.AdvCabinetRows = append(data.AdvCabinetRows, AdvCabinetRow{
			//ID:          rec.Key.Value().String(),
			UserID:      rec.Bins["user-id"].(int),
			Icon:        rec.Bins["icon"].(string),
			Image:       rec.Bins["image"].(string),
			ClickAction: rec.Bins["click-action"].(string),
			Country:     rec.Bins["country"].(string),
			Price:       rec.Bins["price"].(float64),
		})
	}

	tmplCamp.Execute(w, data)
}

func PushidHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	}

	keyAs := fmt.Sprintf("push:%s", keys[0])
	log.Println("keyAs: " + keyAs)

	key, err := aerospike.NewKey(cfg.AerospikeNS, "dsp_proxy_tracker", keyAs)
	if err != nil {
		log.Println("Create AS key error" + err.Error())
		w.Write([]byte("Create AS key error"))
		return
	}

	rec, err := asClient.Get(nil, key)
	if err != nil {
		log.Println("Get rec error: " + err.Error())
		w.Write([]byte("Key not Found"))
		return
	}

	data := DspProxyTrackerRow{
		PageTitle:   "DSP Proxy Tracker",
		PushID:      keys[0],
		SSP:         rec.Bins["ssp"].(string),
		DSP:         rec.Bins["dsp"].(string),
		IconURL:     rec.Bins["icon-url"].(string),
		ClickAction: rec.Bins["click-action"].(string),
		Country:     rec.Bins["country"].(string),
		Node:        rec.Bins["node"].(string),
		PriceSell:   rec.Bins["price-sell"].(string),
		PriceBuy:    rec.Bins["price-buy"].(string),
	}

	//
	//bin["ssp"] = ssp_host[0]
	//bin["dsp"] = win.Name
	//bin["icon-url"] = win.Img
	//bin["click-action"] = win.ClickUrl
	//bin["country"] = win.Country
	//bin["node"] = node
	//// bin["safe-uid"] = safeUID
	//bin["price-sell"] = strconv.FormatFloat(win.Price, 'f', 6, 64)
	//bin["price-buy"] = strconv.FormatFloat(win.Price, 'f', 6, 64)
	//
	err = tmplDsp.Execute(w, data)
	if err != nil {
		log.Println("Template generating error: " + err.Error())
		w.Write([]byte("Server err"))
	}
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	tmplCamp = template.Must(template.ParseFiles("templates/campaigns.html"))
	tmplDsp = template.Must(template.ParseFiles("templates/dspproxy.html"))

	http.Handle("/campaigns", httpauth.SimpleBasicAuth("hello", "6N5mKy4iz0dydOAGxJ7gWptRmo6JyIG6w89Wwe4H")(http.HandlerFunc(CampaignsHandler)))
	http.Handle("/pushid", httpauth.SimpleBasicAuth("hello", "6N5mKy4iz0dydOAGxJ7gWptRmo6JyIG6w89Wwe4H")(http.HandlerFunc(PushidHandler)))

	http.ListenAndServe(":8000", nil)
}

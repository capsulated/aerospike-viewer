package main

import (
	"github.com/aerospike/aerospike-client-go"
	"github.com/caarlos0/env"
	"html/template"
	"log"
	"net/http"
)

var (
	cfg      config
	asClient *aerospike.Client
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
	}

	log.Println("aerospike connection established")
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	tmpl := template.Must(template.ParseFiles("templates/tables.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

		tmpl.Execute(w, data)
	})

	http.ListenAndServe(":8000", nil)
}

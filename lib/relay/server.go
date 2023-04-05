package relay

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/fiatjaf/relayer"
	"github.com/fiatjaf/relayer/storage/postgresql"
	"github.com/kelseyhightower/envconfig"
	"github.com/nbd-wtf/go-nostr"
)

type Relay struct {
	PostgresDatabase string `envconfig:"POSTGRESQL_DATABASE"`

	storage *postgresql.PostgresBackend
}

func (r *Relay) Name() string {
	return "rpctest"
}

func (r *Relay) Storage() relayer.Storage {
	return r.storage
}

func (r *Relay) OnInitialized(*relayer.Server) {
	fmt.Println("Initialized successfully")
}

func (r *Relay) Init() error {
	err := envconfig.Process("", r)
	if err != nil {
		return fmt.Errorf("couldn't process envconfig: %w", err)
	}

	// every hour, delete all very old events
	go func() {
		db := r.Storage().(*postgresql.PostgresBackend)

		for {
			time.Sleep(60 * time.Minute)
			db.DB.Exec(`DELETE FROM event WHERE created_at < $1`, time.Now().AddDate(0, -3, 0).Unix()) // 3 months
		}
	}()

	return nil
}

func (r *Relay) AcceptEvent(evt *nostr.Event) bool {
	// block events that are too large
	jsonb, _ := json.Marshal(evt)
	if len(jsonb) > 10000 {
		return false
	}

	return true
}

func (r *Relay) BeforeSave(evt *nostr.Event) {
	// do nothing
}

func (r *Relay) AfterSave(evt *nostr.Event) {
	// delete all but the 100 most recent ones for each key
	r.Storage().(*postgresql.PostgresBackend).DB.Exec(`DELETE FROM event WHERE pubkey = $1 AND kind = $2 AND created_at < (
      SELECT created_at FROM event WHERE pubkey = $1
      ORDER BY created_at DESC OFFSET 100 LIMIT 1
    )`, evt.PubKey, evt.Kind)
}

func Server() {
	/*
		connCfg := &rpcclient.ConnConfig{
			Host:         "localhost:8332",
			User:         "adam",
			Pass:         "adam260699",
			HTTPPostMode: true,
			DisableTLS:   true,
		}


			client, err := rpcclient.New(connCfg, nil)
			if err != nil {
				fmt.Println("Starting without bitcoin core connection")
			} else {
				defer client.Shutdown()

				// Get the current block count.
				blockCount, err := client.GetDifficulty()
				if err != nil {
					log.Fatal(err)
				}
				log.Printf("Block difficulty: %e", blockCount)
			}
	*/

	r := Relay{}
	if err := envconfig.Process("", &r); err != nil {
		log.Fatalf("failed to read from env: %v", err)
		return
	}

	r.storage = &postgresql.PostgresBackend{DatabaseURL: r.PostgresDatabase} // "postgres://nostr:nostr@localhost:5432/nostr?sslmode=disable"
	if err := relayer.Start(&r); err != nil {
		log.Fatalf("server terminated: %v", err)
	}
}

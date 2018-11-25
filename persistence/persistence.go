/*
    ToDD Persistence Layer

	Copyright 2016 Matt Oswalt. Use or modification of this
	source code is governed by the license provided here:
	https://github.com/toddproject/todd/blob/master/LICENSE
*/

package persistence

import (
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/toddproject/todd/agent/defs"
	pb "github.com/toddproject/todd/api/exp/generated"
	"github.com/toddproject/todd/config"
	"github.com/toddproject/todd/server/objects"

	"github.com/dgraph-io/badger"
)

func NewPersistence(cfg config.Config) (*Persistence, error) {

	var p Persistence

	opts := badger.DefaultOptions
	opts.Dir = "/Users/mierdin/Code/GO/src/github.com/toddproject/todd/tmpdb"
	opts.ValueDir = "/Users/mierdin/Code/GO/src/github.com/toddproject/todd/tmpdb"
	p.db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}

	return &p, nil
}

type persistence struct {
	config config.Config
	db     badger.Db
}

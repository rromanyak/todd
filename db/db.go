/*
    ToDD Database Functions

    This file holds the infrastructure for database abstractions in ToDD.

	Copyright 2016 Matt Oswalt. Use or modification of this
	source code is governed by the license provided here:
	https://github.com/toddproject/todd/blob/master/LICENSE
*/

package db

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

var (
	ErrInvalidDBPlugin = errors.New("Invalid DB plugin in config file")
	ErrNotExist        = errors.New("Value does not exist")
)

// DatabasePackage represents all of the behavior that a ToDD database plugin must support
type DatabasePackage interface {

	// (no args)
	Init() error

	// (agent advertisement to set)
	SetAgent(defs.AgentAdvert) error

	GetAgent(string) (*defs.AgentAdvert, error)
	GetAgents() ([]defs.AgentAdvert, error)

	// (agent advertisement to remove)
	RemoveAgent(defs.AgentAdvert) error

	SetObject(objects.ToddObject) error
	GetObjects(string) ([]objects.ToddObject, error)
	DeleteObject(string, string) error

	GetGroupMap() (map[string]string, error)
	SetGroupMap(map[string]string) error

	// Testing
	InitTestRun(string, map[string]map[string]string) error
	SetAgentTestStatus(string, string, string) error
	GetTestStatus(string) (map[string]string, error)
	SetAgentTestData(string, string, string) error
	GetAgentTestData(string, string) (map[string]string, error)
	WriteCleanTestData(string, string) error
	GetCleanTestData(string) (string, error)

	// New protobufs stuff
	GetGroups() ([]*pb.Group, error)
	CreateGroup(*pb.Group) error
	DeleteGroup(*pb.Group) error
}

// NewToddDB will create a new instance of toddDatabase, and load the desired
// databasePackage-compatible comms package into it.
func NewToddDB(cfg config.Config) (DatabasePackage, error) {

	// // Create toddDatabase instance
	// var tdb DatabasePackage

	// // Load the appropriate DB package based on config file
	// switch cfg.DB.Plugin {
	// case "etcd":
	// 	tdb = newEtcdDB(cfg)
	// default:
	// 	return nil, ErrInvalidDBPlugin
	// }

	return nil, nil
}

type badgerDB struct {
	config config.Config
	// dbLock  mutex
}

func CreateGroup(group *pb.Group) error {
	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	opts := badger.DefaultOptions
	opts.Dir = "/Users/mierdin/Code/GO/src/github.com/toddproject/todd/tmpdb"
	opts.ValueDir = "/Users/mierdin/Code/GO/src/github.com/toddproject/todd/tmpdb"
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// err = db.Update(func(txn *badger.Txn) error {

	// 	groupJson, err := json.Marshal(group)
	// 	if err != nil {
	// 		log.Warn("Error converting group to json")
	// 	}

	// 	err = txn.Set([]byte(fmt.Sprintf("group/%s", group.Name)), groupJson)
	// 	if err != nil {
	// 		log.Warn("Unable to set group in DB")
	// 	}

	// 	return nil
	// })

	// if err != nil {
	// 	log.Warn("Unable to set group in DB")
	// }

	// Start a writable transaction.
	txn := db.NewTransaction(true)
	defer txn.Discard()

	groupJson, err := json.Marshal(group)
	if err != nil {
		log.Warn("Error converting group to json")
	}

	// Use the transaction...
	err = txn.Set([]byte(fmt.Sprintf("group/%s", group.Name)), groupJson)
	if err != nil {
		return err
	}

	// Commit the transaction and check for error.
	if err := txn.Commit(nil); err != nil {
		return err
	}

	return nil
}

func ListGroups() ([]*pb.Group, error) {
	// https://github.com/dgraph-io/badger/issues/436
	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	opts := badger.DefaultOptions
	opts.Dir = "/Users/mierdin/Code/GO/src/github.com/toddproject/todd/tmpdb"
	opts.ValueDir = "/Users/mierdin/Code/GO/src/github.com/toddproject/todd/tmpdb"
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	groups := []*pb.Group{}

	err = db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte("group/")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			v, _ := item.ValueCopy(nil)

			var group pb.Group
			err = json.Unmarshal(v, &group)
			if err != nil {
				log.Warn("Error converting group to json")
			}

			groups = append(groups, &group)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return groups, nil

}

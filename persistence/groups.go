package persistence

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	pb "github.com/toddproject/todd/api/exp/generated"

	"github.com/dgraph-io/badger"
)

func (p *persistence) CreateGroup(group *pb.Group) error {

	txn := db.NewTransaction(true)
	defer txn.Discard()

	groupJson, err := json.Marshal(group)
	if err != nil {
		log.Warn("Error converting group to json")
	}

	err = txn.Set([]byte(fmt.Sprintf("group/%s", group.Name)), groupJson)
	if err != nil {
		return err
	}

	if err := txn.Commit(nil); err != nil {
		return err
	}

	return nil
}

func (p *persistence) ListGroups() ([]*pb.Group, error) {
	// https://github.com/dgraph-io/badger/issues/436
	// Open the Badger database located in the referenced directory.
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

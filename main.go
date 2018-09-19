package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/boltdb/bolt"
)

const dbFile = "sample.db"
const bucket = "todos"

func main() {
	fmt.Println("main start...")

	// ids := []string{"1", "2", "3", "4", "5"}

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	batch(db)
	reference(db)

	fmt.Println("main end...")
}

func batch(db *bolt.DB) {
	var wg sync.WaitGroup
	ids := []string{"1", "2", "3", "4", "5"}

	for _, id := range ids {
		wg.Add(1)
		go func(itemID string) {
			fmt.Printf("goroutine called. itemId = %s\n", itemID)
			err := db.Batch(func(tx *bolt.Tx) error {
				b, err := tx.CreateBucketIfNotExists([]byte(bucket))
				if err != nil {
					return err
				}
				// comment to raise error
				// if itemID == "3" {
				// 	return errors.New(fmt.Sprintf("Error! itemId = %s", itemID))
				// }
				return b.Put([]byte(itemID), []byte(fmt.Sprintf("todo %s", itemID)))
			})
			if err != nil {
				fmt.Printf("Batch() error : %v\n", err)
			}
			wg.Done()
		}(id)
	}

	wg.Wait()
}

func reference(db *bolt.DB) error {
	return db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key:%v, value:%s\n", k, string(v))
		}
		return nil
	})
}

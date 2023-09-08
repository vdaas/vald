package bbolt

import (
	bolt "go.etcd.io/bbolt"
)

type Bbolt struct {
	db     *bolt.DB
	bucket string
}

func New(path string) *Bbolt {
	// TODO: 初期化をここでするか、DIするか。ライフタイムを管理するのだるいのでDIの方がいいかも
	return &Bbolt{}
}

func (b *Bbolt) Set(key string, val []byte) error {
	if err := b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(val)
		err := b.Put([]byte(key), nil)
		return err
	}); err != nil {
		return err
	}

	return nil
}

func (b *Bbolt) Get(key string) ([]byte, bool, error) {
	var val []byte
	if err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		copy(val, b.Get([]byte(key)))
		return nil
	}); err != nil {
		return nil, false, err
	}

	if val == nil {
		return nil, false, nil
	}

	return val, true, nil
}

// func main() {
// 	f, err := os.Create("my.db")
// 	f.Close()
// 	defer os.Remove(f.Name())

// 	// Open the my.db data file in your current directory.
// 	// It will be created if it doesn't exist.
// 	db, err := bolt.Open("my.db", 0600, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	db.Update(func(tx *bolt.Tx) error {
// 		_, err := tx.CreateBucket([]byte("MyBucket"))
// 		if err != nil {
// 			return fmt.Errorf("create bucket: %s", err)
// 		}
// 		return nil
// 	})

// 	db.Update(func(tx *bolt.Tx) error {
// 		b := tx.Bucket([]byte("MyBucket"))
// 		err := b.Put([]byte("answer"), nil)
// 		return err
// 	})

// 	db.View(func(tx *bolt.Tx) error {
// 		b := tx.Bucket([]byte("MyBucket"))
// 		v := b.Get([]byte("answer"))
// 		if v == nil {
// 			fmt.Println("No answer found. We can differentiate this from the empty value")
// 		}
// 		fmt.Println("The key exists")
// 		return nil
// 	})
// }

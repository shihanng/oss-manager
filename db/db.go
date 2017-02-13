package db

import "github.com/boltdb/bolt"

var projectBucket = []byte("project")

func NewDB(path string) (*DB, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &DB{b: db}, nil
}

type DB struct {
	b *bolt.DB
}

func (db *DB) Close() {
	defer db.b.Close()
}

func (db *DB) Path() string {
	return db.b.Path()
}

func (db *DB) AddProject(name, url string) error {
	return db.b.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(projectBucket)
		if err != nil {
			return err
		}

		return b.Put([]byte(name), []byte(url))
	})
}

type Project struct {
	Name string
	URL  string
}

func (db *DB) ListProjects() ([]Project, error) {
	var projects []Project

	db.b.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(projectBucket)

		b.ForEach(func(k, v []byte) error {
			p := Project{
				Name: string(k),
				URL:  string(v),
			}
			projects = append(projects, p)
			return nil
		})

		return nil
	})

	return projects, nil
}

package db

import (
	"bytes"
	"errors"
	"strings"

	"github.com/boltdb/bolt"
)

var (
	projectBucket = []byte("project")
	versionBucket = []byte("version")
)

var (
	errBadVersion      = errors.New("db: bad format of version")
	errProjectNotFound = errors.New("db: project not found")
)

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
	Name     string
	URL      string
	Versions []string
}

func (db *DB) ListProjects() ([]Project, error) {
	var projects []Project

	db.b.View(func(tx *bolt.Tx) error {
		pb := tx.Bucket(projectBucket)
		vc := tx.Bucket(versionBucket).Cursor()

		return pb.ForEach(func(name, url []byte) error {
			p := Project{
				Name: string(name),
				URL:  string(url),
			}

			for k, v := vc.Seek(name); k != nil && bytes.HasPrefix(k, name); k, v = vc.Next() {
				version := bytes.TrimPrefix(k, name)
				version = append(version, v...)
				p.Versions = append(p.Versions, string(version[1:]))
			}

			projects = append(projects, p)

			return nil
		})
	})

	return projects, nil
}

func (db *DB) UpdateVersion(name, version string) error {
	i := strings.LastIndex(version, ".")
	if i < 0 {
		return errBadVersion
	}

	major := version[:i+1]
	minor := version[i+1:]

	return db.b.Update(func(tx *bolt.Tx) error {
		pb := tx.Bucket(projectBucket)
		if v := pb.Get([]byte(name)); v == nil {
			return errProjectNotFound
		}

		vb, err := tx.CreateBucketIfNotExists(versionBucket)
		if err != nil {
			return err
		}

		name += ":" + major
		return vb.Put([]byte(name), []byte(minor))
	})
}

package db

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/boltdb/bolt"
)

var (
	projectsBucket = []byte("projects")
	versionBucket  = []byte("version")
	latestBucket   = []byte("latest")

	urlKey = []byte("url")
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
	tx, err := db.b.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ps, err := tx.CreateBucketIfNotExists(projectsBucket)
	if err != nil {
		return err
	}

	p, err := ps.CreateBucketIfNotExists([]byte(name))
	if err != nil {
		return err
	}

	if err = p.Put(urlKey, []byte(url)); err != nil {
		return err
	}

	return tx.Commit()
}

type Project struct {
	Name     string
	URL      string
	Versions []string
}

func (db *DB) ListProjects() ([]Project, error) {
	var projects []Project

	err := db.b.View(func(tx *bolt.Tx) error {
		ps := tx.Bucket(projectsBucket)
		c := ps.Cursor()

		for projectName, _ := c.First(); projectName != nil; projectName, _ = c.Next() {
			p := ps.Bucket(projectName)
			url := p.Get(urlKey)

			project := Project{
				Name: string(projectName),
				URL:  string(url),
			}

			v := p.Bucket(versionBucket)
			if v == nil {
				projects = append(projects, project)
				continue
			}

			if err := v.ForEach(func(major, minor []byte) error {
				version := fmt.Sprintf("%s%s", major, minor)
				project.Versions = append(project.Versions, version)
				return nil
			}); err != nil {
				return err
			}

			projects = append(projects, project)
		}
		return nil
	})

	return projects, err
}

func (db *DB) UpdateVersion(name, version string) error {
	i := strings.LastIndex(version, ".")
	if i < 0 {
		return errBadVersion
	}

	major := version[:i+1]
	minor := version[i+1:]

	return db.b.Update(func(tx *bolt.Tx) error {
		p := tx.Bucket(projectsBucket).Bucket([]byte(name))
		if p == nil {
			return errProjectNotFound
		}

		l, err := tx.CreateBucketIfNotExists(latestBucket)
		if err != nil {
			return err
		}

		pl, err := l.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return err
		}

		id, _ := pl.NextSequence()
		if err := pl.Put(itob(id), []byte(major+minor)); err != nil {
			return err
		}

		v, err := p.CreateBucketIfNotExists(versionBucket)
		if err != nil {
			return err
		}

		return v.Put([]byte(major), []byte(minor))
	})
}

func (db *DB) FirstLatest() (proj Project, err error) {
	err = db.b.View(func(tx *bolt.Tx) error {
		l := tx.Bucket(latestBucket)
		if l == nil {
			return bolt.ErrBucketNotFound
		}

		c := l.Cursor()
		project, _ := c.First()
		if project == nil {
			return errProjectNotFound
		}

		p := tx.Bucket(projectsBucket).Bucket([]byte(project))
		if p == nil {
			return errProjectNotFound
		}

		proj.Name = string(project)
		proj.URL = string(p.Get(urlKey))

		pl := l.Bucket(project)
		return pl.ForEach(func(_, version []byte) error {
			proj.Versions = append(proj.Versions, string(version))
			return nil
		})
	})
	return
}

func (db *DB) DeleteLatest(projectName string) error {
	return db.b.Update(func(tx *bolt.Tx) error {
		l := tx.Bucket(latestBucket)
		if l == nil {
			return bolt.ErrBucketNotFound
		}

		return l.DeleteBucket([]byte(projectName))
	})
}

// itob returns an 8-byte big endian representation of v.
func itob(x uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, x)
	return b
}

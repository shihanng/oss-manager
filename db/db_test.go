package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjects(t *testing.T) {
	db, err := NewDB("test.db")
	assert.NoError(t, err)
	defer db.Close()
	defer os.Remove(db.Path())

	err = db.AddProject("Project Name", "https://www.example.com")
	assert.NoError(t, err)

	expected := []Project{
		{
			Name: "Project Name",
			URL:  "https://www.example.com",
		},
	}
	actual, err := db.ListProjects()
	assert.NoError(t, err)
	assert.Equal(t, actual, expected)
}

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

	{
		err := db.UpdateVersion("Project Name", "1.0.1-a")
		assert.NoError(t, err)
	}
	{
		err := db.UpdateVersion("Project Name", "2.0.0")
		assert.NoError(t, err)
	}
	{
		err := db.UpdateVersion("Project Name", "1-0-1-a")
		assert.Equal(t, err, errBadVersion)
	}
	{
		err := db.UpdateVersion("Project-Name", "1.0.1-a")
		assert.Equal(t, err, errProjectNotFound)
	}

	expected := []Project{
		{
			Name:     "Project Name",
			URL:      "https://www.example.com",
			Versions: []string{"1.0.1-a", "2.0.0"},
		},
	}
	actual, err := db.ListProjects()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

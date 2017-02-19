package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjects(t *testing.T) {
	assert := assert.New(t)

	db, err := NewDB("test.db")
	assert.NoError(err)
	defer db.Close()
	defer os.Remove(db.Path())

	assert.NoError(db.AddProject("Project One", "https://one.example.com"))
	assert.NoError(db.AddProject("Project Two", "https://two.example.com"))
	assert.NoError(db.AddProject("Project Three", "https://three.example.com"))

	testCases := []struct {
		projectName    string
		projectVersion string
		expectedErr    error
	}{
		{"Project One", "1.0.1-a", nil},
		{"Project One", "2.0.0", nil},
		{"Project Two", "3.0.0", nil},
		{"Project-One", "1.0.1-a", errProjectNotFound},
		{"Project One", "1-0-1-a", errBadVersion},
	}
	for _, tc := range testCases {
		err := db.UpdateVersion(tc.projectName, tc.projectVersion)
		assert.Equal(tc.expectedErr, err)
	}

	expected := []Project{
		{
			Name:     "Project One",
			URL:      "https://one.example.com",
			Versions: []string{"1.0.1-a", "2.0.0"},
		},
		{
			Name: "Project Three",
			URL:  "https://three.example.com",
		},
		{
			Name:     "Project Two",
			URL:      "https://two.example.com",
			Versions: []string{"3.0.0"},
		},
	}
	actual, err := db.ListProjects()
	assert.NoError(err)
	assert.Equal(expected, actual)

	expectedUpdates := []Project{
		{
			Name:     "Project One",
			URL:      "https://one.example.com",
			Versions: []string{"1.0.1-a", "2.0.0"},
		},
		{
			Name:     "Project Two",
			URL:      "https://two.example.com",
			Versions: []string{"3.0.0"},
		},
	}

	for _, expected := range expectedUpdates {
		actual, err := db.FirstLatest()
		assert.NoError(err)
		assert.Equal(expected, actual)

		assert.NoError(db.DeleteLatest(actual.Name))
	}

	// Should be empty.
	project, err := db.FirstLatest()
	assert.Zero(project)
	assert.EqualError(err, errProjectNotFound.Error())

	assert.EqualError(db.DeleteLatest(project.Name), "incompatible value")
}

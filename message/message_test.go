package message

import (
	"testing"

	"github.com/shihanng/oss-manager/db"
	"github.com/stretchr/testify/assert"
)

func TestForUpdate(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		project  db.Project
		expected string
	}{
		{
			project: db.Project{
				Name:     "Project One",
				URL:      "http://one.example.com",
				Versions: []string{"1", "2", "3"},
			},
			expected: `Updates for Project One are available:

  - 1
  - 2
  - 3

 http://one.example.com
`,
		},
		{
			project: db.Project{
				Name:     "Project One",
				URL:      "http://one.example.com",
				Versions: []string{"1"},
			},
			expected: `An update for Project One is available:

  - 1

 http://one.example.com
`,
		},
	}

	for _, tc := range testCases {
		actual, err := ForUpdate(tc.project)
		assert.Equal(tc.expected, actual)
		assert.NoError(err)
	}
}

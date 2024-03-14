package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getValidFileName(t *testing.T) {
	var tests = []struct {
		in   string
		want string
	}{
		{"Go Is Awesome!", "Go Is Awesome"},
		{`ТСЖ "Сабурова, 19"`, `ТСЖ Сабурова 19`},
	}

	for _, test := range tests {
		assert.Equal(t, test.want, GetValidFileName(test.in))
	}
}

package weather

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_makeLocLatLong(t *testing.T) {
	s := `Pune	18.51957	73.85535
Agrahat	20.55241	85.90773`
	lll, err := makeLocLatLong(strings.NewReader(s))
	if assert.NoError(t, err) && assert.Len(t, lll, 2) {
		ll, ok := lll["Pune"]
		assert.True(t, ok)
		assert.Equal(t, "18.51957", ll.Lat)
		assert.Equal(t, "73.85535", ll.Long)

		ll, ok = lll["Agrahat"]
		assert.True(t, ok)
		assert.Equal(t, "20.55241", ll.Lat)
		assert.Equal(t, "85.90773", ll.Long)

		_, ok = lll["Abcd"]
		assert.False(t, ok)
	}
}

package weather

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/pkg/errors"
)

type LatLong struct {
	Lat  string
	Long string
}

type LocLatLong map[string]LatLong

func buildLocLatLong(fileName string) (LocLatLong, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return makeLocLatLong(f)
}

func makeLocLatLong(r io.Reader) (LocLatLong, error) {
	csvr := csv.NewReader(r)
	csvr.Comma = '	'
	records, err := csvr.ReadAll()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	lll := LocLatLong{}
	for _, r := range records {
		lll[r[0]] = LatLong{r[1], r[2]}
	}
	return lll, nil
}

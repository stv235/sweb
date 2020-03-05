package export

import (
	"log"
	"strconv"
	"strings"
	"time"
)

/*
Parses a filename of format <prefix><year>_<month>.zip
*/
func ParseFileName(filename, prefix string) (int, time.Month) {
	if !strings.HasPrefix(filename, prefix) {
		log.Panicln("[LOGIC]", "invalid prefix, expected '" + prefix + "'")
	}

	if !strings.HasSuffix(filename, ".zip") {
		log.Panicln("[LOGIC]", "invalid suffix, expected '.zip'")
	}

	yearMonth := filename
	yearMonth = strings.TrimPrefix(yearMonth, prefix)
	yearMonth = strings.TrimSuffix(yearMonth, ".zip")

	parts := strings.Split(yearMonth, "_")

	if len(parts) != 2 {
		log.Panicln("[LOGIC]", "invalid year_month format, expected '_'")
	}

	year, err := strconv.ParseInt(parts[0], 10, 64)

	if err != nil || year < 0 {
		log.Panicln("[LOGIC]", "invalid year value")
	}

	month, err := strconv.ParseInt(parts[1], 10, 64)

	if err != nil || month < 1 || month > 12 {
		log.Panicln("[LOGIC]", "invalid month value")
	}

	return int(year), time.Month(month)
}

package abel

import (
	"errors"
	"fmt"

	"github.com/buger/jsonparser"
	"github.com/parnurzeal/gorequest"
)

// Abel is the client used to contact Abel to query for the required metrics
type Abel struct {
	URL string
}

var request = gorequest.New()

// Get the count of the metric associated with it's tags
func (w *Abel) Get(metric string, tags []string) (int64, jsonparser.ValueType, error) {
	return w.GetCount(metric, tags, 0, 0, 0)
}

// GetCount gets the count of the metric associated with it's tags. We can filter the responses by
// further start & end. We can also define the window of aggregation using the duration parameter.
func (w *Abel) GetCount(metric string, tags []string, start int64, end int64, duration int64) (int64, jsonparser.ValueType, error) {
	requestSoFar := request.
		Get(w.URL + "/metrics").
		Query("name=" + metric).
		Query(fmt.Sprintf("duration=%d", duration)).
		Query(fmt.Sprintf("start=%d", start)).
		Query(fmt.Sprintf("end=%d", end))

	for _, tag := range tags {
		requestSoFar.Query("tags=" + tag)
	}

	response, body, errs := requestSoFar.End()
	if response.StatusCode != 200 {
		return -1, jsonparser.Unknown, fmt.Errorf("Status: %d, Body: %s", response.StatusCode, body)
	}

	if len(errs) > 0 {
		return -1, jsonparser.Unknown, combineErrors(errs)
	}

	value, datatype, _, err := jsonparser.Get([]byte(body), "responses", "[0]", "value", "value")
	if err != nil {
		return -1, datatype, err
	}

	valueInInt, err := jsonparser.ParseInt(value)
	if err != nil {
		return -1, datatype, err
	}

	return valueInInt, datatype, err
}

func combineErrors(errs []error) error {
	if len(errs) == 1 {
		return errs[0]
	} else if len(errs) > 1 {
		msg := "Error(s):"
		for _, err := range errs {
			msg += " " + err.Error()
		}
		return errors.New(msg)
	} else {
		return nil
	}
}

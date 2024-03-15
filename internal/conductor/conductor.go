package conductor

import (
	"time"

	"github.com/kanthorlabs/common/validator"
	dbentities "github.com/kanthorlabs/kanthor/internal/database/entities"
	dsentities "github.com/kanthorlabs/kanthor/internal/datastore/entities"
)

type Destination struct {
	Endpoint *dbentities.Endpoint
	SignKey  string
	Routes   []*dbentities.Route
}

func (d *Destination) Validate() error {
	return validator.Validate(
		validator.PointerNotNil("Endpoint", d.Endpoint),
		validator.StringRequired("SignKey", d.SignKey),
		validator.SliceRequired("Routes", d.Routes),
	)
}

func Many(
	msg *dsentities.Message,
	destinations map[string]*Destination,
	now time.Time,
) (map[string]*dsentities.Request, error) {
	requests := map[string]*dsentities.Request{}

	for _, destination := range destinations {
		req, err := One(msg, destination, now)
		if err != nil {
			return nil, err
		}

		if req != nil {
			requests[destination.Endpoint.Id] = req
		}
	}

	return requests, nil
}

// Request is a function that generates a request based on the given message and destination
// It only generates request if at least one non-exclusionary condition is matched
func One(
	msg *dsentities.Message,
	destination *Destination,
	now time.Time,
) (*dsentities.Request, error) {
	// if any destination got error, reject the whole request
	if err := destination.Validate(); err != nil {
		return nil, err
	}

	for i := range destination.Routes {
		source := &ConditionSource{Source: destination.Routes[i].ConditionSource}
		data := source.ExtractMessage(msg)

		expression := &ConditionExression{Expression: destination.Routes[i].ConditionExpression}
		matched, err := expression.Match(data)
		if err != nil {
			return nil, err
		}

		if !matched {
			continue
		}

		// if the condition is matched, and it's exclusionary, then no request will be generate
		if destination.Routes[i].Exclusionary {
			return nil, nil
		}

		return Request(msg, destination.Endpoint, destination.Routes[i], now, destination.SignKey), nil
	}

	return nil, nil
}

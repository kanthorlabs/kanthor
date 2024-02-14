package routing_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kanthorlabs/kanthor/internal/entities"
	"github.com/kanthorlabs/kanthor/internal/routing"
	"github.com/kanthorlabs/kanthor/internal/tester"
	"github.com/kanthorlabs/kanthor/mocks/timer"
	"github.com/stretchr/testify/assert"
)

func TestPlanRequests(t *testing.T) {
	now := time.Now().UTC()
	timer := timer.NewTimer(t)
	timer.On("Now").Return(now)

	app := tester.Application(timer)
	ep := tester.EndpointOfApp(timer, app)

	t.Run("success", func(st *testing.T) {
		route := &routing.Route{Endpoint: ep, Rules: make([]entities.EndpointRule, 0)}
		route.Rules = append(route.Rules, entities.EndpointRule{
			ConditionSource:     routing.ConditionSourceType,
			ConditionExpression: fmt.Sprintf("%s%s", routing.ConditionExpressionAny, routing.ConditionExpressionDivider),
		})

		msg := tester.MessageOfApp(timer, app)
		req, trace := routing.PlanRequest(timer, msg, route)
		assert.True(st, len(trace) == 0)
		assert.NotNil(st, req)
	})

	t.Run("error of condition expression", func(st *testing.T) {
		route := &routing.Route{Endpoint: ep, Rules: make([]entities.EndpointRule, 0)}
		route.Rules = append(route.Rules, entities.EndpointRule{
			ConditionExpression: string(routing.ConditionExpressionDivider[0]),
		})

		msg := tester.MessageOfApp(timer, app)

		_, trace := routing.PlanRequest(timer, msg, route)
		assert.True(st, len(trace) > 0)
		assert.Equal(st, "ROUTING.PLAN.RULE.CE.ERROR", trace[0])
	})

	t.Run("error of condition source empty", func(st *testing.T) {
		route := &routing.Route{Endpoint: ep, Rules: make([]entities.EndpointRule, 0)}
		route.Rules = append(route.Rules, entities.EndpointRule{
			ConditionSource:     routing.ConditionSourceType,
			ConditionExpression: fmt.Sprintf("%s%s", routing.ConditionExpressionAny, routing.ConditionExpressionDivider),
		})

		msg := tester.MessageOfApp(timer, app)
		msg.Type = ""
		_, trace := routing.PlanRequest(timer, msg, route)
		assert.True(st, len(trace) > 0)
		assert.Equal(st, "ROUTING.PLAN.RULE.CS.EMPTY.ERROR", trace[0])
	})

	t.Run("error of exclustionary check", func(st *testing.T) {
		route := &routing.Route{Endpoint: ep, Rules: make([]entities.EndpointRule, 0)}
		route.Rules = append(route.Rules, entities.EndpointRule{
			ConditionSource:     routing.ConditionSourceType,
			ConditionExpression: fmt.Sprintf("%s%s", routing.ConditionExpressionAny, routing.ConditionExpressionDivider),
			Exclusionary:        true,
		})

		msg := tester.MessageOfApp(timer, app)
		_, trace := routing.PlanRequest(timer, msg, route)
		assert.True(st, len(trace) > 0)
		assert.Equal(st, "ROUTING.PLAN.RULE.EXCLUSIONARY.ERROR", trace[0])
	})

	t.Run("error of not matched any rule", func(st *testing.T) {
		route := &routing.Route{Endpoint: ep, Rules: make([]entities.EndpointRule, 0)}
		route.Rules = append(
			route.Rules,
			entities.EndpointRule{
				ConditionSource:     routing.ConditionSourceType,
				ConditionExpression: fmt.Sprintf("%s%sunable", routing.ConditionExpressionEqual, routing.ConditionExpressionDivider),
				Exclusionary:        false,
			},
			entities.EndpointRule{
				ConditionSource:     routing.ConditionSourceAppId,
				ConditionExpression: fmt.Sprintf("%s%sunable", routing.ConditionExpressionEqual, routing.ConditionExpressionDivider),
				Exclusionary:        false,
			})

		msg := tester.MessageOfApp(timer, app)
		_, trace := routing.PlanRequest(timer, msg, route)
		assert.True(st, len(trace) > 0)
		assert.Equal(st, "ROUTING.PLAN.NOT_MATCH.ERROR", trace[0])
	})
}

func TestPlanRequest(t *testing.T) {
	now := time.Now().UTC()
	timer := timer.NewTimer(t)
	timer.On("Now").Return(now)

	app := tester.Application(timer)

	t.Run("success", func(st *testing.T) {
		routes := make(map[string]*routing.Route)

		match := tester.EndpointOfApp(timer, app)
		routes[match.Id] = &routing.Route{
			Endpoint: match,
			Rules: []entities.EndpointRule{
				{
					ConditionSource:     routing.ConditionSourceType,
					ConditionExpression: fmt.Sprintf("%s%s", routing.ConditionExpressionAny, routing.ConditionExpressionDivider),
				},
			},
		}

		matchother := tester.EndpointOfApp(timer, app)
		routes[matchother.Id] = &routing.Route{
			Endpoint: matchother,
			Rules: []entities.EndpointRule{
				{
					ConditionSource:     routing.ConditionSourceType,
					ConditionExpression: fmt.Sprintf("%s%s", routing.ConditionExpressionAny, routing.ConditionExpressionDivider),
				},
			},
		}

		not := tester.EndpointOfApp(timer, app)
		routes[not.Id] = &routing.Route{
			Endpoint: tester.EndpointOfApp(timer, app),
			Rules: []entities.EndpointRule{
				{
					ConditionSource:     routing.ConditionSourceType,
					ConditionExpression: fmt.Sprintf("%s%s", routing.ConditionExpressionAny, routing.ConditionExpressionDivider),
					Exclusionary:        true,
				},
			},
		}

		msg := tester.MessageOfApp(timer, app)
		reqs, traces := routing.PlanRequests(timer, msg, routes)
		assert.True(st, len(traces) == 1)
		assert.True(st, len(reqs) == 2)
	})
}

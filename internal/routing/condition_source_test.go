package routing_test

import (
	"testing"

	"github.com/jaswdr/faker"
	"github.com/kanthorlabs/common/idx"
	"github.com/kanthorlabs/kanthor/internal/entities"
	"github.com/kanthorlabs/kanthor/internal/routing"
	"github.com/stretchr/testify/assert"
)

func TestConditionSource(t *testing.T) {
	fake := faker.New()

	t.Run("error", func(st *testing.T) {
		st.Run("unknown cs", func(sst *testing.T) {
			cs := &entities.EndpointRule{ConditionSource: fake.Gender().Name()}
			assert.Empty(st, routing.ConditionSource(cs, nil))
		})

		st.Run("empty message type", func(sst *testing.T) {
			cs := &entities.EndpointRule{ConditionSource: routing.ConditionSourceType}
			assert.Empty(st, routing.ConditionSource(cs, &entities.Message{Type: ""}))
		})
	})

	msg := &entities.Message{
		AppId: idx.New(entities.IdNsApp),
		Type:  "testing.routing",
	}

	t.Run("return message type", func(st *testing.T) {
		cs := &entities.EndpointRule{ConditionSource: routing.ConditionSourceType}
		assert.Equal(st, msg.Type, routing.ConditionSource(cs, msg))
	})

	t.Run("return message app_id", func(st *testing.T) {
		cs := &entities.EndpointRule{ConditionSource: routing.ConditionSourceAppId}
		assert.Equal(st, msg.AppId, routing.ConditionSource(cs, msg))
	})
}

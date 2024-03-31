package api

import (
	"github.com/go-chi/chi/v5"
	httpxmw "github.com/kanthorlabs/common/gateway/httpx/middleware"
	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/kanthor/internal/datastore/entities"
	"github.com/kanthorlabs/kanthor/services/sdk/config"
	"github.com/kanthorlabs/kanthor/services/sdk/usecase"
)

func RegisterMessageRoutes(router chi.Router, service *sdk) {
	router.Route("/message", func(sr chi.Router) {
		sr.Use(httpxmw.Authz(service.infra.Gatekeeper(), config.ServiceName))
		// this API need achieve the best performance,
		// so we pass the application verification into the handler
		// By that way, we can apply cache technique to the application verification
		sr.Post("/", UseMessageCreate(service))
		sr.Route("/{id}", func(ssr chi.Router) {
			ssr.Get("/", UseMessageGet(service))
		})
	})
}

type Message struct {
	Id        string         `json:"id" example:"msg_2ePVr2tTfiJA20mN8wkc8EkGZu4"`
	CreatedAt int64          `json:"created_at" example:"1728925200000"`
	Tier      string         `json:"tier" example:"default"`
	AppId     string         `json:"app_id" example:"msg_2ePVr2tTfiJA20mN8wkc8EkGZu4"`
	Type      string         `json:"type" example:"testing.openapi"`
	Body      string         `json:"body" example:"{\"app_id\":\"msg_2ePVr2tTfiJA20mN8wkc8EkGZu4\",\"type\":\"testing.openapi\",\"object\":{\"from_client\":\"openapi\",\"say\":\"hello\"}}"`
	Metadata  *safe.Metadata `json:"metadata"  example:"kanthor.server.version:v2024.1014.1700" swaggertype:"object,string"`
} // @name Message

func (app *Message) Map(entity *entities.Message) {
	app.Id = entity.Id
	app.CreatedAt = entity.CreatedAt
	app.Tier = entity.Tier
	app.AppId = entity.AppId
	app.Type = entity.Type
	app.Body = entity.Body
	app.Metadata = &safe.Metadata{}
	app.Metadata.Merge(entity.Metadata)
}

type MessageEndpoint struct {
	Endpoint  *Endpoint   `json:"endpoint"`
	Requests  []*Request  `json:"requests"`
	Responses []*Response `json:"responses"`
} // @name MessageEndpoint

func (ep *MessageEndpoint) Map(endpoint *usecase.MessageEndpoint) {
	ep.Endpoint = &Endpoint{}
	ep.Endpoint.Map(endpoint.Endpoint)

	ep.Requests = make([]*Request, len(endpoint.Requests))
	for i := range endpoint.Requests {
		ep.Requests[i] = &Request{}
		ep.Requests[i].Map(endpoint.Requests[i])
	}

	ep.Responses = make([]*Response, len(endpoint.Responses))
	for i := range endpoint.Responses {
		ep.Responses[i] = &Response{}
		ep.Responses[i].Map(endpoint.Responses[i])
	}
}

type Request struct {
	Id        string `json:"id" example:"req_2ePVrMU69SGTlX0QC9Lvqkma82x"`
	CreatedAt int64  `json:"created_at" example:"1728925200000"`

	MsgId string `json:"msg_id" example:"msg_2ePVr2tTfiJA20mN8wkc8EkGZu4"`
	Tier  string `json:"tier" example:"default"`
	AppId string `json:"app_id" example:"msg_2ePVr2tTfiJA20mN8wkc8EkGZu4"`
	Type  string `json:"type" example:"testing.openapi"`
	Body  string `json:"body" example:"{\"app_id\":\"msg_2ePVr2tTfiJA20mN8wkc8EkGZu4\",\"type\":\"testing.openapi\",\"object\":{\"from_client\":\"openapi\",\"say\":\"hello\"}}"`

	EpId   string `json:"ep_id" example:"ep_2dZRCcnumVTMI9eHdmep89IpOgY"`
	Method string `json:"method" example:"POST"`
	Uri    string `json:"uri" example:"https://postman-echo.com/post"`

	Headers  *safe.Metadata `json:"headers" example:"Content-Type:application/json,Idempotency-Key:ik_2eR0d3ySDxK0ZjA35zdMswsF6HG,User-Agent:Kanthor/v2024.1014.1700,Webhook-Id:req_2ePVrMU69SGTlX0QC9Lvqkma82x,Webhook-Signature:v1=d0c41af2d916cf09225288ddebeb5fbb42a0a635f059b777bf4d4e787b3c5714,Webhook-Timestamp:1711806397376,Webhook-Type:testing.openapi" swaggertype:"object,string"`
	Metadata *safe.Metadata `json:"metadata" example:"kanthor.server.version:v2024.1014.1700,kanthor.rt.id:rt_2ePVcGq0hlMi1xBohNLuJvyIVHW" swaggertype:"object,string"`
} // @name Request

func (req *Request) Map(entity *entities.Request) {
	req.Id = entity.Id
	req.CreatedAt = entity.CreatedAt

	req.MsgId = entity.MsgId
	req.Tier = entity.Tier
	req.AppId = entity.AppId
	req.Type = entity.Type
	req.Body = entity.Method

	req.EpId = entity.EpId
	req.Method = entity.Method
	req.Uri = entity.Uri

	req.Headers = &safe.Metadata{}
	req.Headers.Merge(entity.Headers)
	req.Metadata = &safe.Metadata{}
	req.Metadata.Merge(entity.Metadata)
}

type Response struct {
	Id        string `json:"id" example:"res2eQwR88Z4xujEzwjUBmtQxK5uHB"`
	CreatedAt int64  `json:"created_at" example:"1728925200000"`

	MsgId string `json:"msg_id" example:"msg_2ePVr2tTfiJA20mN8wkc8EkGZu4"`
	Tier  string `json:"tier" example:"default"`
	AppId string `json:"app_id" example:"msg_2ePVr2tTfiJA20mN8wkc8EkGZu4"`
	Type  string `json:"type" example:"testing.openapi"`

	ReqId    string         `json:"req_id" example:"req_2ePVrMU69SGTlX0QC9Lvqkma82x"`
	EpId     string         `json:"ep_id" example:"ep_2dZRCcnumVTMI9eHdmep89IpOgY"`
	Method   string         `json:"method" example:"POST"`
	Headers  *safe.Metadata `json:"headers" example:"Content-Type:application/json" swaggertype:"object,string"`
	Metadata *safe.Metadata `json:"metadata"  example:"kanthor.server.version:v2024.1014.1700,kanthor.rt.id:rt_2ePVcGq0hlMi1xBohNLuJvyIVHW" swaggertype:"object,string"`

	Status int    `json:"status" example:"200"`
	Uri    string `json:"uri" example:"https://postman-echo.com/post"`
	Body   string `json:"body" example:"{\"args\":{},\"headers\":{},\"url\":\"https://postman-echo.com/post\"}"`
} // @name Response

func (res *Response) Map(entity *entities.Response) {
	res.Id = entity.Id
	res.CreatedAt = entity.CreatedAt

	res.MsgId = entity.MsgId
	res.Tier = entity.Tier
	res.AppId = entity.AppId
	res.Type = entity.Type

	res.ReqId = entity.ReqId
	res.Headers = &safe.Metadata{}
	res.Headers.Merge(entity.Headers)
	res.Metadata = &safe.Metadata{}
	res.Metadata.Merge(entity.Metadata)

	res.Status = entity.Status
	res.Uri = entity.Uri
	res.Body = entity.Body
}

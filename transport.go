package main

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"net/http"
)

type echoRequest struct {
	S string `json:"s"`
}

type echoResponse struct {
	V   string `json:"v"`
	Err string `json:"err, omitempty"`
}

func makeEchoEndpoint(svc EchoService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(echoRequest)
		v, err := svc.Echo(req.S)
		if err != nil {
			return echoResponse{v, err.Error()}, nil
		}
		return echoResponse{v, ""}, nil
	}
}

func decodeEchoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request echoRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

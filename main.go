package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

type Request struct {
	VideosIDs []string `json:"videos_ids"`
}

type Response struct {
	Message string `json:"message"`
}

func main() {
	// Connect to redis
	components, err := NewComponents()
	if err != nil {
		panic(err)
	}

	// Create download service
	svc := NewService(components)

	// Build serverEndpoints
	downloadVideosHandler := httptransport.NewServer(
		downloadVideosEndpoint(svc),
		decodeRequest,
		encodeResponse,
	)

	http.Handle("/video-download", downloadVideosHandler)
	http.ListenAndServe(":8080", nil)
}

func downloadVideosEndpoint(svc ServiceI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(Request)
		if !ok {
			panic(ok)
		}

		err := svc.DownloadVideos(req.VideosIDs)
		if err != nil {
			panic(err)
		}

		return Response{Message: "Downloads completed."}, nil
	}
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := Request{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

package main

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	svc := &addService{}
	sumEndpoint := makeSumEndpoint(svc)
	concatEndpoint := makeConcatEndpoint(svc)

	sumHandler := httptransport.NewServer(
		sumEndpoint,
		decodeSumRequest,
		encodeResponse,
	)

	concatHandler := httptransport.NewServer(
		concatEndpoint,
		decodeConcatRequest,
		encodeResponse,
	)

	http.Handle("/sum", sumHandler)
	http.Handle("/concat", concatHandler)

	http.ListenAndServe(":8080", nil)
}

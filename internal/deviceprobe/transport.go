package deviceprobe

import "errors"

type Response struct {
	Status int
	Name   string
}

func ToResponse(name string, err error) Response {
	if err == nil {
		return Response{Status: 200, Name: name}
	}
	if errors.Is(err, ErrDeviceMissing) {
		return Response{Status: 404}
	}
	return Response{Status: 500}
}

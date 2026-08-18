package deviceprobe

import "errors"

type Response struct {
	Status int
	Name   string
}

func ToResponse(name string, err error) Response {
	switch {
	case err == nil:
		return Response{Status: 200, Name: name}
	case errors.Is(err, ErrDeviceMissing):
		return Response{Status: 404}
	default:
		return Response{Status: 500}
	}
}

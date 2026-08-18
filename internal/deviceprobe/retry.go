package deviceprobe

import "context"

func Request(ctx context.Context, service *Service, code string, attempts *int) Response {
	var response Response
	for try := 0; try < 3; try++ {
		*attempts++
		name, err := service.Lookup(ctx, code)
		response = ToResponse(name, err)
		if response.Status < 500 {
			return response
		}
	}
	return response
}

package request

import (
	"strconv"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
)

func GenObjectLocations(num int, name string, ip string) *payload.Object_Locations {
	result := &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, num),
	}

	for i := 0; i < num; i++ {
		result.Locations[i] = &payload.Object_Location{
			Name: name,
			Uuid: "uuid-" + strconv.Itoa(i+1),
			Ips:  []string{ip},
		}
	}
	return result
}

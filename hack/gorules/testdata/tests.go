package target

import (
	"fmt"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
)

func payloadObjectAccesses() {
	sc := &payload.Search_Config{}
	_ = sc.Radius      // want `\QAvoid to access struct fields directly`
	_ = sc.GetRadius() // OK: use function to access the field

	loc := &payload.Object_Location{}
	_ = loc.Ips       // want `\QAvoid to access struct fields directly`
	_ = loc.Name      // want `\QAvoid to access struct fields directly`
	_ = loc.GetIps()  // OK: use function to access the field
	_ = loc.GetName() // OK: use function to access the field

	ireq := &payload.Insert_Request{}
	_ = ireq.Config  // want `\QAvoid to access struct fields directly`
	_ = ireq.Vector  // want `\QAvoid to access struct fields directly`
	_ = ireq.GetConfig()  // OK: use function to access the field
	_ = ireq.GetVector()  // OK: use function to access the field
}

func printFmts() {
	fmt.Println("%v", "") // want `\Qfound formatting directive in non-formatting call`
	fmt.Print("%s", "")   // want `\Qfound formatting directive in non-formatting call`

	fmt.Printf("%s", "") // OK: use Printf
}

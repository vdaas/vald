package target

import (
	"fmt"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
)

func payloadObjectAccesses() {
	sc := &payload.Search_Config{}
	_ = sc.Radius                  // want `\QAvoid to access struct fields directly`
	_, _ = sc.Epsilon, "test"      // want `\QAvoid to access struct fields directly`
	_, _ = "test", sc.Timeout      // want `\QAvoid to access struct fields directly`
	_ = sc.GetRadius()             // OK: use function to access the field
	_, _ = sc.GetEpsilon(), "test" // OK: use function to access the field
	_, _ = "test", sc.GetTimeout() // OK: use function to access the field

	loc := &payload.Object_Location{}
	_ = loc.Ips         // want `\QAvoid to access struct fields directly`
	_ = loc.Name        // want `\QAvoid to access struct fields directly`
	fmt.Print(loc.Ips)  // want `\QAvoid to access struct fields directly`
	fmt.Print(loc.Name) // want `\QAvoid to access struct fields directly`
	_ = loc.GetIps()    // OK: use function to access the field
	_ = loc.GetName()   // OK: use function to access the field

	fmt.Printf("%s, %%", "test: ", loc.Uuid) // want `\QAvoid to access struct fields directly`
	_, _, _ = "test", "test", loc.Uuid       // want `\QAvoid to access struct fields directly`
	_, _, _ = "test", "test", loc.GetUuid()  // OK: use function to access the field

	ireq := &payload.Insert_Request{}
	_ = ireq.Config      // want `\QAvoid to access struct fields directly`
	_ = ireq.Vector      // want `\QAvoid to access struct fields directly`
	_ = ireq.GetConfig() // OK: use function to access the field
	_ = ireq.GetVector() // OK: use function to access the field

	loc.Name = "newname" // OK: it is used in LHS
	ic := &payload.Insert_Config{}
	ireq.Config = ic // OK: it is used in LHS

	ireq.Config, _ = ic, "test" // OK: it is used in LHS
	ireq.Config.Timestamp = 0   // OK: it is used in LHS

	if loc.Name != "" { // want `\QAvoid to access struct fields directly`
	}

	if loc.Name == "" { // want `\QAvoid to access struct fields directly`
	} else {
	}

	if loc != nil && loc.Name != "" { // want `\QAvoid to access struct fields directly`
	}

	if ireq != nil && ireq.Vector.Id != "" { // want `\QAvoid to access struct fields directly`
	}
}

func printFmts() {
	fmt.Println("%v", "") // want `\Qfound formatting directive in non-formatting call`
	fmt.Print("%s", "")   // want `\Qfound formatting directive in non-formatting call`

	fmt.Printf("%s", "") // OK: use Printf
}

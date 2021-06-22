package gorules

import (
	"github.com/quasilyte/go-ruleguard/dsl"
)

func CheckPayloadObjectAccess(m dsl.Matcher) {
	m.Import("github.com/vdaas/vald/apis/grpc/v1/payload")

	m.Match(
		`$*_ = $*_, $x.$y, $*_`,
		`$*_, $x.$y, $*_`,
		`$x.$y != $_`,
		`$x.$y == $_`,
		`$x.$y.$_`,
	).Where(!m["y"].Text.Matches(`Get.+`) &&
		m["x"].Type.Implements(`payload.Payload`)).
		At(m["y"]).
		Report("Avoid to access struct fields directly").
		Suggest(`Get$y()`)
}

// from Ruleguard by example
// https://go-ruleguard.github.io/by-example/matching-comments
func PrintFmt(m dsl.Matcher) {
	// Here we scan a $s string for %s, %d and %v.
	m.Match(`fmt.Println($s, $*_)`,
		`fmt.Print($s, $*_)`).
		Where(m["s"].Text.Matches(`%[sdv]`)).
		Report("found formatting directive in non-formatting call")
}

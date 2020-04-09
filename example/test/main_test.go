package main_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/test"
	"github.com/vdaas/vald/internal/test/caser"
)

type client struct {
	timeout time.Duration
}

func (s *client) Send(addr string) (string, error) {
	return "", nil
}

func TestSend(t *testing.T) {
	test.New(
		test.WithCase(
			// Case_1（this case use deepequals. ）
			func() test.Caser {
				return caser.New(
					caser.WithName("test_case_1"),
					caser.WithArg("192.168.33.10", "80"),
					caser.WithField(time.Second),
					caser.WithWant("ok", nil),
				)
			}(),

			// Case_2 （this case not use deepequals. user check with `WithAssertFunc` method）
			func() test.Caser {
				return caser.New(
					caser.WithName("test_case_1"),
					caser.WithArg("192.168.33.10", "80"),
					caser.WithField(time.Second),
					caser.WithWant("ok", nil),
					caser.WithAssertFunc(
						func(gots, wants []interface{}) error {
							return errors.New("any error occurs")
						},
					),
				)
			}(),
		),

		// Test target
		test.WithDriverFunc(
			func(ctx context.Context, dp test.DataProvider) []interface{} {
				c := client{
					timeout: dp.Fields()[0].(time.Duration),
				}
				txt, err := c.Send(dp.Args()[0].(string))
				return []interface{}{txt, err}
			},
		),
	).Run(context.Background(), t)
}

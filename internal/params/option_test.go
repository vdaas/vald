package params

import (
	"testing"

	"github.com/vdaas/vald/internal/errors"
)

func TestWithConfigFilePathKey(t *testing.T) {
	type test struct {
		name      string
		key       string
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success",
			key:  "key",
			checkFunc: func(opt Option) error {
				got := new(parser)
				opt(got)

				if got.filePath.key != "key" {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithConfigFilePathKey(tt.key)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithConfigFilePathDefault(t *testing.T) {
	type test struct {
		name      string
		path      string
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success",
			path: "path",
			checkFunc: func(opt Option) error {
				got := new(parser)
				opt(got)

				if got.filePath.defaultPath != "path" {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithConfigFilePathDefault(tt.path)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithConfigFileDescription(t *testing.T) {
	type test struct {
		name      string
		desc      string
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success",
			desc: "desc",
			checkFunc: func(opt Option) error {
				got := new(parser)
				opt(got)

				if got.filePath.description != "desc" {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithConfigFileDescription(tt.desc)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithVersionKey(t *testing.T) {
	type test struct {
		name      string
		key       string
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success",
			key:  "key",
			checkFunc: func(opt Option) error {
				got := new(parser)
				opt(got)

				if got.version.key != "key" {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithVersionKey(tt.key)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithVersionFlagDefault(t *testing.T) {
	type test struct {
		name      string
		flag      bool
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success",
			flag: true,
			checkFunc: func(opt Option) error {
				got := new(parser)
				opt(got)

				if got.version.defaultFlag != true {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithVersionFlagDefault(tt.flag)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithVersionDescription(t *testing.T) {
	type test struct {
		name      string
		desc      string
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success",
			desc: "desc",
			checkFunc: func(opt Option) error {
				got := new(parser)
				opt(got)

				if got.version.description != "desc" {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithVersionDescription(tt.desc)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

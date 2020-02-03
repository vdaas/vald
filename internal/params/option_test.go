package params

import (
	"testing"

	"github.com/vdaas/vald/internal/errors"
)

func TestWithConfigFilePathKey(t *testing.T) {
	type test struct {
		name      string
		keys      []string
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success when keys is not empty",
			keys: []string{
				"key",
			},
			checkFunc: func(opt Option) error {
				got := new(parser)
				opt(got)

				if len(got.filePath.keys) != 1 || got.filePath.keys[0] != "key" {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "not set when keys is empty",
			keys: nil,
			checkFunc: func(opt Option) error {
				got := &parser{
					filePath: filePath{
						keys: []string{
							"key",
						},
					},
				}
				opt(got)

				if len(got.filePath.keys) != 1 || got.filePath.keys[0] != "key" {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithConfigFilePathKeys(tt.keys...)
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
			name: "set success when path is not empty",
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

		{
			name: "not set when path is empty",
			checkFunc: func(opt Option) error {
				got := &parser{
					filePath: filePath{
						defaultPath: "path",
					},
				}
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
			name: "set success when desc is not empty",
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

		{
			name: "not set when desc is empty",
			desc: "",
			checkFunc: func(opt Option) error {
				got := &parser{
					filePath: filePath{
						description: "desc",
					},
				}
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
		keys      []string
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success when keys is not empty",
			keys: []string{
				"key",
			},
			checkFunc: func(opt Option) error {
				got := new(parser)
				opt(got)

				if len(got.version.keys) != 1 && got.version.keys[0] != "key" {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},

		{
			name: "not set when keys is empty",
			keys: nil,
			checkFunc: func(opt Option) error {
				got := &parser{
					version: version{
						keys: []string{
							"key",
						},
					},
				}
				opt(got)

				if len(got.version.keys) != 1 && got.version.keys[0] != "key" {
					return errors.New("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithVersionKeys(tt.keys...)
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
			name: "set success when desc is not empty",
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

		{
			name: "not set when desc is empty",
			checkFunc: func(opt Option) error {
				got := &parser{
					version: version{
						description: "desc",
					},
				}
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

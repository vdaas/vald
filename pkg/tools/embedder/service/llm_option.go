//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package service

import (
	openai "github.com/sashabaranov/go-openai"
	"github.com/vdaas/vald/internal/errors"
)

// Option represent the functional option for ngt.
type OpenAIOption func(o *openAI) error

var defaultOpenAIOptions = []OpenAIOption{
	WithOpenAIModel("adav2"),
}

// WithOpenAIModel returns the functional option to set the openai model.
func WithOpenAIModel(mstr string) OpenAIOption {
	var model openai.EmbeddingModel
	switch mstr {
	case string(openai.AdaEmbeddingV2), "adav2", "ada-v2", "ada002", "ada-002", "ada":
		model = openai.AdaEmbeddingV2
	case string(openai.SmallEmbedding3), "small3", "small-3", "small003", "small-003", "small":
		model = openai.SmallEmbedding3
	case string(openai.LargeEmbedding3), "large3", "large-3", "large003", "large-003", "large":
		model = openai.LargeEmbedding3
	}
	return func(o *openAI) error {
		o.model = model
		return nil
	}
}

// WithToken returns the functional option to set the openai token.
func WithToken(token string) OpenAIOption {
	return func(o *openAI) error {
		if token == "" {
			return errors.New("token is empty")
		}
		o.token = token
		return nil
	}
}

// Copyright 2021 ecodeclub
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package retry

import (
	"testing"
	"time"

	"github.com/ecodeclub/ekit/internal/errs"

	"github.com/stretchr/testify/assert"
)

func TestFixedIntervalRetryStrategy_Next(t *testing.T) {

	testCases := []struct {
		name     string
		s        *FixedIntervalRetryStrategy
		interval time.Duration

		isContinue bool
	}{
		{
			name: "init case, retries 0",
			s: &FixedIntervalRetryStrategy{
				maxRetries: 3,
				interval:   time.Second,
			},

			interval:   time.Second,
			isContinue: true,
		},
		{
			name: "retries equals to MaxRetries 3 after the increase",
			s: &FixedIntervalRetryStrategy{
				maxRetries: 3,
				interval:   time.Second,
				retries:    2,
			},
			interval:   time.Second,
			isContinue: true,
		},
		{
			name: "retries over MaxRetries after the increase",
			s: &FixedIntervalRetryStrategy{
				maxRetries: 3,
				interval:   time.Second,
				retries:    3,
			},
			interval:   0,
			isContinue: false,
		},
		{
			name: "MaxRetries equals to 0",
			s: &FixedIntervalRetryStrategy{
				maxRetries: 0,
				interval:   time.Second,
			},
			interval:   time.Second,
			isContinue: true,
		},
		{
			name: "negative MaxRetries",
			s: &FixedIntervalRetryStrategy{
				maxRetries: -1,
				interval:   time.Second,
				retries:    0,
			},
			interval:   time.Second,
			isContinue: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			interval, isContinue := tt.s.Next()
			assert.Equal(t, tt.interval, interval)
			assert.Equal(t, tt.isContinue, isContinue)
		})
	}
}

func TestFixedIntervalRetryStrategy_New(t *testing.T) {
	testCases := []struct {
		name       string
		maxRetries int32
		interval   time.Duration

		want    *FixedIntervalRetryStrategy
		wantErr error
	}{
		{
			name:       "no error",
			maxRetries: 5,
			interval:   time.Second,

			want: &FixedIntervalRetryStrategy{
				maxRetries: 5,
				interval:   time.Second,
			},
			wantErr: nil,
		},
		{
			name:       "returns error, interval equals to 0",
			maxRetries: 5,
			interval:   0,

			want:    nil,
			wantErr: errs.NewErrInvalidIntervalValue(0),
		},
		{
			name:       "returns error, interval equals to -1",
			maxRetries: 5,
			interval:   -1,

			want:    nil,
			wantErr: errs.NewErrInvalidIntervalValue(-1),
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFixedIntervalRetryStrategy(tt.maxRetries, tt.interval)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

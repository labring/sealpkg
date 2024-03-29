// Copyright © 2023 sealos.
//
// Licensed under the Apache License, DefaultVersion 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import "testing"

func TestCompare(t *testing.T) {
	type args struct {
		v1 string
		v2 string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "default",
			args: args{
				v1: "v1.26",
				v2: "v1.26",
			},
			want: true,
		},
		{
			name: "default",
			args: args{
				v1: "v1.26.0",
				v2: "v1.26.1",
			},
			want: false,
		},
		{
			name: "default",
			args: args{
				v1: "v1.26.3",
				v2: "v1.26",
			},
			want: true,
		},
		{
			name: "default",
			args: args{
				v1: "v1.25.3",
				v2: "v1.26",
			},
			want: false,
		},
		{
			name: "default",
			args: args{
				v1: "v1.25.3",
				v2: "1.26",
			},
			want: false,
		},
		{
			name: "default",
			args: args{
				v1: "v4.1.5-rc1",
				v2: "v4.1.3",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Compare(tt.args.v1, tt.args.v2); got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToBigVersion(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "cc",
			args: args{v: "20.10"},
			want: "20.10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToBigVersion(tt.args.v); got != tt.want {
				t.Errorf("ToBigVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

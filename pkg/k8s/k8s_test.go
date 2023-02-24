package k8s

import "testing"

func TestFetchLatestVersion(t *testing.T) {
	type args struct {
		kubeVersion string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				kubeVersion: "v1.23",
			},
			want:    "1.23.16",
			wantErr: false,
		},
		{
			name: "default",
			args: args{
				kubeVersion: "v1.23.4",
			},
			want:    "1.23.4",
			wantErr: false,
		},
		{
			name: "default",
			args: args{
				kubeVersion: "1.23.6",
			},
			want:    "1.23.6",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FetchFinalVersion(tt.args.kubeVersion)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchFinalVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FetchFinalVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

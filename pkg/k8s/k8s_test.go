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
				kubeVersion: "1.23",
			},
			want:    "1.23.17",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetchVersion(tt.args.kubeVersion)
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

func Test_fetchAllVersion(t *testing.T) {
	c := fetchAllVersion("1.23.1")
	t.Log(c)
}

package k8sclient

import (
	"testing"

	v1 "k8s.io/api/apps/v1"
)

func TestListContains(t *testing.T) {
	type args struct {
		haystack interface{}
		needle   interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty - an empty list should not contain anything",
			args: args{
				haystack: v1.DeploymentList{
					Items: make([]v1.Deployment, 0),
				},
				needle: v1.Deployment{},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ListContains(tt.args.haystack, tt.args.needle); got != tt.want {
				t.Errorf("ListContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

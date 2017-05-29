package num

import "testing"

func TestSpaces(t *testing.T) {
	type args struct {
		num int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"     ", args{5}, "     "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Spaces(tt.args.num); got != tt.want {
				t.Errorf("Spaces() = %v, want %v", got, tt.want)
			}
		})
	}
}

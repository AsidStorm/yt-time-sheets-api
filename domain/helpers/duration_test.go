package helpers

import "testing"

func TestDurationFromString(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid date",
			args: args{
				in: "1w 2h",
			},
			want:    "P1WT2H",
			wantErr: false,
		},
		{
			name: "only time",
			args: args{
				in: "2h 3m 1s",
			},
			want:    "PT2H3M1S",
			wantErr: false,
		},
		{
			name: "only date",
			args: args{
				in: "2w 1d",
			},
			want:    "P2W1D",
			wantErr: false,
		},
		{
			name: "invalid parts",
			args: args{
				in: "1y 1m",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "invalid parts (doubled)",
			args: args{
				in: "1m 1m",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DurationFromString(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("DurationFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DurationFromString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

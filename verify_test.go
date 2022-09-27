package kurma

import (
	"crypto/ecdsa"
	"testing"
)

func TestVerifyToken(t *testing.T) {
	key := LoadKey("testdata/key")
	type args struct {
		token string
		key   *ecdsa.PrivateKey
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "good token",
			args: args{
				token: "3046022100e2e339d78f4cfde6a2a996aaf54092c228b5300b1b92016faf6f961f307bf030022100a3b7236487a01ec43e3593c000d6ccc7e28ee5b7c0e2c1f9e9af352693ef2055",
				key:   key,
			},
			want: true,
		},
		{
			name: "bad token",
			args: args{
				token: "3046022100e2e339d78f4cfde6a2a996aaf54092c228b5300b1b92016faf6f961f307bf030022100a3b7236487a01ec43e3593c000d6ccc7e28ee5b7c0e2c1f9e9af352693ef20ff",
				key:   key,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VerifyToken(tt.args.token, tt.args.key); got != tt.want {
				t.Errorf("VerifyToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

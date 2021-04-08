package console

import (
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestToken_IsValid(t1 *testing.T) {
	type fields struct {
		Value     string
		ExpiresAt time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "unsetValue",
			want: false,
			fields: fields{
				ExpiresAt: time.Now(),
			},
		},
		{
			name: "unsetExpiry",
			want: false,
			fields: fields{
				Value: "secret",
			},
		},
		{
			name: "expiredToken",
			want: false,
			fields: fields{
				Value:     "secret",
				ExpiresAt: time.Now().Add(-5 * time.Minute),
			},
		},
		{
			name: "validToken",
			want: true,
			fields: fields{
				Value:     "secret",
				ExpiresAt: time.Now().Add(5 * time.Minute),
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Token{
				Value:     tt.fields.Value,
				ExpiresAt: tt.fields.ExpiresAt,
			}
			if got := t.IsValid(); got != tt.want {
				t1.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToken_SignInURL(t1 *testing.T) {
	type fields struct {
		Value     string
		ExpiresAt time.Time
	}
	type args struct {
		dst string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *url.URL
		wantErr bool
	}{
		{
			name:    "invalidDestination",
			want:    nil,
			wantErr: true,
			fields: fields{
				Value:     "secret",
				ExpiresAt: time.Now().Add(5 * time.Minute),
			},
			args: args{
				dst: "example.com",
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Token{
				Value:     tt.fields.Value,
				ExpiresAt: tt.fields.ExpiresAt,
			}
			got, err := t.SignInURL(tt.args.dst)
			if (err != nil) != tt.wantErr {
				t1.Errorf("SignInURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("SignInURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

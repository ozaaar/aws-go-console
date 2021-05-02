package console

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

type mockedSTS struct {
	stsiface.STSAPI
}

func (*mockedSTS) GetFederationToken(*sts.GetFederationTokenInput) (*sts.GetFederationTokenOutput, error) {
	return &sts.GetFederationTokenOutput{
		Credentials: &sts.Credentials{
			AccessKeyId:     aws.String("foo"),
			SecretAccessKey: aws.String("bar"),
			SessionToken:    aws.String("foobar"),
		},
	}, nil
}

type mockedHTTPClient struct{}

func (*mockedHTTPClient) Do(*http.Request) (*http.Response, error) {
	return &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"SigninToken": "very secret token"}`)),
	}, nil
}

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

func TestConsole_SignInTokenWithArn(t *testing.T) {
	type fields struct {
		STS    stsiface.STSAPI
		Client HTTPClient
	}
	type args struct {
		name *string
		arn  *string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantValue string
		wantErr   bool
	}{
		{
			name: "getToken",
			fields: fields{
				STS:    &mockedSTS{},
				Client: &mockedHTTPClient{},
			},
			args: args{
				name: aws.String("foo"),
				arn:  aws.String("bar"),
			},
			wantValue: "very secret token",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Console{
				STS:    tt.fields.STS,
				Client: tt.fields.Client,
			}
			got, err := c.SignInTokenWithArn(tt.args.name, tt.args.arn)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignInTokenWithArn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Value, tt.wantValue) {
				t.Errorf("SignInTokenWithArn() got = %v, want %v", got, tt.wantValue)
			}
		})
	}
}

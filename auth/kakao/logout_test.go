package kakao

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

func TestLogout(t *testing.T) {
	tests := []struct {
		desc string
		req  *events.APIGatewayProxyRequest
		want int
	}{
		{
			"has cookie",
			&events.APIGatewayProxyRequest{Headers: map[string]string{"Cookie": "refresh_token=IepM4VixXMcP6V-eH-w-QC639h4GMpRlZQdz4go9dGkAAAF0K8iHcA"}},
			http.StatusOK,
		},
		{"no cookie", &events.APIGatewayProxyRequest{Headers: map[string]string{}}, http.StatusBadRequest},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			resp := Logout(test.req)
			if result := resp.StatusCode; result != test.want {
				t.Errorf("want: %d, got: %d", test.want, result)
			}
			// if response status is ok, check if cookie is expired
			if resp.StatusCode == http.StatusOK {
				// get date of expiry of cookie
				cookies := resp.MultiValueHeaders["Set-Cookie"]
				var c string
				for _, val := range cookies {
					if strings.HasPrefix(val, "refresh_token") {
						c = val
						break
					}
				}
				attrs := strings.Split(c, "; ")
				i := strings.Index(attrs[1], "=")
				e := attrs[1][i+1:]

				// parse it to time and get duration
				expires, _ := time.ParseInLocation(time.RFC1123, e, time.Local)
				period := time.Now().Local().Sub(expires)
				dur := int(period.Hours())
				if dur != 24 {
					t.Error("cookie set incorrectly")
				}
			}
		})
	}
}

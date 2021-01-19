package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/geo/s2"
)

func TestTopicName2Token(t *testing.T) {
	type args struct {
		topic string
	}
	type want struct {
		token string
		err   error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Test 01",
			args: args{
				topic: "/0/2/0/2/2/0/0/0/1/2/2/3/1/1/1/2/1/1/3/2/0/1/2/0/1/1/0/2/1/3/0",
			},
			want: want{
				token: "114035ab2f0c2939",
				err:   nil,
			},
		},
		{
			name: "Test 02",
			args: args{
				topic: "/1/2/2/2/1/3/0/1/2/1/0/0/2/3/0/2/1/3/2/2/2/2/0/0/0/0/2/1/1/2/0",
			},
			want: want{
				token: "3538c8593d5004b1",
				err:   nil,
			},
		},
		{
			name: "Test 03",
			args: args{
				topic: "/2/2/0/2/2/1/0/2/2/0/1/1/2/2/0/0/1/1/3/3/1/2/2/2/2/2/3/2/0/0/2",
			},
			want: want{
				token: "514942d02fb55705",
				err:   nil,
			},
		},
		{
			name: "Test 04",
			args: args{
				topic: "/3/1/0/3/0/3/0/2/3/3/0/3/1/3",
			},
			want: want{
				token: "699979bc",
				err:   nil,
			},
		},
		{
			name: "Test 05",
			args: args{
				topic: "/4/0/3/0/2/2/2/3/0/3/2/1/1/3",
			},
			want: want{
				token: "86559cbc",
				err:   nil,
			},
		},
		{
			name: "Test 06",
			args: args{
				topic: "/5/1/3/0/0/2/2/0/3/1/3/2/2/3/1/1/1/1/0/3/3/1/1/3/3/3/2/2/0/1/2",
			},
			want: want{
				token: "ae146f5aa9ebfd0d",
				err:   nil,
			},
		},
		{
			name: "Test 07",
			args: args{
				topic: "/0/1/1/3/2/2/3/1/0/0/2/2/2/3",
			},
			want: want{
				token: "0bd6855c",
				err:   nil,
			},
		},
		{
			name: "Test 08",
			args: args{
				topic: "0/1132231002223",
			},
			want: want{
				token: "0bd6855c",
				err:   nil,
			},
		},
		{
			name: "Test 09 (ERROR 01)",
			args: args{
				topic: "6/1132231002223",
			},
			want: want{
				token: "",
				err:   TopicNameError{""},
			},
		},
		{
			name: "Test 11 (ERROR 02)",
			args: args{
				topic: "0/1132236002223",
			},
			want: want{
				token: "",
				err:   TopicNameError{""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := TopicName2Token(tt.args.topic)
			if tt.want.err == nil {
				if err != nil {
					t.Errorf("Expected: %v, Result: %v", tt.want.err, err)
				}
				if token != tt.want.token {
					t.Errorf("Expected: %v, Result: %v", tt.want.token, token)
				}
			} else {
				if err == nil || reflect.ValueOf(err).Type() != reflect.ValueOf(tt.want.err).Type() {
					t.Errorf("Expected: %v, Result: %v", tt.want.err, err)
				}
			}
		})
	}

	// ランダムテスト
	rand.Seed(int64(20200118))
	for i := 0; i < 10000; i++ {
		t.Run(fmt.Sprintf("Random test %v", i), func(t *testing.T) {
			lat, lng := ((rand.Float64() * 360.) - 180.), ((rand.Float64() * 360.) - 180.)
			point := s2.CellIDFromLatLng(s2.LatLngFromDegrees(lat, lng))
			correctToken := point.ToToken()
			clientID := point.String()
			topic := celID2TopicName(point)

			convertedToken, _ := TopicName2Token(clientID)
			if correctToken != convertedToken {
				t.Errorf("Expected: %v, Result: %v", correctToken, convertedToken)
			}
			convertedToken, _ = TopicName2Token(topic)
			if correctToken != convertedToken {
				t.Errorf("Expected: %v, Result: %v", correctToken, convertedToken)
			}
		})
	}
}

func celID2TopicName(id s2.CellID) string {
	idString := strings.Replace(id.String(), "/", "", 1)
	return strings.Replace(idString, "", "/", len(idString))
}

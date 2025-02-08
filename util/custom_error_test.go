package util

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestNewCustomErrorWithKeys(t *testing.T) {
	type args struct {
		ctx  context.Context
		code string
		err  error
		keys map[string]string
	}
	tests := []struct {
		name string
		args args
		want *CustomError
	}{
		// TODO: Add test cases.
		{
			name: "New Error with no keys",
			args: args{
				ctx:  context.Background(),
				code: "code",
				err:  errors.New("error"),
				keys: map[string]string{},
			},
			want: &CustomError{
				Message: "error",
				Code:    "code",
			},
		},
		{
			name: "New Error with extra keys",
			args: args{
				ctx:  context.Background(),
				code: "code",
				err:  errors.New("error"),
				keys: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			want: &CustomError{
				Message: "error",
				Code:    "code",
			},
		},
		{
			name: "New Error with missing keys",
			args: args{
				ctx:  context.Background(),
				code: "code",
				err:  errors.New("error: {error}"),
				keys: map[string]string{},
			},
			want: &CustomError{
				Message: "error: {error}",
				Code:    "code",
			},
		},
		{
			name: "New Error with correct keys",
			args: args{
				ctx:  context.Background(),
				code: "code",
				err:  errors.New("{custom} error: {error}"),
				keys: map[string]string{
					"custom": "this is a custom message",
					"error":  "error occured while processing",
				},
			},
			want: &CustomError{
				Message: "this is a custom message error: error occured while processing",
				Code:    "code",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCustomErrorWithKeys(tt.args.ctx, tt.args.code, tt.args.err, tt.args.keys); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCustomErrorWithKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

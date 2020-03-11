package main

import (
	"fmt"
	"testing"
)

func Test_getMsdsByCas(t *testing.T) {
	type args struct {
		cas string
	}
	tests := []struct {
		name string
		args args
	}{{
		name: "maoyan",
		args: args{cas: "115-07-1"},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			json := getMsdsByCas(tt.args.cas)
			fmt.Printf("%+v", json)
		})
	}
}

package main

import (
	"testing"

	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
)

func TestCreateToAddress(t *testing.T) {
	testCases := []struct {
		desc   string
		addr   string
		want   []string
		hasErr bool
	}{
		{
			desc:   "happy case",
			addr:   "bob@bob.com",
			want:   []string{"bob@bob.com"},
			hasErr: false,
		},
		{
			desc:   "blank",
			addr:   "",
			hasErr: true,
		},
		{
			desc:   "malformed",
			addr:   "jhgf",
			hasErr: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual, err := toList(tC.addr)
			if tC.hasErr {
				then.AssertThat(t, err, is.Not(is.Nil()))
			} else {
				then.AssertThat(t, actual, is.EqualTo(tC.want))
			}
		})
	}
}

func TestBuildBody(t *testing.T) {
	testCases := []struct {
		desc    string
		subject string
		body    string
		want    string
	}{
		{
			desc:    "happy case",
			subject: "asubject",
			body:    "a body\nline 2",
			want: `Subject: asubject

a body
line 2`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := buildBody(tC.subject, tC.body)
			then.AssertThat(t, string(actual), is.EqualTo(tC.want))
		})
	}
}

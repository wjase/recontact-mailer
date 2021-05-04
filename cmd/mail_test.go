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
		to      []string
		from    string
		subject string
		body    string
		want    string
	}{
		{
			desc:    "happy case",
			subject: "asubject",
			to:      []string{"recipient@host.com"},
			from:    "frodo@baggins.com",
			body:    "a body\nline 2",
			want: `To: recipient@host.com
From: frodo@baggins.com
Subject: asubject
Content-Type: text/html; charset="UTF-8"
Content-Transfer-Encoding: base64

YSBib2R5CmxpbmUgMg==`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := buildBody(tC.to, tC.from, tC.subject, tC.body)
			then.AssertThat(t, whitespaceRemover.Replace(string(actual)), is.EqualTo(whitespaceRemover.Replace(tC.want)))
		})
	}
}

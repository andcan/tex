package tex

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestToGo(t *testing.T) {
	tests := []struct {
		name, want string
	}{
		{"foo_bar", "FooBar"},
		{"foo_bar_baz", "FooBarBaz"},
		{"Foo_bar", "FooBar"},
		{"foo_WiFi", "FooWiFi"},
		{"id", "ID"},
		{"Id", "ID"},
		{"foo_id", "FooID"},
		{"fooId", "FooID"},
		{"fooUid", "FooUID"},
		{"idFoo", "IDFoo"},
		{"uidFoo", "UIDFoo"},
		{"midIdDle", "MidIDDle"},
		{"APIProxy", "APIProxy"},
		{"ApiProxy", "APIProxy"},
		{"apiProxy", "APIProxy"},
		{"_Leading", "Leading"},
		{"___Leading", "Leading"},
		{"trailing_", "Trailing"},
		{"trailing___", "Trailing"},
		{"a_b", "AB"},
		{"a__b", "AB"},
		{"a___b", "AB"},
		{"Rpc1150", "RPC1150"},
		{"case3_1", "Case3_1"},
		{"case3__1", "Case3_1"},
		{"IEEE802_16bit", "IEEE802_16bit"},
		{"IEEE802_16Bit", "IEEE802_16Bit"},
		{"", ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			g := NewGomegaWithT(tt)

			got := ToGo(test.name)
			g.Expect(got).To(Equal(test.want))
		})
	}
}

func TestToGoPrivate(t *testing.T) {
	tests := []struct {
		name, want string
	}{
		{"foo_bar", "fooBar"},
		{"foo_bar_baz", "fooBarBaz"},
		{"Foo_bar", "fooBar"},
		{"foo_WiFi", "fooWiFi"},
		{"id", "id"},
		{"Id", "id"},
		{"foo_id", "fooID"},
		{"fooId", "fooID"},
		{"fooUid", "fooUID"},
		{"idFoo", "idFoo"},
		{"uidFoo", "uidFoo"},
		{"midIdDle", "midIDDle"},
		{"APIProxy", "apiProxy"},
		{"ApiProxy", "apiProxy"},
		{"apiProxy", "apiProxy"},
		{"_Leading", "_Leading"},
		{"___Leading", "_Leading"},
		{"trailing_", "trailing"},
		{"trailing___", "trailing"},
		{"a_b", "aB"},
		{"a__b", "aB"},
		{"a___b", "aB"},
		{"Rpc1150", "rpc1150"},
		{"case3_1", "case3_1"},
		{"case3__1", "case3_1"},
		{"IEEE802_16bit", "ieee802_16bit"},
		{"IEEE802_16Bit", "ieee802_16bit"},
		{"_", "_Arg"},
		{"", ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			g := NewGomegaWithT(tt)

			got := ToGoPrivate(test.name)
			g.Expect(got).To(Equal(test.want))
		})
	}
}

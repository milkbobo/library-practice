package client

import (
	. "github.com/fishedee/language"
	. "library/models/common"
	"testing"
)

type ClientLoginAoTest struct {
	BaseTest
	ClientLoginAo ClientLoginAoModel
}

func (this *ClientLoginAoTest) TestDuplicateName() {
	this.Benchmark(1000, 1000, func() {
		defer CatchCrash(func(e Exception) {

		})
		this.ClientLoginAo.Register("nnnnnn", "password", "password")
	})
}

func TestClientLoginAo(t *testing.T) {
	InitTest(t, &ClientLoginAoTest{})
}

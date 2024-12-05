//go:build integration

package mq_test

import (
	"os"
	"slices"
	"testing"

	"github.com/guackamolly/zero-monitor/internal/mq"
)

var rsaPem = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAw358uANE26qvIs51rjjBRo3v58MJFLja9/padm7iM+ZSchAj
LTNksmsJI1pbSRvA6gbmMxqqrvaWXu2TfQqponGc/iwBPzIa4dl8gHi1CLMME1Nl
7yJw+uKK15XGo/Qtaz7KYdtHCpwRs5Gc+8Ww7Cq7jEYPu+KL/TrrqxXldUV/Hocv
LRkydtRDeuSdefGjokTi2LEQspdpZ4HRbBncl/ntKLboN2R3RAomx1bWTqA/M0xx
QqOiHO9zai9yrYuFYUuSmeyusqv3wH2mhLrEhPPQYzG/BLvq4rmZautJDnbZV613
/WKzqg2hhMgwJtAzbm7yCMSlyTB69uNLHpVK9QIDAQABAoIBAADwsUsYMJ9lrnu3
8rWihsgNvYxd1wjDhIayQDKXFMcxgacn8AU3+lO7pw0s5hylyQXtRPPFvMp9Awm+
OU3KOhGhb2f5EuYnVPWwax9yl6GH9upTfowziIy5RJWZscSu8OoOeXHbRNpkGCYp
XNnhSizPXhxkLHwyuidmSC7gpzJv6Ysy6+3PJubro1IyBIuMmQdiKtVbWL9Razk6
IPhCMuOCknCHgZbOETh8NeqXh7Ztiragb58h5Fc8RSONdYsmkmfirWXFBArSs/Xw
YNyYZGwfZ5Tyt2TggA2iFBF6oPJtQk9cOnyLOhjmCLc7w+TF6crOPxZRx6usOkP2
ElduSBsCgYEA9gGloYDhK2oEmtU/3ZtI+NT9wEVW/W23aPaG4QqB7yI5XaW3zX9P
Mq3pAFxR9xjOhvpW8ed+f9qg+5toc73fLLuMDWFuA5Pr34dxV0DeFdNZ8h8qPBtz
j7kvUTtH2h+zTc4IvVpbcBKkqT2SJpr1wO3xThBajfOCh+8V/VdOwncCgYEAy2+J
A2EuPAjHtybyJpcZp/lyMWzwN4c3N511e/PvJqT0TSQjRrM9M4TREmX4eQghE2Va
QkKkm+RoN4mxFP7C6sjq7kmvFbZG4JC9kBtDqV79mmEGw8/IvaSTXXG0RL8tKGIg
bul8Q0VVWlZgkx780bII9TMYOZRZP9FUSicb7PMCgYAsnN3ZrRKomeBd5+BeIuQX
5CBkdu6wpO4HBfYt54bqxA0dM4lipfzJ1woTO6rNod0KU2njErU5IH/jQSqvGrbX
WOesIYge8/tpnRlr1mKwGJUOOKKjJeNOJCo1lAeSwf71VDD3jeRZLbhYzMatY5q/
syb4njSd25RHbI9TUzsAPwKBgDdy/C5uo5J7diwmsmPwVW7iX8y2+7a25UcEZQxX
Db1Dws7v5amUmz7amb3hC1u56oIF4xciYQmYtQtGPX0Sf4BNKTOv48gQObtl2DVa
KRQWLxuQDK78iKOgIwaaQl9mmGFkdaClhVg0orIPzxzqmlBxrV1gAt9W3wi0/ruD
c2ofAoGBAL+GDfuCrJPbnp37C4ZEMgWiJ73Mgu7tFVh53irll1IXp4pqtn7YHRNL
P7dplKCTY0RUyckE8WGqLl2sY5HrA3RrB+SzjFJcESsBrWxk56JIquEQPCY1kaFZ
+6J3fC5Nz8oA5jVUtyFGj156hP/A0grfWLQnauJpFNhOLS/VATBx
-----END RSA PRIVATE KEY-----`

var rsaPub = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAw358uANE26qvIs51rjjB
Ro3v58MJFLja9/padm7iM+ZSchAjLTNksmsJI1pbSRvA6gbmMxqqrvaWXu2TfQqp
onGc/iwBPzIa4dl8gHi1CLMME1Nl7yJw+uKK15XGo/Qtaz7KYdtHCpwRs5Gc+8Ww
7Cq7jEYPu+KL/TrrqxXldUV/HocvLRkydtRDeuSdefGjokTi2LEQspdpZ4HRbBnc
l/ntKLboN2R3RAomx1bWTqA/M0xxQqOiHO9zai9yrYuFYUuSmeyusqv3wH2mhLrE
hPPQYzG/BLvq4rmZautJDnbZV613/WKzqg2hhMgwJtAzbm7yCMSlyTB69uNLHpVK
9QIDAQAB
-----END PUBLIC KEY-----`

var rsaPemKeypath string
var rsaPubKeypath string

func init() {
	rsaPemFile, err := os.CreateTemp("", "")
	if err != nil {
		panic(err)
	}
	rsaPemFile.WriteString(rsaPem)
	defer rsaPemFile.Close()

	rsaPubFile, err := os.CreateTemp("", "")
	if err != nil {
		panic(err)
	}
	rsaPubFile.WriteString(rsaPub)
	defer rsaPemFile.Close()

	rsaPemKeypath = rsaPemFile.Name()
	rsaPubKeypath = rsaPubFile.Name()
}

func TestEncryptAndDecryptAsymmetric(t *testing.T) {
	data := []byte("zero-monitor")

	err := mq.LoadAsymmetricBlock(rsaPubKeypath)
	if err != nil {
		t.Fatalf("didn't expect load public key block to fail, %v", err)
	}

	encrypted, err := mq.EncryptAsymmetric(data)
	if err != nil {
		t.Fatalf("didn't expect encrypt asymmetric to fail, %v", err)
	}

	err = mq.LoadAsymmetricBlock(rsaPemKeypath)
	if err != nil {
		t.Fatalf("didn't expect load private key block to fail, %v", err)
	}

	decrypted, err := mq.DecryptAsymmetric(encrypted)
	if err != nil {
		t.Fatalf("didn't expect decrypt asymmetric to fail, %v", err)
	}

	if !slices.Equal(decrypted, data) {
		t.Errorf("expected decrypt to return %v, but got %v", data, decrypted)
	}
}

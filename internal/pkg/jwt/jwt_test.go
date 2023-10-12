package jwt_test

import (
	"testing"

	"github.com/ic3network/mccs-alpha/internal/pkg/jwt"
	"github.com/stretchr/testify/require"
)

const (
	TEST_PRIVATE_KEY = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA1vl/vENflHIgomF0qxfPg9l9FqwguCYtGNIbAfvVXSsLueAV
D23ZDtazisA67y+dO8TQk+KTGeCbM2Otcvvs7mdsRsbFe+demdSfydQHAp2tGb7O
DBQPGJeyCnyfePB2GmqjVrBB1BjjWgnwwNaXJQkmtLruG8Sgrwl3nKjTeC3x8LjB
l0gGU6UFON6SBF+/CovbOHn+P8eUC/LJrvX8dGpXfGoTzk18WKU5GzThrSgoGkL7
CToAN8/JoW8G/1gtiuhFRAi53oxpt058An8ORBP8PkxGj+enjq59C/5YMYLwfKHj
lylu0uRpfbmih4/46pgbYjjUyOhSXxHd1k5VmwIDAQABAoIBAQCMb29X0HeXJTtW
eO3be3GQA7to3Ud+pUnephsIn7iR5bYCVnXLn4ol3HJr2QpnCKbhzcAoa+KHDCi3
WI2NyS/Nynh8gAuw1sQBIFrGYaG2vsS/RdubHluCSE8B9MnFGuk8dp9/2SMX6K5V
Opsxjr4sbp7/gAJe14PU9Q1TpSKIpeBvH3li+mF405YIo6JTcxZToJtJcRMB7V7P
E8KI6F+cwczAmc282pMtpT8/TL1PHq7JYH8LJ18+EwQOvEvo6A8KJIDlFqEaK+yu
oQhHdUxxO9ZxAlTcb48g6o2laLd+v6k8N9yQagEyqHof/cSNXA24ZnqGHS3sECri
TpV1gDypAoGBAPLzXtxV2zseEDH5dYQ3vh6bdEudBFy/3YGFIMrtbMZCP8dx2Lie
P9ABWJqYEUrQpSNnKk6XdNQQGiuwFmygOuZMvyw4svuGKHIIQsSFQRtn+oO6lqSe
cOW+59UgBipBBjRuNvdSh3g4i+JI33bDwVedO5Qp+OinenVHMx27NHNNAoGBAOKF
c2/W7PYllHMGPIfjW1/+otkHdwPyLiBleamUgs33do8YGJekddX0+2BgR1ZIoIWQ
MsiWA/FcsTKERaZv220s7iz58w0GTcpbHQQW7e6D9cl+5DIXnEyG6vQ+hSiUOQXe
LjblgGQJHitrH2wUW/eEjQvXYLIduKlTcOWGZ6iHAoGBAJUtLJUcPsX4+rbE1wy9
cYa3q1v2aMROp0MtLGqOCJlf+muLkyghO0uMWAxszUlj/dJUOV0SkJDZ5kfnEo3W
gPQCMeyEUBozUUhbnCuxKr4aRW93NaKVCvt3EkECLebqEFZHSobobPg7uGDUoCn7
nw8eI4QhlY29sGqssk1SMq2NAoGBAIW12n8w8e0WH7uJ+d8IoJ5Yc44CbwlgQkQT
Qi6MoG2t3kj3I0UX6gqismOgUVuoQUC17pQioS8u1NYJ6AcnzfFy7SCVZhfRGcgR
4l3QnyAEuuf2xAKhlzxBA52q7fUXEVXaYZM8A36JN0rPz9t/ZQ4FKzDLMKPTEXa5
71E89iEvAoGASgN2hEjcF4lwe5ahgrLAheibPzC+6DFdSCBw9CbL183C5s3r4JDh
VDH3H2SHpB0qmBa+YLRwRvHxpWU9uq/unaJvc+AQ3JwZX3bQ8ixvyVfpeBXZF6Dh
KoQ0MewWRNtrpGFa5qdWBfcenKdhWgWrdMnroNhqCHfXEIiYsj3qqWs=
-----END RSA PRIVATE KEY-----`

	TEST_PUBLIC_KEY = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1vl/vENflHIgomF0qxfP
g9l9FqwguCYtGNIbAfvVXSsLueAVD23ZDtazisA67y+dO8TQk+KTGeCbM2Otcvvs
7mdsRsbFe+demdSfydQHAp2tGb7ODBQPGJeyCnyfePB2GmqjVrBB1BjjWgnwwNaX
JQkmtLruG8Sgrwl3nKjTeC3x8LjBl0gGU6UFON6SBF+/CovbOHn+P8eUC/LJrvX8
dGpXfGoTzk18WKU5GzThrSgoGkL7CToAN8/JoW8G/1gtiuhFRAi53oxpt058An8O
RBP8PkxGj+enjq59C/5YMYLwfKHjlylu0uRpfbmih4/46pgbYjjUyOhSXxHd1k5V
mwIDAQAB
-----END PUBLIC KEY-----`
)

func TestJWT(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		isAdmin bool
	}{
		{
			name:    "Valid User Admin",
			userID:  "123",
			isAdmin: true,
		},
		{
			name:    "Valid User Not Admin",
			userID:  "456",
			isAdmin: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("jwt.private_key", TEST_PRIVATE_KEY)
			t.Setenv("jwt.public_key", TEST_PUBLIC_KEY)
			j := jwt.NewJWTManager()

			token, err := j.Generate(tt.userID, tt.isAdmin)
			require.NoError(t, err)

			claims, err := j.Validate(token)
			require.NoError(t, err)

			require.Equal(t, tt.userID, claims.UserID)
			require.Equal(t, tt.isAdmin, claims.Admin)
		})
	}
}

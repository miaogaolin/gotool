package rsax

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublicDecrypt(t *testing.T) {
	data := "JAyMbv1wqoLPPp8fRY4KBMPfpNHto58NzOU2oMZf9W/tl/1AmEloPTlPpWmrTxY6z2ubbKl6wU34dDV3B5XhejkUgpw3DY9Rt0buAlqeOdEuKmaTfrxEtaGrdg7M0l91ttZzZobZ335dUXr6EXfw4iy5puaLQqNhxgDNDNxJYc9flwRr6QJzJbPxml1X3B83EwB4cJ10BtNpyUlDZ26ysn/Ja99f7MZuZpaMgCp2RtEv98yK0PN+RjNuom1MXZa406p1yy25AvGF9TZzrgkNh7QGxfGJ5g0lVjlTB8bDrBH9J+vDbdcnLNAcWLPUSq0ncpAmEPGOzYHjp2EEmixkmQ=="
	publicKey := "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwQ81QuYHJo4alQ9WURZH\n87W0sEdSC98Q/2nCXLjrix2Hmjoc7utwfqzj9+fKb9jLKfNWpU/kUaINhvja9IbV\nROMSjNvL/ECbwO3/NSybY5rn5LGyGN3ftAabCrJ77GyDzKZaYNOC+mzfANYHfoJX\nHcLJEXYJvv8jsQqVlnNyB2roj+Jznjl1wyHleSAVPpFYvzHnDe4zI0lvy+5PNxKZ\nRY1ltv13ngRdSCNE5MY6JTmcnwyvNEGVeHnaUfLUa0/yPQvqwZ36RrZHI2bB7NDB\nJ2rI3K2qD3tilB4lD/eKMb2Um5hMgRNeSeAIAC2R/upSuRz17qbEWfoX89Kl+R8R\niwIDAQAB\n-----END PUBLIC KEY-----"
	res, err := PublicDecrypt(data, publicKey)
	if err != nil {
		assert.Error(t, err)
	}
	expected := `{"date":"2020-12-16 18:22:17","rand":512,"ip":"192.168.1.185"}`
	assert.Equal(t, expected, res)
}

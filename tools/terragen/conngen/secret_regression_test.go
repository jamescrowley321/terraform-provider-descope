package conngen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEightByEightAPIKeysAreSecrets(t *testing.T) {
	for _, id := range []string{"eight-by-eight-viber", "eight-by-eight-whatsapp"} {
		c := Connector{ID: id, Fields: []*Field{{Name: "apiKey", Type: FieldTypeString, Required: true}}}
		c.Prepare()
		require.Equal(t, FieldTypeSecret, c.Fields[0].Type)
		require.Equal(t, `stringattr.SecretRequired()`, c.Fields[0].AttributeType())
	}
}

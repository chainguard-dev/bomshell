package shell

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadRecipeFile(t *testing.T) {
	for uCase, tc := range map[string]struct {
		data      []byte
		expected  string
		shouldErr bool
	}{
		"single line:": {[]byte("sboms[0].files().ToDocument()"), "sboms[0].files().ToDocument()\n", false},
		"multiline":    {[]byte("// Here are some comments!\nsboms[0].files().ToDocument()"), "// Here are some comments!\nsboms[0].files().ToDocument()\n", false},
		"shebang":      {[]byte("#!/usr/bin/bomshell\nsboms[0].files().ToDocument()"), "sboms[0].files().ToDocument()\n", false},
		"null string":  {[]byte(""), "", true},
		"just shebang": {[]byte("#!/usr/bin/bomshell\n"), "", true},
	} {
		i := DefaultBomshellImplementation{}
		buf := bytes.NewBuffer(tc.data)
		res, err := i.ReadRecipeFile(buf)
		if tc.shouldErr {
			require.Error(t, err, uCase)
			continue
		}

		require.Equal(t, tc.expected, res, uCase)
	}
}

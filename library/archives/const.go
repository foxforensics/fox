package archives

import "errors"

var ErrCorruptPassword = errors.New("archive corrupt or password wrong")

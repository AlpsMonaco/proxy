package vpn

import "errors"

var ErrServerRejected = errors.New("remote server reject connection")
var ErrPrefixNotMacth = errors.New("prefix not match")
var ErrVersionDismatch = errors.New("version not match")

// +build !cc1101,!rfm69

package medtronic

import "github.com/thecubic/cc111x"

var radioInterface = cc111x.Open

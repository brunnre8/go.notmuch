package notmuch

import "unsafe"

type cStruct struct {
	cptr unsafe.Pointer
}

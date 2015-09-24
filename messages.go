package notmuch

// Copyright © 2015 The go.bindings Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

// #cgo LDFLAGS: -lnotmuch
// #include <stdlib.h>
// #include <notmuch.h>
import "C"

// Messages represents notmuch messages.
type Messages struct {
	cptr   *C.notmuch_messages_t
	thread *Thread
}

func (ms *Messages) toC() *C.notmuch_messages_t {
	return (*C.notmuch_messages_t)(ms.cptr)
}
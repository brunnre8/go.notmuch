package notmuch

// Copyright © 2015 The go.bindings Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

var dbPath string

func init() {
	var err error
	dbPath, err = filepath.Abs("fixtures/database-v1")
	if err != nil {
		panic(err)
	}
}

func TestOpenNotFound(t *testing.T) {
	_, err := Open("/not-found", DBReadOnly)
	if err == nil {
		t.Errorf("Open(%q): expected error got nil", "/not-found")
	}
}

func TestCreate(t *testing.T) {
	tmp := os.TempDir()
	tmpDbPath := path.Join(tmp, ".notmuch")
	defer func() {
		os.RemoveAll(tmpDbPath)
	}()

	db, err := Create(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if want, got := 1, db.Version(); got < want {
		t.Errorf("db.Version(): want at least %d got %d", want, got)
	}
}

func TestOpen(t *testing.T) {
	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatalf("Open(%q): unexpected error: %s", dbPath, err)
	}
	defer db.Close()
	if want, got := 1, db.Version(); got < want {
		t.Errorf("db.Version(): want at least %d got %d", want, got)
	}
}

func TestLastStatus(t *testing.T) {
	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatalf("Open(%q): unexpected error: %s", dbPath, err)
	}
	defer db.Close()
	if want, got := "", db.LastStatus(); want != got {
		t.Errorf("db.LastStatus(): want %s got %s", want, got)
	}
	// TODO(kalbasit): use add_message later to cause an error and add a test for it
}

func TestPath(t *testing.T) {
	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatalf("Open(%q): unexpected error: %s", dbPath, err)
	}
	defer db.Close()
	if want, got := "go.notmuch/fixtures/database-v1", db.Path(); !strings.HasSuffix(got, want) {
		t.Errorf("db.Path(): want %s got %s", want, got)
	}
}

func TestNeedsUpgrade(t *testing.T) {
	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatalf("Open(%q): unexpected error: %s", dbPath, err)
	}
	defer db.Close()
	if want, got := false, db.NeedsUpgrade(); want != got {
		t.Errorf("db.NeedsUpgrade(): want %b got %b", want, got)
	}
}

func TestUpgrade(t *testing.T) {
	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatalf("Open(%q): unexpected error: %s", dbPath, err)
	}
	defer db.Close()
	if want, got := ErrReadOnlyDB, db.Upgrade(); want != got {
		t.Errorf("db.Upgrade(): want error %q got %q", want, got)
	}

	db, err = Open(dbPath, DBReadWrite)
	if err != nil {
		t.Fatalf("Open(%q): unexpected error: %s", dbPath, err)
	}
	defer db.Close()
	if want, got := error(nil), db.Upgrade(); want != got {
		t.Errorf("db.Upgrade(): want error %q got %q", want, got)
	}
}

func TestAddMessage(t *testing.T) {
	fp, err := filepath.Abs("fixtures/emails/notmuch0202:2,")
	if err != nil {
		t.Fatalf("error getting the absolute path: %s", err)
	}
	nfp, err := filepath.Abs("fixtures/database-v1/new/notmuch0202:2,")
	if err != nil {
		t.Fatalf("error getting the absolute path: %s", err)
	}

	f, err := os.Open(fp)
	if err != nil {
		t.Fatalf("error opening the new email: %s", err)
	}
	defer f.Close()
	nf, err := os.Create(nfp)
	if err != nil {
		t.Fatalf("error creating the new email: %s", err)
	}
	defer nf.Close()
	if _, err := io.Copy(nf, f); err != nil {
		t.Fatalf("error copying the email: %s", err)
	}
	defer os.Remove(nfp)

	db, err := Open(dbPath, DBReadWrite)
	if err != nil {
		t.Fatalf("Open(%q): unexpected error: %s", dbPath, err)
	}
	defer db.Close()
	msg, err := db.AddMessage(nfp)
	if err != nil {
		t.Fatalf("AddMessage(%q): got error: %s", nfp, err)
	}
	defer db.RemoveMessage(nfp)
	if msg == nil {
		t.Errorf("expecting msg to not be nil")
	}
}

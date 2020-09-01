package porcelain

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

const gitOutput string = `
# branch.oid 51c9c58e2175b768137c1e38865f394c76a7d49d
# branch.head master
# branch.upstream origin/master
# branch.ab +1 -10
1 .M N... 100644 100644 100644 3e2ceb914cf9be46bf235432781840f4145363fd 3e2ceb914cf9be46bf235432781840f4145363fd Gopkg.lock
1 .M N... 100644 100644 100644 cecb683e6e626bcba909ddd36d3357d49f0cfd09 cecb683e6e626bcba909ddd36d3357d49f0cfd09 Gopkg.toml
1 .M N... 100644 100644 100644 aea984b7df090ce3a5826a854f3e5364cd8f2ccd aea984b7df090ce3a5826a854f3e5364cd8f2ccd porcelain.go
1 .D N... 100644 100644 000000 6d9532ba55b84ec4faf214f9cdb9ce70ec8f4f5b 6d9532ba55b84ec4faf214f9cdb9ce70ec8f4f5b porcelain_test.go
2 R. N... 100644 100644 100644 44d0a25072ee3706a8015bef72bdd2c4ab6da76d 44d0a25072ee3706a8015bef72bdd2c4ab6da76d R100 hm.rb     hw.rb
u UU N... 100644 100644 100644 100644 ac51efdc3df4f4fd328d1a02ad05331d8e2c9111 36c06c8752c78d2aff89571132f3bf7841a7b5c3 e85207e04dfdd5eb0a1e9febbc67fd837c44a1cd hw.rb
? _porcelain_test.go
? git.go
? git_test.go
? goreleaser.yml
? vendor/
`

var expectedPorcInfo = PorcInfo{
	Branch:    "master",
	Commit:    "51c9c58e2175b768137c1e38865f394c76a7d49d",
	Remote:    "",
	Upstream:  "origin/master",
	Ahead:     1,
	Behind:    10,
	Untracked: 5,
	Unmerged:  1,
	Unstaged: GitArea{
		Modified: 3,
		Added:    0,
		Deleted:  1,
		Renamed:  0,
		Copied:   0,
	},
	Staged: GitArea{
		Modified: 0,
		Added:    0,
		Deleted:  0,
		Renamed:  1,
		Copied:   0,
	},
}

func TestParsePorcInfo(t *testing.T) {
	var pi = new(PorcInfo)
	if err := pi.ParsePorcInfo(strings.NewReader(gitOutput)); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(&expectedPorcInfo, pi) {
		t.Logf("%#+v\n", pi)
		t.FailNow()
	}
}

func TestParseGetGitOutput(t *testing.T) {
	cwd, _ := os.Getwd()
	gitOut, err := GetGitOutput(cwd)
	if err != nil {
		log.Printf("error: %s", err)
		if err == ErrNotAGitRepo {
			os.Exit(0)
		}
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}

	var pi = new(PorcInfo)
	if err := pi.ParsePorcInfo(gitOut); err != nil {
		t.Fatal(err)
	}
	// Print parse results
	pi.Fmt()
}

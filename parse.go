package porcelain

import (
	"bufio"
	"io"
	"log"
	"strconv"
	"strings"
)

func consumeNext(s *bufio.Scanner) string {
	if s.Scan() {
		return s.Text()
	}
	return ""
}

func (pi *PorcInfo) ParsePorcInfo(r io.Reader) error {
	log.Println("parsing git output")

	var err error
	var s = bufio.NewScanner(r)

	for s.Scan() {
		if len(s.Text()) < 1 {
			continue
		}

		pi.ParseLine(s.Text())
	}

	return err
}

func (pi *PorcInfo) ParseLine(line string) error {
	s := bufio.NewScanner(strings.NewReader(line))
	// switch to a word based scanner
	s.Split(bufio.ScanWords)

	for s.Scan() {
		switch s.Text() {
		case "#":
			pi.parseBranchInfo(s)
		case "1":
			pi.parseTrackedFile(s)
		case "2":
			pi.parseRenamedFile(s)
		case "u":
			pi.Unmerged++
		case "?":
			pi.Untracked++
		}
	}
	return nil
}

func (pi *PorcInfo) parseBranchInfo(s *bufio.Scanner) (err error) {
	// uses the word based scanner from ParseLine
	for s.Scan() {
		switch s.Text() {
		case "branch.oid":
			pi.Commit = consumeNext(s)
		case "branch.head":
			pi.Branch = consumeNext(s)
		case "branch.upstream":
			pi.Upstream = consumeNext(s)
		case "branch.ab":
			err = pi.parseAheadBehind(s)
		}
	}
	return err
}

func (pi *PorcInfo) parseAheadBehind(s *bufio.Scanner) error {
	// uses the word based scanner from ParseLine
	for s.Scan() {
		i, err := strconv.Atoi(s.Text()[1:])
		if err != nil {
			return err
		}

		switch s.Text()[:1] {
		case "+":
			pi.Ahead = i
		case "-":
			pi.Behind = i
		}
	}
	return nil
}

// parseTrackedFile parses the porcelain v2 output for tracked entries
// doc: https://git-scm.com/docs/git-status#_changed_tracked_entries
//
func (pi *PorcInfo) parseTrackedFile(s *bufio.Scanner) error {
	// uses the word based scanner from ParseLine
	var index int
	for s.Scan() {
		switch index {
		case 0: // xy
			pi.parseXY(s.Text())
		default:
			continue
			// case 1: // sub
			// 	if s.Text() != "N..." {
			// 		log.Println("is submodule!!!")
			// 	}
			// case 2: // mH - octal file mode in HEAD
			// 	log.Println(index, s.Text())
			// case 3: // mI - octal file mode in index
			// 	log.Println(index, s.Text())
			// case 4: // mW - octal file mode in worktree
			// 	log.Println(index, s.Text())
			// case 5: // hH - object name in HEAD
			// 	log.Println(index, s.Text())
			// case 6: // hI - object name in index
			// 	log.Println(index, s.Text())
			// case 7: // path
			// 	log.Println(index, s.Text())
		}
		index++
	}
	return nil
}

func (pi *PorcInfo) parseXY(xy string) error {
	switch xy[:1] { // parse staged
	case "M":
		pi.Staged.Modified++
	case "A":
		pi.Staged.Added++
	case "D":
		pi.Staged.Deleted++
	case "R":
		pi.Staged.Renamed++
	case "C":
		pi.Staged.Copied++
	}

	switch xy[1:] { // parse unstaged
	case "M":
		pi.Unstaged.Modified++
	case "A":
		pi.Unstaged.Added++
	case "D":
		pi.Unstaged.Deleted++
	case "R":
		pi.Unstaged.Renamed++
	case "C":
		pi.Unstaged.Copied++
	}
	return nil
}

func (pi *PorcInfo) parseRenamedFile(s *bufio.Scanner) error {
	return pi.parseTrackedFile(s)
}

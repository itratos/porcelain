package porcelain

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/robertgzr/color"
)

var (
	Commit  string = "invalid"
	Version string = "invalid"
	Date    string = "invalid"
)

var (
	Cwd         string
	NoColorFlag bool
	FmtFlag     bool
	DebugFlag   bool
	ZshFmtFlag  bool
	BashFmtFlag bool
	TmuxFmtFlag bool
	VersionFlag bool
)

type GitArea struct {
	Modified int
	Added    int
	Deleted  int
	Renamed  int
	Copied   int
}

func (a *GitArea) hasChanged() bool {
	var changed bool
	if a.Added != 0 {
		changed = true
	}
	if a.Deleted != 0 {
		changed = true
	}
	if a.Modified != 0 {
		changed = true
	}
	if a.Copied != 0 {
		changed = true
	}
	if a.Renamed != 0 {
		changed = true
	}
	return changed
}

type PorcInfo struct {
	WorkingDir string

	Branch   string
	Commit   string
	Remote   string
	Upstream string
	Ahead    int
	Behind   int

	Untracked int
	Unmerged  int

	Unstaged GitArea
	Staged   GitArea
}

func (pi *PorcInfo) hasUnmerged() bool {
	if pi.Unmerged > 0 {
		return true
	}
	gitDir, err := PathToGitDir(Cwd)
	if err != nil {
		log.Printf("error calling PathToGitDir: %s", err)
		return false
	}
	// TODO figure out if output of MERGE_HEAD can be useful
	if _, err := ioutil.ReadFile(path.Join(gitDir, "MERGE_HEAD")); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Printf("error reading MERGE_HEAD: %s", err)
		return false
	} else {
		return true
	}
}
func (pi *PorcInfo) hasModified() bool {
	return pi.Unstaged.hasChanged()
}
func (pi *PorcInfo) isDirty() bool {
	return pi.Staged.hasChanged()
}

func (pi *PorcInfo) Debug() string {
	return fmt.Sprintf("%#+v", pi)
}

// Fmt formats the output for the shell
// TODO should be configurable by the user
//
func (pi *PorcInfo) Fmt() string {
	log.Printf("formatting output: %s", pi.Debug())

	var (
		branchGlyph   string = ""
		modifiedGlyph string = "Δ"
		// deletedGlyph   string = "＊"
		dirtyGlyph     string = "✘"
		cleanGlyph     string = "✔"
		untrackedGlyph string = "?"
		unmergedGlyph  string = "‼"
		aheadArrow     string = "↑"
		behindArrow    string = "↓"
	)

	if NoColorFlag {
		color.NoColor = true
	} else {
		color.NoColor = false
		color.EscapeBashPrompt = BashFmtFlag
		color.EscapeZshPrompt = ZshFmtFlag
		color.TmuxMode = TmuxFmtFlag
	}
	branchFmt := color.New(color.FgBlue).SprintFunc()
	commitFmt := color.New(color.FgGreen, color.Italic).SprintFunc()

	aheadFmt := color.New(color.Faint, color.BgYellow, color.FgBlack).SprintFunc()
	behindFmt := color.New(color.Faint, color.BgRed, color.FgWhite).SprintFunc()

	modifiedFmt := color.New(color.FgBlue).SprintFunc()
	// deletedFmt := color.New(color.FgYellow).SprintFunc()
	dirtyFmt := color.New(color.FgRed).SprintFunc()
	cleanFmt := color.New(color.FgGreen).SprintFunc()

	untrackedFmt := color.New(color.Faint).SprintFunc()
	unmergedFmt := color.New(color.FgCyan).SprintFunc()

	return fmt.Sprintf("%s %s@%s %s %s %s",
		branchGlyph,
		branchFmt(pi.Branch),
		func() string {
			if pi.Commit == "(initial)" {
				return commitFmt(pi.Commit)
			}
			return commitFmt(pi.Commit[:7])
		}(),
		func() string {
			var buf bytes.Buffer
			if pi.Ahead > 0 {
				buf.WriteString(aheadFmt(" ", aheadArrow, pi.Ahead, " "))
			}
			if pi.Behind > 0 {
				buf.WriteString(behindFmt(" ", behindArrow, pi.Behind, " "))
			}
			return buf.String()
		}(),
		func() string {
			var buf bytes.Buffer
			if pi.Untracked > 0 {
				buf.WriteString(untrackedFmt(untrackedGlyph))
			} else {
				buf.WriteRune(' ')
			}
			if pi.hasUnmerged() {
				buf.WriteString(unmergedFmt(unmergedGlyph))
			} else {
				buf.WriteRune(' ')
			}
			if pi.hasModified() {
				buf.WriteString(modifiedFmt(modifiedGlyph))
			} else {
				buf.WriteRune(' ')
			}
			// TODO star glyph
			return buf.String()
		}(),
		// dirty/clean
		func() string {
			if pi.isDirty() {
				return dirtyFmt(dirtyGlyph)
			} else {
				return cleanFmt(cleanGlyph)
			}
		}(),
	)
}

func Run() *PorcInfo {
	gitOut, err := GetGitOutput(Cwd)
	if err != nil {
		log.Printf("error: %s", err)
		if err == ErrNotAGitRepo {
			os.Exit(0)
		}
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}

	var porcInfo = new(PorcInfo)
	porcInfo.WorkingDir = Cwd

	if err := porcInfo.ParsePorcInfo(gitOut); err != nil {
		log.Printf("error: %s", err)
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}

	return porcInfo
}

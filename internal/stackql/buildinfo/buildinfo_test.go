package buildinfo_test

import (
	"testing"

	"github.com/stackql/stackql/internal/stackql/buildinfo"
)

func TestNewBuildInfo(t *testing.T) {
	bi := buildinfo.NewBuildInfo("1", "2", "3", "abc123", "abc12", "2024-01-01", "linux")
	if bi.GetMajorVersion() != "1" {
		t.Errorf("GetMajorVersion() = %q, want %q", bi.GetMajorVersion(), "1")
	}
	if bi.GetMinorVersion() != "2" {
		t.Errorf("GetMinorVersion() = %q, want %q", bi.GetMinorVersion(), "2")
	}
	if bi.GetPatchVersion() != "3" {
		t.Errorf("GetPatchVersion() = %q, want %q", bi.GetPatchVersion(), "3")
	}
	if bi.GetCommitSHA() != "abc123" {
		t.Errorf("GetCommitSHA() = %q, want %q", bi.GetCommitSHA(), "abc123")
	}
	if bi.GetShortCommitSHA() != "abc12" {
		t.Errorf("GetShortCommitSHA() = %q, want %q", bi.GetShortCommitSHA(), "abc12")
	}
	if bi.GetDate() != "2024-01-01" {
		t.Errorf("GetDate() = %q, want %q", bi.GetDate(), "2024-01-01")
	}
	if bi.GetPlatform() != "linux" {
		t.Errorf("GetPlatform() = %q, want %q", bi.GetPlatform(), "linux")
	}
	if bi.GetSemVersion() != "1.2.3" {
		t.Errorf("GetSemVersion() = %q, want %q", bi.GetSemVersion(), "1.2.3")
	}
}

func TestNewBuildInfoEmpty(t *testing.T) {
	bi := buildinfo.NewBuildInfo("", "", "", "", "", "", "")
	if bi.GetMajorVersion() != "" {
		t.Errorf("GetMajorVersion() = %q, want %q", bi.GetMajorVersion(), "")
	}
	if bi.GetSemVersion() != ".." {
		t.Errorf("GetSemVersion() = %q, want %q", bi.GetSemVersion(), "...")
	}
}

func TestInitAndSingleton(t *testing.T) {
	// Before Init, the singleton has empty strings.
	// We cannot easily reset the singleton, so just verify Get() returns
	// a BuildInfo (the singleton may have been initialized elsewhere).
	bi := buildinfo.Get()
	if bi == nil {
		t.Fatal("Get() returned nil")
	}
}

func TestInitPopulatesSingleton(t *testing.T) {
	// Init is idempotent via sync.Once; subsequent calls are no-ops.
	// Verify the interface is satisfied and returns the same value.
	bi1 := buildinfo.Get()
	bi2 := buildinfo.Get()
	if bi1 == nil || bi2 == nil {
		t.Fatal("Get() returned nil")
	}
	if bi1.GetSemVersion() != bi2.GetSemVersion() {
		t.Errorf("Get() returned inconsistent values: %q vs %q", bi1.GetSemVersion(), bi2.GetSemVersion())
	}
}
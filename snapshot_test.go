package virtualbox

import (
	"errors"
	"testing"
)

func TestSnapshotList(t *testing.T) {
	Setup(t)
	var (
		snapshotListOut string
	)
	nameE, idE := "baseline", "0d50e9a7-0a6b-4357-bd71-73b30ca9add8"
	if ManageMock != nil {
		snapshotListOut = ReadTestData("vboxmanage-snapshot-list.out")
		ManageMock.EXPECT().run("snapshot", VM, "list").Return(snapshotListOut, "", nil)
	}
	b, _, err := Manage().run("snapshot", VM, "list")
	if err != nil {
		t.Fatal(err)
	}

	snapshots, err := parseSnapshots(b)
	if err != nil {
		t.Fatal(err)
	}

	if len(snapshots) != 1 {
		t.Fatal(errors.New("snapshot count mismatch"))
	}

	snapshot := snapshots[0]
	if snapshot.Name != nameE || snapshot.ID != idE {
		t.Fatal(errors.New("snapshot mismatch"))
	}

	for _, s := range snapshots {
		t.Logf("%+v", s)
	}

	Teardown()
}

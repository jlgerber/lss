package lss

import (
	//"fmt"
	"testing"
)

func TestDirItem_Prefix(t *testing.T) {
	tests := []string{
		"*RD100_main_robot.0001.exr",
		".RD100_main_robot_v2.001.mb",
		".%RD100_main.100.robot_v2.1.exr",
		"RD100_main.01",
	}
	for _, testStr := range tests {
		retval := NewDirItemFromString(testStr)
		if len(retval.Prefix) <= 0 {
			t.Error("Prefix is screwed:", retval.Prefix)
		}
	}
}

func TestDirItem_Number(t *testing.T) {
	tests := []string{
		"RD100_main_robot.0001.exr",
		"RD100_main_robot_v2.001.mb",
		"RD100_main.100.robot_v2.1.exr",
		"RD100_main.01",
	}
	for _, testStr := range tests {
		retval := NewDirItemFromString(testStr)
		if retval.Number != 1 {
			t.Error("Wrong number identified:", retval.Number)
		}
	}
}

func TestDirItem_Padding(t *testing.T) {
	tests := []string{
		"RD100_main_robot.0001.exr",
		"RD100_main_robot_v2.1111.mb",
		"RD100_main.100.robot_v2.1000.exr",
		"RD100_main.1101",
	}
	for _, testStr := range tests {
		retval := NewDirItemFromString(testStr)
		if retval.Padding != 4 {
			t.Error("Wrong padding identified:", retval.Padding)
		}
	}
}

func TestDirItem_Padded(t *testing.T) {
	padded := NewDirItemFromString("foo.0001")
	notpadded := NewDirItemFromString("foo.0")
	either := NewDirItemFromString("foo.100")
	if padded.Padded() != PADDED_YES {
		t.Error("DirItem.Padded() did not recognize padded item")
	}
	if notpadded.Padded() != PADDED_NO {
		t.Error("DirItem.Padded() call did not recognize unpadded item")
	}
	if either.Padded() != PADDED_EITHER {
		t.Error("DirItem.Padded() did not regocnize ambiguous padding")
	}
}
func TestDirItem_Match(t *testing.T) {
	lhs := "RD100_main_robot.0001.exr"
	rhs := "RD100_main_robot.0002.exr"
	if !DirItemsStrMatch(lhs, rhs) {
		t.Error("Provided strings should match", lhs, rhs)
	}
}

func TestDirItem_MisMatchExt(t *testing.T) {
	lhs := "RD100_main_robot.0001.exr"
	rhs := "RD100_main_robot.0002.md"
	if DirItemsStrMatch(lhs, rhs) {
		t.Error("Provided strings should not match", lhs, rhs)
	}
}

func TestDirItem_MisMatchPadding(t *testing.T) {
	lhs := "RD100_main_robot.01.md"
	rhs := "RD100_main_robot.0002.md"
	if DirItemsStrMatch(lhs, rhs) {
		t.Error("Provided strings should not match", lhs, rhs)
	}
}

func TestDirItem_MisMatchCapitalization(t *testing.T) {
	lhs := "rd100_main_robot.0001.md"
	rhs := "RD100_main_robot.0002.md"
	if DirItemsStrMatch(lhs, rhs) {
		t.Error("Provided strings should not match", lhs, rhs)
	}
}

func TestDirItem_NewDirItemFromString(t *testing.T) {
	il := NewDirItemListFromSlice([]string{
		"foo.0001.mb",
		"foo.0002.mb",
	})
	if len(il) != 2 {
		t.Error("impropper construction. should produce 2 items. we have:",
			len(il), "items:", il)
	}
}

// Tests that we convert the padding/number and extension correctly, among other
// things
func TestDirItem_StringConversion(t *testing.T) {
	tests := []string{
		"foo.1.mb",
		"foo.0100.mb",
		"foo.1",
		"foo.001",
	}
	// loop over the strings, convert to DirItem, and verify that
	// DirItem.String() == string
	for _, val := range tests {
		tdi := NewDirItemFromString(val)
		if tdi.String() != val {
			t.Error("failed to convert", val, "to", val, ".Results:",
				tdi.String())
		}
	}
}

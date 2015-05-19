package lss

/*
dirItem provdes a struct modeling a directory item tailored for lss's purposes. A DirItem has a prefix,
and, optionally, a padded number and an extension. It also provides a number of functions and methods
including altername constructors, and special comparison functions.

The DirItem helps trnasform:
foo_bar.0001.mb
foo_bar.0002.mb
foo_bar.0003.mb

Into:
foo_bar.%04d.mb  1-3

and this:
foo_bar.0001.mb
foo_bar.0003.mb
foo_bar.0004.mb

Into this:
foo_bar.%04d.mb 1,3-4
*/

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

//import lg "lsscmd/lsslog"

//var Logr = lg.NewDefaultLogger(lg.DEBUG)

/*
There are cases where lss cannot solve the problem
foo.099.exr
foo.99.exr

foo.100.exr which does it belong to?

*/

//-------------------------
// Type Padding
//-------------------------
type Padding int

const (
	PADDED_NO     Padding = iota
	PADDED_YES            // yes we are padded (start with 0)
	PADDED_EITHER         // the padding count equals the string rep of the number (eg 3,100)
	PADDED_FAIL           // if padding test fails
)

//---------------------------
// Type DirItem
//---------------------------

// DirItem
//     Represents the components of an item in a directory for the purposes of calculating ranges
//
// Vars:
//     Prefix  string   - The name of the item up to the '.' prior to the range number
//     Number  int      - The range's number, if representing a range item. Otherwise -1.
//     Padding int      - The size of the padding for the range. -1 if the item is not a range item.
//     Extension string - The file extension, if any, or ""
//     Processed bool   - Book keeping - whether the DirItem has been processed or not
type DirItem struct {
	Prefix    string
	Number    int
	Padding   int
	Extension string
}

//-----------------------------
// DirItem Constructors
//-----------------------------

// NewDirRangeItem
//     Alternate constructor
//     Given a prefix, number, padding, and extension, return a pointer to a new DirItem
//
// Vars:
//     prefix string    - The name of the item preceding the number.
//     number int       - The number of the range item.
//     padding int      - The padding size of the range item.
//     extension string - The extension name
//
// Returns:
//     *DirItem - a pointer to a DirItem
func NewDirRangeItem(prefix string, number int, padding int, extension string) *DirItem {
	di := DirItem{prefix, number, padding, extension}
	return &di
}

// NewDirItem
//     Alternate Constructor
//     This constructor builds a DirItem initialized appropriately for a non-range directory item.
//     Chiefly, this means initializing Number and Padding to -1.
//
// Args:
//     name string - The name of the DirItem. The other values will be initialized appropriately.
//
// Returns:
//     *DirItem - pointer to a DirItem instance.
func NewDirItem(name string) *DirItem {
	di := new(DirItem)
	di.Number = -1
	di.Padding = -1
	di.Prefix = name
	return di
}

// NewDirItemFromSlice
//     AlternateConstructor
//     Use a string slice to construct a DirItem. The slice is expected to follow this convention:
//     []string{prefix[,number,extension]}
//
// Vars:
//     pieces []string - A string slice whose members should correspond to the components of a DirItem.
//                       That is, they should contain a prefix at minimum, and optionally a number, and extension
//
// Returns:
//     *DirItem - a pointer to a new DirItem instance.
func NewDirItemFromSlice(pieces []string) *DirItem {
	di := new(DirItem)
	di.Number = -1
	di.Padding = -1

	di.Prefix = pieces[0]
	if len(pieces) < 2 {
		return di
	}

	tmp, err := strconv.ParseInt(pieces[1], 10, 0)

	if err != nil {
		return di
	}
	di.Number = int(tmp)

	di.Padding = len(pieces[1])
	if len(pieces) < 3 {
		return di
	}
	di.Extension = pieces[2]
	return di
}

// NewDirItemFromString
//     AlternateConstructor
//     Construct a DirITem from a string representation of a directory item.
//
// Vars:
//     item string - The string representation of the directory item.
//
// Returns:
//     *DirItem - Pointer ot a DirItem struct instance.
func NewDirItemFromString(item string) *DirItem {
	results := re.FindAllStringSubmatch(item, -1)
	if len(results) > 0 {
		return NewDirItemFromSlice(results[0][1:])
	} else {
		return NewDirItem(item)
	}
}

//-------------------------
// DirItem Methods
//-------------------------

// Padded Method returns a Padding enum indicating whether we are padded, not padded
// or context dependent ( eg 100 is either 3 padded or not padded depending)
func (di *DirItem) Padded() Padding {

	switch {
	// 0-9
	case di.Padding == 1:
		return PADDED_NO
	// any number >9 that doesn't start with a 0
	case len(strconv.Itoa(di.Number)) == di.Padding:
		return PADDED_EITHER
	// any number starting with a 0, 1 or greater
	default:
		return PADDED_YES
	}
	return PADDED_FAIL
}

// GetPaddedNumber takes paddinginto consideration and
// reconstructs an appropriate number string
// eg if Padding = 4 and Number =1, we return "0001"
func (di *DirItem) GetPaddedNumber() string {
	numstr := strconv.Itoa(di.Number)
	padding := di.Padding - len(numstr)
	if padding > 0 {
		return strings.Repeat("0", padding) + numstr
	}
	return numstr
}

// GetExtension returns a normalized extension, prefixing
// the extension with ".". In the future, we will do the
// opposite, stripping "." off it is the first letter...
func (di *DirItem) GetExtension() string {
	if len(di.Extension) > 0 {
		if di.Extension[0] == '.' {
			return di.Extension
		} else {
			return "." + di.Extension
		}
	}
	return ""
}

// ApproxLen method used to calculate the approximate length of the DirItem for calculating padding
func (di *DirItem) ApproxLen() int {
	if di.Padding > 1 {
		return len(di.Prefix) + len(di.Extension) + 6 // 1 period for ext. 1 for period before num, 4 for %0#d
	} else if di.Padding >= 0 {
		return len(di.Prefix) + len(di.Extension) + 3
	}
	return len(di.Prefix)
}

// GetName
//     *DirItem method which returns a string representation of the struct.
//
// Returns:
//     string - The string representation of the DirItem
func (di *DirItem) String() string {
	// if we don't have a number, we are a non-matching item
	if di.Number < 0 {
		return di.Prefix
	}
	// pad number appropriately
	num := di.GetPaddedNumber()

	// determine suffix
	ext := di.GetExtension() //""

	return fmt.Sprintf("%s.%s%s", di.Prefix, num, ext)
}

//-------------------------------
// DirItem Functions
//-------------------------------

// DirItemsMatch
//     Takes two DirItems and attempts to determine whether they "match".
//     By match I mean that they have the same Prefix, the same padding, and the same
//     extension
//
// Vars:
//    lhs *DirItem - A pointer to the first DirItem
//    rhs *DirItem - A pointer to the second DirItem
//
// Returns:
//     bool - indicating whether the two DirItems match or not.
func DirItemsMatch(lhs *DirItem, rhs *DirItem) bool {
	if lhs.Prefix == rhs.Prefix &&
		lhs.Padding == rhs.Padding &&
		lhs.Extension == rhs.Extension {
		return true
	}
	return false
}

// DirItemsStrMatch
//     Give two string representations of DirItems, determine if they "match",
//     using the logic described in DirItemsMatch above.
//
// Vars:
//     lhs string - The first string representation of a DirItem.
//     rhs string - The second string representation of a DirItem.
//
// Returns:
//     bool - A boolean indicating whether or not the DirItems match
func DirItemsStrMatch(lhs string, rhs string) bool {
	lhsDi := NewDirItemFromString(lhs)
	rhsDi := NewDirItemFromString(rhs)
	return DirItemsMatch(lhsDi, rhsDi)
}

// DirItemsContiguous
//     given two DirItemsdetermine
//     if they are contiguous. We define contiguous in this case to be
//     abs(rhs.Number - lhs.Number) = 1
//
// Vars:
//     lhs *DirItem - Pointer to first DirItem.
//     rhs *DirItem - Pointer to second DirItem.
//
// Returns:
//     bool - Whether the DirItems are contiguous or not.
func DirItemsContiguous(lhs *DirItem, rhs *DirItem) bool {
	// these have to be the same other than the count
	if lhs.Prefix != rhs.Prefix ||
		lhs.Extension != rhs.Extension {
		// if the left hand side padding is not eqal to the right hand side padding
		// return false unless either of the following is true:
		// the left hand side padding is 1 and the leading number of the right hand side is 0
		// or
		// the right hand side padding is 1 and the leading number of the left hand side is 0
		if lhs.Padding != rhs.Padding &&
			!((lhs.Padding == 1 && string(strconv.Itoa(rhs.Number)[0]) != "0") ||
				(rhs.Padding == 1 && string(strconv.Itoa(lhs.Number)[0]) != "0")) {
			return false
		}
	}

	distance := lhs.Number - rhs.Number
	if distance < 0 {
		distance *= -1
	}
	if distance == 1 {
		return true
	}
	return false
}

//-----------------------------------------
// Private Utility Functions & Variables
//-----------------------------------------

// compileRegex
//     Precompile the regular expression matching prefix, number, and extension
//
// Returns:
//     *regexp.Regexp - pointer to the compiled regexp object
func compileRegex() *regexp.Regexp {
	var re, err = regexp.Compile("(.*)\\.([0-9]+)(\\..*){0,1}")

	if err != nil {
		fmt.Printf("ERROR compiling regex")
		os.Exit(1)
	}
	return re
}

var re = compileRegex()

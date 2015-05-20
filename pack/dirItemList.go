package lss

import (
	"fmt"
	"github.com/xlab/handysort"
	"strconv"
	"strings"
)

//-------------------------
// Type DirItemList
//-------------------------

type DirItemList []DirItem

//-------------------------
// DirItemList constructors
//-------------------------

func NewDirItemListFromSlice(items []string) DirItemList {
	di := make(DirItemList, 0)
	for _, val := range items {
		di = append(di, *NewDirItemFromString(val))
	}
	return di

}

func NewSliceFromDirItemList(dil DirItemList) []string {
	ss := []string{}
	for _, val := range dil {
		ss = append(ss, val.String())
	}
	return ss
}

//--------------------------
// DirITemList methods
//--------------------------

func (dil *DirItemList) Len() int {
	return len(*dil)
}

func (dil DirItemList) Swap(i, j int) {
	dil[i], dil[j] = dil[j], dil[i]
}

func (dil DirItemList) Less(i, j int) bool {
	switch {
	case dil[i].Prefix == dil[j].Prefix:
		if dil[i].Number == dil[j].Number {
			return dil[i].Padding < dil[j].Padding
		} else {
			return dil[i].Number < dil[j].Number
		}
	default:
		return handysort.StringLess(dil[i].Prefix, dil[j].Prefix)

	}
}

func (dil DirItemList) String() string {
	ret := ""
	for _, val := range dil {
		ret += val.String() + "\n"
	}
	return ret
}

// SortDirItemList - This function is designed to take a sorted
// DirItemList and separate it into two lists - one for padded items
// and one for non-padded items. The trick is that the non-padded items
// are context dependent
func SortDirItemList(list DirItemList) (DirItemList, DirItemList, int) {
	maxlen := 0 // max length
	clen := 0   // current length
	padded := make(DirItemList, 0)
	unpadded := make(DirItemList, 0)

	// here we go
	for _, diritem := range list {

		clen = diritem.ApproxLen()
		if clen > maxlen {
			maxlen = clen
		}

		padding := diritem.Padded()
		switch padding {
		case PADDED_NO:
			unpadded = append(unpadded, diritem)
		case PADDED_YES:
			padded = append(padded, diritem)
		case PADDED_EITHER:

			if pl := len(padded); pl > 0 && padded[pl-1].Padding == diritem.Padding {
				padded = append(padded, diritem)
			} else if upl := len(unpadded); upl > 0 {
				unpadded = append(unpadded, diritem)
			} else {
				padded = append(padded, diritem)
			}
		}
	}

	// debuging nonsense
	if 1 == 0 {
		println(padded.String())
		println("$$$$")
		println(unpadded.String())
		println("####")
	}

	return padded, unpadded, maxlen
}

func same(di1 *DirItem, di2 *DirItem) bool {
	// if the prefixes and extensions are the same
	if di1.Prefix == di2.Prefix &&
		di1.Extension == di2.Extension {
		// if the padding is the same or
		// if one of the items is unpadded and the other's
		// leading number is not a zero ( eg 9 & 10 )
		if di1.Padding == di2.Padding ||
			(di1.Padding == 1 && string(strconv.Itoa(di2.Number)[0]) != "0") ||
			(di2.Padding == 1 && string(strconv.Itoa(di1.Number)[0]) != "0") {
			return true
		}
	}
	return false
}

// DivideByType filters a DirItemList into multiple DirItemLists,
// each one representing a single range item. Each unique
// DirItemList is returned to the channel provided as the return type.
func DivideByType(list DirItemList) chan DirItemList {
	// create empty DirItemSlice
	ch := make(chan DirItemList)
	go func() {
		sz := len(list)
		// if someone was a jerk and passed in an empty list
		if sz == 0 {
			close(ch)
			return
		}

		// grab the first one
		outslice := NewDirItemListFromSlice([]string{})
		outslice = append(outslice, list[0])
		cnt := 1

		for {
			if cnt == len(list) {
				break
			}

			if same(&list[cnt], &list[cnt-1]) {
				outslice = append(outslice, list[cnt])
			} else {
				ch <- outslice
				outslice = NewDirItemListFromSlice([]string{})
				outslice = append(outslice, list[cnt])
			}
			cnt++
		}

		if len(outslice) > 0 {
			ch <- outslice
		}

		close(ch)
	}()

	return ch
}

// BuildRangeStringPrefix
//     Given a pointer to a DirItem, build an appropriate range string and return it.
//     Generally, the range string will be in the form of:
//     DirItem.Prefix + '.' + "%0" + Diritem.Padding + "d" + DirItem.Extension
//
// Args:
//     item *dirItem.DirItem - pointer to DirItem instance, from which we will build the string.
//
// Returns:
//     string - A string with range formatting ( %04d)
func BuildRangeStringPrefix(item *DirItem) string {
	switch {
	case item.Padding < 0:
		return item.Prefix // copy
	case len(item.Extension) == 0:
		if item.Padding <= 1 {
			return item.Prefix + ".%d" //+ strconv.Itoa(item.Padding)
		}
		return item.Prefix + ".%0" + strconv.Itoa(item.Padding) + "d"
	default:
		retstr := item.Prefix
		if item.Padding <= 1 {
			retstr += ".%d" // + strconv.Itoa(item.Padding)
		} else {
			retstr += ".%0" + strconv.Itoa(item.Padding) + "d"
		}

		if string(item.Extension[0]) == "." {
			retstr += item.Extension
			return retstr
		} else {
			retstr += "." + item.Extension
			return retstr
		}
	}
}

// PadInt - generate a string from the input integer, padding it
// with "0s" appropriately.
func PadInt(x, padding int) string {
	xPadding := NumDigits(x)
	sz := padding - xPadding

	if sz <= 0 {
		return strconv.Itoa(x)
	}

	return fmt.Sprintf("%d%s", x, strings.Repeat(" ", sz))
}

// PadToSize - given a source string, a target size, and an indication
// of where the padding comes ( prefix or not), generate a space padded
// version of the source string.
// egs
// PadToSize("foo",5,true) return "  foo"
// PadToSize("foo",5,false) return "foo  "
func PadToSize(source string, toSize int, prefix bool) string {
	sz := len(source)
	padding := toSize - sz
	if padding <= 0 {
		return source
	}

	// do i stick the padding in front?
	if prefix {
		return fmt.Sprintf("%s%s", strings.Repeat(" ", padding), source)
	}
	return fmt.Sprintf("%s%s", source, strings.Repeat(" ", padding))
}

// BuildRangeSlice takes an hemogenous DirItemList and returns
// a channel of type string. Each string is in condensed range form
// ie foo.%04d.mb  1-4,10,100-122
func BuildRangeString(list DirItemList, rangePadding int) string {
	padding := 5 // arbitrary true

	if len(list) == 1 {
		return fmt.Sprintf("%s %s", PadInt(1, padding), list[0].String())
	}

	rangestr := PadToSize(BuildRangeStringPrefix(&list[0]), rangePadding, false)
	rangestr = PadInt(len(list), padding) + " " + rangestr + "    "

	// is there a range at all?
	if list[0].Padding == -1 && list[0].Number == -1 {
		return rangestr
	}

	last := len(list) - 1
	lastcontiguous := -1
	for i, diritem := range list {
		// special case the first time through
		if i == 0 {
			rangestr += strconv.Itoa(diritem.Number)
			continue
		}
		// special case the last time through
		// TODO: take care of len(list) == 1
		if i == last {
			if DirItemsContiguous(&list[i-1], &diritem) {
				rangestr += "-" + strconv.Itoa(diritem.Number)
				break
			}
		}

		if DirItemsContiguous(&list[i-1], &diritem) {
			lastcontiguous = diritem.Number
		} else {
			if lastcontiguous > 0 {
				rangestr += "-" +
					strconv.Itoa(lastcontiguous)
				lastcontiguous = 0
			}
			rangestr += "," +
				strconv.Itoa(diritem.Number)
		}

	}
	return rangestr
}

// RangeStringsFromSortedItemList takes a sorted DirItemList and returns a channel on which sorted strings are returned.
// Your first question is probably why? Well, it does some internal kungfu to take that item list and reduce it
// into its final output.
func RangeStringsFromSortedItemList(itemList DirItemList) chan string {

	// sort DirItems into padded and nonPadded
	padded, unpadded, maxlen := SortDirItemList(itemList)
	paddedChan := DivideByType(padded)
	unpaddedChan := DivideByType(unpadded)

	ch := make(chan string)

	go func() {
		prBuff := make([]string, 0)  // padded Range buffer
		uprBuff := make([]string, 0) // unpadded Range buffer

		for {

			paddedRange, prOk := <-paddedChan
			unpaddedRange, uprOk := <-unpaddedChan

			// if we fetched data from the channel
			if prOk {
				prBuff = append(prBuff, BuildRangeString(paddedRange, maxlen))
			}
			// if we fetched data from the channel
			if uprOk {
				uprBuff = append(uprBuff, BuildRangeString(unpaddedRange, maxlen))
			}
			// if both channels are empty lets beat it
			if len(prBuff) == 0 && len(uprBuff) == 0 {
				break
			}
			// ok, we move on.
			switch {
			case len(prBuff) > 0 && len(uprBuff) > 0:
				// buffer entries are in the form
				// count name range eg
				// 2 foo.%04d.exr  4-5
				// we use the fields command to split the string up into a slice
				// and take the middle index (1)
				if strings.Fields(prBuff[0])[1] < strings.Fields(uprBuff[0])[1] {
					ch <- prBuff[0]
					prBuff = prBuff[1:]
				} else {
					ch <- uprBuff[0]
					uprBuff = uprBuff[1:]
				}
			case len(prBuff) > 0:
				ch <- prBuff[0]
				prBuff = prBuff[1:]
			case len(uprBuff) > 0:
				ch <- uprBuff[1]
				uprBuff = uprBuff[1:]
			}

		}
		close(ch)
	}()

	return ch
}

func RangesChanFromStringSlice(contents []string) chan string {

	// cast to a Stringlist and call NaturalSort()
	Stringlist(contents).NaturalSort()

	// build a dirItemList from the contents slice
	dil := NewDirItemListFromSlice(contents)

	return RangeStringsFromSortedItemList(dil)
}

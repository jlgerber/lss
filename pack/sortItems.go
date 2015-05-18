package lss

/*
Defines two types:
StringChunks - a []string representing chunks of a string split into number and non-number strings.
StringChunkList

*/
import "fmt"
import "strings"
import "sort"
import "strconv"

// lastIndex - assuming strval starts with a member of the contains string - returns the index of the last
// contiguous character which matches one of the characters in the 'contains' string.
//
// Args:
//     strval string   - The string to test against.
//     contains string - The string of characters, one of which we assert strval starts with
//                       and which we are looking for the last contiguous match of in "strval"
//
// Returns:
//     int - the matching index of the last contiguous character from the start of strval which
//           is in the set of characters which form "contains".
//
// We are specifically using this to find the range of a padded or unpadded number sequence whithin
// the input string strval. FYI we guarantee that strval starts with a number prior to calling this
// function
func lastIndex(strval string, contains string) int {
	last := -1
	for i, val := range strval {
		if strings.ContainsAny(string(val), contains) == false {
			return i - 1
		}
		last = i
	}
	return last
}

// chunk
//     Given an input string, determine whether the fist character is a number or not. If it is a number,
//     return a string representing the longest contigous substring from the start which contains only
//     numbers, followed by a substring containing the remaining characters.
//     If the first character of val is NOT a number, return the longest substring from the begining
//     comprised of characters which are not numbers, followed by the remaining substring.
//
// Args:
//     val string - the string value to split.
//
// Returns:
//     chunkval string - The longest substring starting from the first character, which contains either
//                       no numbers, or all numbers, depending upon what the first character is.
//     remainder string - The remaining characters once "chunkval" has been extracted.
func chunk(val string) (chunkval string, remainder string) {
	i := strings.IndexAny(val, "0123456789")
	if i < 0 {
		return val, ""
	}
	if i == 0 {
		i2 := lastIndex(val, "0123456789")
		if i2 >= len(val)-1 {
			return val, ""
		}
		return val[:i2+1], val[i2+1:]
	}
	return val[:i], val[i:]
}

//------------------------------------
// StringChunks Type
//------------------------------------
//     a list of strings
type StringChunks []string

// String
//     Method to convert StringChunks into a string, joining all of the slices into a contiguous string.
//
// Returns:
//     string - Eg given foo defined as StringChunks{"foo.","0001"}, foo.String() => foo.0001
func (s StringChunks) String() string {
	return strings.Join(s, "")
}

// GenChunks
//     Given a string, break it up and return a slice of strings, using the chunk function to split the
//     input up into substrings consisting of contiguous numbers or non-numbers. (eg ["foo","010",".exr"])
//
// Vars:
//     val string - The input string
//
// Returns:
//     StringChunks - a list of strings containing chunks of the original val, broken up into character and number
//     strings in order.
func GenChunks(val string) StringChunks {
	retv := []string{}
	var chnk, remain string
	remain = val
	for {
		chnk, remain = chunk(remain)
		retv = append(retv, chnk)
		if remain == "" {
			break
		}
	}
	return retv
}

//------------------------------------
//  StringChunksList type
//------------------------------------
//
//A slice of StringChunks - effectively [][]string
type StringChunksList []StringChunks

// NewStringChunksList
//     New up a StringChunksList eh....
func NewStringChunksList(inStrings []string) StringChunksList {
	chlist := StringChunksList{}
	for _, val := range inStrings {
		retv := GenChunks(val)
		chlist = append(chlist, retv)
	}
	return chlist

}

// ChunkStringsToChan generates StringChunks for each input string
// in the inStrings slice, and sends them to a channel via an anonymous
// go routine, before closing the channel.
//
// The function returns a StringChunks chan.
func ChunkStringsToChan(inStrings []string) chan StringChunks {
	ch := make(chan StringChunks)
	go func() {
		for _, val := range inStrings {
			ch <- GenChunks(val)
		}
		close(ch)
	}()
	return ch
}

// SliceOfStrings
//     Method of StringChunksList which returns a traditional slice of strings
//
// Returns
//     Slice of strings from slice of chunks
func (s StringChunksList) SliceOfStrings() []string {
	ret := []string{}
	for _, val := range s {
		ret = append(ret, val.String())
	}
	return ret
}

// StringChan method returns a channel which we broadcast all
// of the stringchunks too
func (s StringChunksList) StringChan() chan string {
	ch := make(chan string)
	go func() {
		for _, value := range s {
			ch <- value.String()
		}
		close(ch)
	}()
	return ch
}

// Len
//     Method returns the length of the StringChunksList
//
// Returns:
//     int - the length of the StringChunksList ( []string )
func (s StringChunksList) Len() int {
	return len(s)
}

// Swap
//     Swaps the items at index i & j when called on a StringChunksList
//
// Args:
//     i int - one of the indices of the stringlist to swap.
//     j int - The other index of the stringlist to swap.
func (s StringChunksList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func isPadded(s string) bool {
	if string(s[0]) == "0" {
		return true
	}
	return false
}

// Less
//     Method which compares the value at two indices and determines whether or not the
//     first one is less than the second one.
//
// Args:
//     i int - index of first item to compare.
//     j int - index of second item to compare.
//
// Returns:
//     bool - s[i] < s[j]. This method must take into account whether or not s[i] &/or s[j] represents
//     a number or not and perform appropriate conversions so that the comparison is applied in the write
//     context ( ie numbers should be compared as such, and not as strings.
func (s StringChunksList) Less(i, j int) bool {
	i_len := len(s[i])
	j_len := len(s[j])
	min_l := i_len

	if j_len < min_l {
		min_l = j_len
	}

	for c := 0; c < min_l; c++ {
		// try and convert to a number
		val_i, err_i := strconv.Atoi(s[i][c])
		val_j, err_j := strconv.Atoi(s[j][c])
		ic_len := len(s[i][c])
		jc_len := len(s[j][c])
		switch {
		// neither are numbers
		case err_i != nil && err_j != nil:
			if s[i][c] < s[j][c] {
				return true
			}
			if s[i][c] > s[j][c] {
				return false
			}
		// numbers but padding not equal
		case err_i == nil && err_j == nil && ic_len != jc_len:
			return ic_len < jc_len
		// numbers but not equal
		case err_i == nil && err_j == nil && val_i != val_j:
			// we are both numbers
			return val_i < val_j

		/*// equal numbers
		case err_i == nil && err_j == nil && val_i == val_j:
			return ic_len < jc_len*/
		// left is a number
		case err_i == nil && err_j != nil:
			return true
		// right is a number
		case err_j == nil && err_i != nil:
			return false
		}
	}
	// is one stringlist longer than the other? because
	// they are equivalent so far...
	return i_len < j_len
}

// PrintStringChunksList
//     Utility function to print a string list
func PrintStringChunksList(s StringChunksList) {
	for _, val := range s {
		fmt.Println(val)
	}
}

// Sort Method sorts a StringChunksList in place.
func (s StringChunksList) Sort() {
	sort.Sort(s)
}

// NaturalSort
//     Given a list of strings, return a sorted list of strings.
func NaturalSort(items []string) []string {
	itemlist := NewStringChunksList(items)
	sort.Sort(itemlist)
	return itemlist.SliceOfStrings()

}

//------------------------------------
//  Stringlist type
//------------------------------------
type Stringlist []string

func (s Stringlist) NaturalSort() {
	itemlist := NewStringChunksList([]string(s))
	sort.Sort(itemlist)
	for i, val := range itemlist {
		s[i] = val.String()
	}
}

func (s Stringlist) String() string {
	str := ""
	for _, val := range s {
		str += val + "\n"
	}
	return str
}

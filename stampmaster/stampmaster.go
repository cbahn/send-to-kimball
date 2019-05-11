// Utility program to confirm 

package stampmaster

import (
	"time"
	"strings"
	"strconv"
	"regexp"
	"fmt"
	"errors"
	"crypto/sha256"
	"math/bits"
)

/*
	// Challenge form, suggested
	  <difficulty>:<timestamp>:<authentication>:<userdefined>
	*   4:20190509T071655Z:c09799c24f7d50f4:1KomQXRI9  <- valid stamp, where secret=="sendtokimball"
	*   difficulty = number of zeros needed
	*   timestamp = date+time as Format("20060102T150405Z")
	*	authentication = sha256("<difficulty>:<timestamp>:<secret>") truncated to the first 16 hex characters
	*	userdefined = any value that matches /[a-zA-Z0-9]{0,64}/
*/

type Stamp struct {
	difficulty int
	timestamp time.Time
	authentication string
	counter string
}

const (
	TIMESTAMP_LAYOUT = "20060102T150405Z"
)

func ParseStamp(s string) (Stamp, error) {
	var stamp Stamp
	var err error
	stampPatternMatch := regexp.MustCompile(`^[:0-9a-zA-Z]{1,101}$`)
	if !stampPatternMatch.MatchString(s) {
		return stamp, errors.New("input contains invalid characters or length")
	}

	components := strings.Split(s,":")

	if len(components) != 4 {
		return stamp, errors.New("Could not parse stamp, invalid number of delimiters")
	}

	// Validate the difficulty string. It must parse to an int within [1,255]
	stamp.difficulty, err = strconv.Atoi(components[0])
	if err != nil {
		return stamp, errors.New("Could not parse stamp difficulty")
	}
	if 0 >= stamp.difficulty || stamp.difficulty > 256 {
		return stamp, errors.New("Difficulty not in range")
	}

	// Parse timestamp
	stamp.timestamp, err = time.Parse(TIMESTAMP_LAYOUT, components[1])
	if err != nil {
		return stamp, errors.New("could not parse stamp timestamp")
	}

	// Validate that authentication is 16 hex digits
	stamp.authentication = components[2]
	sixteenHexDigitsMatch := regexp.MustCompile(`^[0-9a-f]{16}$`)
	if !sixteenHexDigitsMatch.MatchString(stamp.authentication) {
		return stamp, errors.New("Error: stamp authentication is not 16 hex digits")
	}

	// Validate user-defined counter
	stamp.counter = components[3]
	counterMatch := regexp.MustCompile(`^[0-9a-zA-Z]{0,64}$`)
	if !counterMatch.MatchString(stamp.counter) {
		return stamp, errors.New("Error: stamp counter invalid characetrs or length")
	}

	// Stamp has parsed and is valid
	return stamp, nil
}

func CreateNewStamp(difficulty int, secret string) Stamp {
	var s Stamp
	s.difficulty = difficulty
	s.timestamp = time.Now()
	s.authentication = CreateAuthentication(s, secret)
	return s
}

func (s Stamp) ToString() string {
	return fmt.Sprintf("%d:%s:%s:%s", s.difficulty, s.timestamp.Format(TIMESTAMP_LAYOUT), s.authentication, s.counter)	
}

func CreateAuthentication(s Stamp, secret string) string {
	stringToHash := fmt.Sprintf("%d:%s:%s", s.difficulty, s.timestamp.Format(TIMESTAMP_LAYOUT), secret)
	hash := sha256.Sum256( []byte(stringToHash) )
	return fmt.Sprintf("%x", hash)[:16] // Byte array to hex, trim to 16 characters
}

func (s Stamp) CheckAuthentication(secret string) bool {
	return s.authentication == CreateAuthentication(s, secret)
}

func CountLeadingZeros(bytelist [32]byte) int {
	bitcount := 0
	for _, b := range bytelist {
		lz := bits.LeadingZeros8(b)
		bitcount += lz
		if lz < 8 { 
			// When a byte is encountered with less than 7 leading zeros, we're done counting
			return bitcount
		}
	}
	return 256
}

func (s Stamp) CheckHash() bool {
	hash := sha256.Sum256([]byte( s.ToString() ))
	return CountLeadingZeros(hash) >= s.difficulty
}

/*
func main() {

	// myStamp, err := ParseStamp("4:20190509T071655Z:c09799c24f7d50f4:1KomQouXR35")

	//s := CreateNewStamp(4,"sendtokimball")
	s, _ := ParseStamp("4:20190511T103438Z:e7d4ec2c7ff6692c:")

	fmt.Println(s.ToString())
	fmt.Printf("Auth valid? %t\n", s.CheckAuthentication("sendtokimball"))
	fmt.Printf("Hash valid? %t\n", s.CheckHash())

}
*/
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
 * Stamps are a form of proof that work has been done by the client.
 *
 * A valid stamp is a string of the form "<difficulty>:<timestamp>:<authentication>:<userdefined>"
 * Example of a valid stamp where <salt>=="sendtokimball"    4:20190509T071655Z:c09799c24f7d50f4:1KomQXRI9
 *
 * - Difficulty is the number of zeros needed in the stamp's hash.
 *
 * - Timestamp is the date+time that the stamp computation was started after.
 * This is formatted as time.Format("20060102T150405Z").
 *
 * - Authentication is a code which verifies that the difficulty and timestamp were actually requested
 * from the server. This prevents precomputing stamps or computing stamps of a different difficulty.
 * The authentication is computed as the first 16 hex digits of sha256("<difficulty>:<timestamp>:<salt>").
 *
 * - Salt is a server chosen random value which personalizes the authentication codes to that specific server.
 *
 * - Userdefined can be any string which matches /[a-zA-Z0-9]{0,64}/ . A client must try many values for
 * this field before they happen upon one which hashes to a valid value
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


/* ParseStamp is used to read stamp into its component values
 * This confirms that the stamp is formatted properly, but does not validate it
 */
func ParseStamp(s string) (Stamp, error) {
	var stamp Stamp
	var err error
	stampPatternMatch := regexp.MustCompile(`^[:0-9a-zA-Z]{1,101}$`)
	if !stampPatternMatch.MatchString(s) {
		return stamp, errors.New("input contains invalid characters or length")
	}

	// Split string via delimiter ":"
	components := strings.Split(s,":")
	if len(components) != 4 {
		return stamp, errors.New("Could not parse stamp, invalid number of delimiters")
	}

	// Confirm that the difficulty is an integer in the range [1,255]
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

	// Check that authentication consists of 16 hex digits
	stamp.authentication = components[2]
	sixteenHexDigitsMatch := regexp.MustCompile(`^[0-9a-f]{16}$`)
	if !sixteenHexDigitsMatch.MatchString(stamp.authentication) {
		return stamp, errors.New("Error: stamp authentication is not 16 hex digits")
	}

	// Check that user-defined counter is of correct characters and length
	stamp.counter = components[3]
	counterMatch := regexp.MustCompile(`^[0-9a-zA-Z]{0,64}$`)
	if !counterMatch.MatchString(stamp.counter) {
		return stamp, errors.New("Error: stamp counter invalid characetrs or length")
	}

	// Stamp has parsed and is well formed
	return stamp, nil
}


/* Returns an unfinished stamp with the provided difficulty and current timestamp
 */
func CreateNewStamp(difficulty int, salt string) Stamp {
	var s Stamp
	s.difficulty = difficulty
	s.timestamp = time.Now()
	s.authentication = CreateAuthentication(s, salt)
	return s
}


func (s Stamp) ToString() string {
	return fmt.Sprintf("%d:%s:%s:%s", s.difficulty, s.timestamp.Format(TIMESTAMP_LAYOUT), s.authentication, s.counter)	
}

/* Computes a stamp's correct authentication given the server's salt
 */
func CreateAuthentication(s Stamp, salt string) string {
	stringToHash := fmt.Sprintf("%d:%s:%s", s.difficulty, s.timestamp.Format(TIMESTAMP_LAYOUT), salt)
	hash := sha256.Sum256( []byte(stringToHash) )
	return fmt.Sprintf("%x", hash)[:16] // Byte array to hex, trim to 16 characters
}

/* Confirm if the stamp's authentication is correct for a given salt
 */
func (s Stamp) CheckAuthentication(salt string) bool {
	return s.authentication == CreateAuthentication(s, salt)
}


/* Returns the number of leading zeros in a byte array, starting from smallest index and largest bit
 */
func CountLeadingZeros(bytelist [32]byte) int {
	bitcount := 0
	for _, b := range bytelist {
		lz := bits.LeadingZeros8(b)
		bitcount += lz
		// When a byte is encountered with less than 8 leading zeros, we're done counting
		if lz < 8 { 
			return bitcount
		}
	}
	return 256
}


/* Compares the leading zero count of sha256(stamp) against the difficulty.
 * True means the stamp is valid and should be accepted.
 */
func (s Stamp) CheckHash() bool {
	hash := sha256.Sum256([]byte( s.ToString() ))
	return CountLeadingZeros(hash) >= s.difficulty
}

/* artifacts of testing
func main() {

	//s := CreateNewStamp(4,"sendtokimball")
	s, _ := ParseStamp("16:20190511T193115Z:8812dc7a8a620f6a:q9LUJ0VJ17538")

	fmt.Println(s.ToString())
	fmt.Printf("Auth valid? %t\n", s.CheckAuthentication("sendtokimball"))
	fmt.Printf("Hash valid? %t\n", s.CheckHash())
}
*/
package utils

import (
	"database/sql"
	"fmt"
	"hash/fnv"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func RandomString(numberOnly bool, length int) string {
	choices := "abcdefghijklmnnopqrst0123456789"
	if numberOnly {
		choices = "0123456789"
	}
	out := make([]byte, length)
	for i := range out {
		out[i] = choices[rand.Intn(len(choices))]
	}
	return string(out[:])
}

// RandomID generates a random ID, for which `isTaken` is false.
func RandomID(numberOnly bool, length int, isTaken func(string) bool) string {
	newID := RandomString(numberOnly, length)
	// avoid (unlikely) collisions
	for taken := isTaken(newID); taken; newID = RandomString(numberOnly, length) {
		taken = isTaken(newID)
	}
	return newID
}

// WebsocketError format `err` and send a Control message to `ws`
func WebsocketError(ws *websocket.Conn, err error) {
	message := websocket.FormatCloseMessage(websocket.CloseUnsupportedData, err.Error())
	ws.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
}

// SQLError wraps [*pq.Error] errors only
func SQLError(err error) error {
	if err, ok := err.(*pq.Error); ok {
		//lint:ignore ST1005 Erreur utilisateur
		return fmt.Errorf("La requête SQL a échoué : %s (table %s)", err, err.Table)
	}
	return err
}

// InTx démarre une transaction, execute [fn] et
// COMMIT. En cas d'erreur, la transaction est ROLLBACK,
// et l'erreur renvoyé est passée à [SQLError]
func InTx(db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return SQLError(err)
	}
	err = fn(tx)
	if err != nil {
		_ = tx.Rollback()
		return SQLError(err)
	}
	err = tx.Commit()
	if err != nil {
		return SQLError(err)
	}
	return nil
}

// QueryParamInt64 parse the query param `name` to an int64
func QueryParamInt64(c echo.Context, name string) (int64, error) {
	idS := c.QueryParam(name)
	id, err := strconv.ParseInt(idS, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid ID parameter %s : %s", idS, err)
	}
	return id, nil
}

// QueryParamBool parse the query param `name` to a boolean
func QueryParamBool(c echo.Context, name string) bool {
	idS := c.QueryParam(name)
	return idS != ""
}

// Sample choose an index between 0 and len(weights)-1 at random, with the given weights,
// which must sum up to 1.
func SampleIndex(weights []float64) int {
	cumWeights := make([]float64, len(weights)) // last entry is 1
	sum := 0.
	for i, w := range weights {
		sum += w
		cumWeights[i] = sum
	}
	alea := rand.Float64()
	for i, cumWeight := range cumWeights {
		if alea < cumWeight {
			return i
		}
	}
	return len(weights) - 1
}

// BuildUrl returns the url composed of <host><path>?<query>.
func BuildUrl(host, path string, query map[string]string) string {
	pm := url.Values{}
	for k, v := range query {
		pm.Add(k, v)
	}
	u := url.URL{
		Host:     host,
		Scheme:   "https",
		Path:     path,
		RawQuery: pm.Encode(),
	}
	if strings.HasPrefix(host, "localhost") {
		u.Scheme = "http"
	}
	return u.String()
}

// Shuffler is a permutation used to shuffle a slice,
// storing the shuffled to original indices map.
type Shuffler []int

// Shuffle shuffles a slice using `moveTo`, which should be, in pseudo code :
// destination[dst] = source[src]
func (sh Shuffler) Shuffle(moveTo func(dst, src int)) {
	for i, mappedIndex := range sh {
		moveTo(i, mappedIndex)
	}
}

// OriginalToShuffled reverse the permutation, returning
// the original index -> shuffled index map.
func (sh Shuffler) OriginalToShuffled() []int {
	// build the reversed permutation
	out := make([]int, len(sh))
	for shuffled, original := range sh {
		out[original] = shuffled
	}
	return out
}

func newSeed(hash []byte) int64 {
	s := fnv.New32()
	s.Write(hash)
	return int64(s.Sum32())
}

func NewDeterministicRand(hash []byte) *rand.Rand {
	return rand.New(rand.NewSource(newSeed(hash)))
}

// NewDeterministicShuffler returns a random permutation of length `n`, seeded
// with a value computed from `hash`.
// The seed is adjusted to make sure the permuation is not the identity (unless n <= 1).
func NewDeterministicShuffler(hash []byte, n int) Shuffler {
	startSeed := newSeed(hash)
	out := rand.New(rand.NewSource(startSeed)).Perm(n)
	if n <= 1 { // nothing to check
		return out
	}

	isIdentity := func(perm []int) bool {
		for i, v := range perm {
			if i != v {
				return false
			}
		}
		return true
	}

	for ; isIdentity(out); startSeed++ {
		out = rand.New(rand.NewSource(startSeed)).Perm(n)
	}

	return out
}

func RemoveAccents(s string) string {
	noAccent := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(noAccent, s)
	if e != nil {
		return s
	}
	return output
}

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](values ...T) Set[T] {
	out := make(Set[T], len(values))
	for _, v := range values {
		out.Add(v)
	}
	return out
}

func (s Set[T]) Has(key T) bool { _, ok := s[key]; return ok }

func (s Set[T]) Add(key T) { s[key] = struct{}{} }

func (s Set[T]) Delete(key T) { delete(s, key) }

func (s Set[T]) Keys() []T {
	out := make([]T, 0, len(s))
	for k := range s {
		out = append(out, k)
	}
	return out
}

package ceintures

// Code generated by gomacro/generator/go/sqlcrud. DO NOT EDIT.

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	"github.com/lib/pq"
)

type scanner interface {
	Scan(...interface{}) error
}

// DB groups transaction like objects, and
// is implemented by *sql.DB and *sql.Tx
type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
}

func scanOneBeltevolution(row scanner) (Beltevolution, error) {
	var item Beltevolution
	err := row.Scan(
		&item.IdStudent,
		&item.Level,
		&item.Advance,
		&item.Stats,
	)
	return item, err
}

func ScanBeltevolution(row *sql.Row) (Beltevolution, error) { return scanOneBeltevolution(row) }

// SelectAll returns all the items in the beltevolutions table.
func SelectAllBeltevolutions(db DB) (Beltevolutions, error) {
	rows, err := db.Query("SELECT * FROM beltevolutions")
	if err != nil {
		return nil, err
	}
	return ScanBeltevolutions(rows)
}

type Beltevolutions []Beltevolution

func ScanBeltevolutions(rs *sql.Rows) (Beltevolutions, error) {
	var (
		item Beltevolution
		err  error
	)
	defer func() {
		errClose := rs.Close()
		if err == nil {
			err = errClose
		}
	}()
	structs := make(Beltevolutions, 0, 16)
	for rs.Next() {
		item, err = scanOneBeltevolution(rs)
		if err != nil {
			return nil, err
		}
		structs = append(structs, item)
	}
	if err = rs.Err(); err != nil {
		return nil, err
	}
	return structs, nil
}

func InsertBeltevolution(db DB, item Beltevolution) error {
	_, err := db.Exec(`INSERT INTO beltevolutions (
			idstudent, level, advance, stats
			) VALUES (
			$1, $2, $3, $4
			);
			`, item.IdStudent, item.Level, item.Advance, item.Stats)
	if err != nil {
		return err
	}
	return nil
}

// Insert the links Beltevolution in the database.
// It is a no-op if 'items' is empty.
func InsertManyBeltevolutions(tx *sql.Tx, items ...Beltevolution) error {
	if len(items) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("beltevolutions",
		"idstudent",
		"level",
		"advance",
		"stats",
	))
	if err != nil {
		return err
	}

	for _, item := range items {
		_, err = stmt.Exec(item.IdStudent, item.Level, item.Advance, item.Stats)
		if err != nil {
			return err
		}
	}

	if _, err = stmt.Exec(); err != nil {
		return err
	}

	if err = stmt.Close(); err != nil {
		return err
	}
	return nil
}

// Delete the link Beltevolution from the database.
// Only the foreign keys IdStudent fields are used in 'item'.
func (item Beltevolution) Delete(tx DB) error {
	_, err := tx.Exec(`DELETE FROM beltevolutions WHERE IdStudent = $1;`, item.IdStudent)
	return err
}

// ByIdStudent returns a map with 'IdStudent' as keys.
func (items Beltevolutions) ByIdStudent() map[teacher.IdStudent]Beltevolution {
	out := make(map[teacher.IdStudent]Beltevolution, len(items))
	for _, target := range items {
		out[target.IdStudent] = target
	}
	return out
}

// IdStudents returns the list of ids of IdStudent
// contained in this link table.
// They are not garanteed to be distinct.
func (items Beltevolutions) IdStudents() []teacher.IdStudent {
	out := make([]teacher.IdStudent, len(items))
	for index, target := range items {
		out[index] = target.IdStudent
	}
	return out
}

// SelectBeltevolutionByIdStudent return zero or one item, thanks to a UNIQUE SQL constraint.
func SelectBeltevolutionByIdStudent(tx DB, idStudent teacher.IdStudent) (item Beltevolution, found bool, err error) {
	row := tx.QueryRow("SELECT * FROM beltevolutions WHERE idstudent = $1", idStudent)
	item, err = ScanBeltevolution(row)
	if err == sql.ErrNoRows {
		return item, false, nil
	}
	return item, true, err
}

func SelectBeltevolutionsByIdStudents(tx DB, idStudents_ ...teacher.IdStudent) (Beltevolutions, error) {
	rows, err := tx.Query("SELECT * FROM beltevolutions WHERE idstudent = ANY($1)", teacher.IdStudentArrayToPQ(idStudents_))
	if err != nil {
		return nil, err
	}
	return ScanBeltevolutions(rows)
}

func DeleteBeltevolutionsByIdStudents(tx DB, idStudents_ ...teacher.IdStudent) (Beltevolutions, error) {
	rows, err := tx.Query("DELETE FROM beltevolutions WHERE idstudent = ANY($1) RETURNING *", teacher.IdStudentArrayToPQ(idStudents_))
	if err != nil {
		return nil, err
	}
	return ScanBeltevolutions(rows)
}

func scanOneBeltquestion(row scanner) (Beltquestion, error) {
	var item Beltquestion
	err := row.Scan(
		&item.Id,
		&item.Domain,
		&item.Rank,
		&item.Parameters,
		&item.Enonce,
		&item.Correction,
	)
	return item, err
}

func ScanBeltquestion(row *sql.Row) (Beltquestion, error) { return scanOneBeltquestion(row) }

// SelectAll returns all the items in the beltquestions table.
func SelectAllBeltquestions(db DB) (Beltquestions, error) {
	rows, err := db.Query("SELECT * FROM beltquestions")
	if err != nil {
		return nil, err
	}
	return ScanBeltquestions(rows)
}

// SelectBeltquestion returns the entry matching 'id'.
func SelectBeltquestion(tx DB, id IdBeltquestion) (Beltquestion, error) {
	row := tx.QueryRow("SELECT * FROM beltquestions WHERE id = $1", id)
	return ScanBeltquestion(row)
}

// SelectBeltquestions returns the entry matching the given 'ids'.
func SelectBeltquestions(tx DB, ids ...IdBeltquestion) (Beltquestions, error) {
	rows, err := tx.Query("SELECT * FROM beltquestions WHERE id = ANY($1)", IdBeltquestionArrayToPQ(ids))
	if err != nil {
		return nil, err
	}
	return ScanBeltquestions(rows)
}

type Beltquestions map[IdBeltquestion]Beltquestion

func (m Beltquestions) IDs() []IdBeltquestion {
	out := make([]IdBeltquestion, 0, len(m))
	for i := range m {
		out = append(out, i)
	}
	return out
}

func ScanBeltquestions(rs *sql.Rows) (Beltquestions, error) {
	var (
		s   Beltquestion
		err error
	)
	defer func() {
		errClose := rs.Close()
		if err == nil {
			err = errClose
		}
	}()
	structs := make(Beltquestions, 16)
	for rs.Next() {
		s, err = scanOneBeltquestion(rs)
		if err != nil {
			return nil, err
		}
		structs[s.Id] = s
	}
	if err = rs.Err(); err != nil {
		return nil, err
	}
	return structs, nil
}

// Insert one Beltquestion in the database and returns the item with id filled.
func (item Beltquestion) Insert(tx DB) (out Beltquestion, err error) {
	row := tx.QueryRow(`INSERT INTO beltquestions (
		domain, rank, parameters, enonce, correction
		) VALUES (
		$1, $2, $3, $4, $5
		) RETURNING *;
		`, item.Domain, item.Rank, item.Parameters, item.Enonce, item.Correction)
	return ScanBeltquestion(row)
}

// Update Beltquestion in the database and returns the new version.
func (item Beltquestion) Update(tx DB) (out Beltquestion, err error) {
	row := tx.QueryRow(`UPDATE beltquestions SET (
		domain, rank, parameters, enonce, correction
		) = (
		$1, $2, $3, $4, $5
		) WHERE id = $6 RETURNING *;
		`, item.Domain, item.Rank, item.Parameters, item.Enonce, item.Correction, item.Id)
	return ScanBeltquestion(row)
}

// Deletes the Beltquestion and returns the item
func DeleteBeltquestionById(tx DB, id IdBeltquestion) (Beltquestion, error) {
	row := tx.QueryRow("DELETE FROM beltquestions WHERE id = $1 RETURNING *;", id)
	return ScanBeltquestion(row)
}

// Deletes the Beltquestion in the database and returns the ids.
func DeleteBeltquestionsByIDs(tx DB, ids ...IdBeltquestion) ([]IdBeltquestion, error) {
	rows, err := tx.Query("DELETE FROM beltquestions WHERE id = ANY($1) RETURNING id", IdBeltquestionArrayToPQ(ids))
	if err != nil {
		return nil, err
	}
	return ScanIdBeltquestionArray(rows)
}

// SelectBeltquestionsByDomainAndRank selects the items matching the given fields.
func SelectBeltquestionsByDomainAndRank(tx DB, domain Domain, rank Rank) (item Beltquestions, err error) {
	rows, err := tx.Query("SELECT * FROM beltquestions WHERE Domain = $1 AND Rank = $2", domain, rank)
	if err != nil {
		return nil, err
	}
	return ScanBeltquestions(rows)
}

// DeleteBeltquestionsByDomainAndRank deletes the item matching the given fields, returning
// the deleted items.
func DeleteBeltquestionsByDomainAndRank(tx DB, domain Domain, rank Rank) (item Beltquestions, err error) {
	rows, err := tx.Query("DELETE FROM beltquestions WHERE Domain = $1 AND Rank = $2 RETURNING *", domain, rank)
	if err != nil {
		return nil, err
	}
	return ScanBeltquestions(rows)
}

func loadJSON(out interface{}, src interface{}) error {
	if src == nil {
		return nil //zero value out
	}
	bs, ok := src.([]byte)
	if !ok {
		return errors.New("not a []byte")
	}
	return json.Unmarshal(bs, out)
}

func dumpJSON(s interface{}) (driver.Value, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return driver.Value(string(b)), nil
}

func IdBeltquestionArrayToPQ(ids []IdBeltquestion) pq.Int64Array {
	out := make(pq.Int64Array, len(ids))
	for i, v := range ids {
		out[i] = int64(v)
	}
	return out
}

// ScanIdBeltquestionArray scans the result of a query returning a
// list of ID's.
func ScanIdBeltquestionArray(rs *sql.Rows) ([]IdBeltquestion, error) {
	defer rs.Close()
	ints := make([]IdBeltquestion, 0, 16)
	var err error
	for rs.Next() {
		var s IdBeltquestion
		if err = rs.Scan(&s); err != nil {
			return nil, err
		}
		ints = append(ints, s)
	}
	if err = rs.Err(); err != nil {
		return nil, err
	}
	return ints, nil
}

type IdBeltquestionSet map[IdBeltquestion]bool

func NewIdBeltquestionSetFrom(ids []IdBeltquestion) IdBeltquestionSet {
	out := make(IdBeltquestionSet, len(ids))
	for _, key := range ids {
		out[key] = true
	}
	return out
}

func (s IdBeltquestionSet) Add(id IdBeltquestion) { s[id] = true }

func (s IdBeltquestionSet) Has(id IdBeltquestion) bool { return s[id] }

func (s IdBeltquestionSet) Keys() []IdBeltquestion {
	out := make([]IdBeltquestion, 0, len(s))
	for k := range s {
		out = append(out, k)
	}
	return out
}

func (s *Advance) Scan(src interface{}) error  { return loadJSON(s, src) }
func (s Advance) Value() (driver.Value, error) { return dumpJSON(s) }

func (s *Stats) Scan(src interface{}) error  { return loadJSON(s, src) }
func (s Stats) Value() (driver.Value, error) { return dumpJSON(s) }
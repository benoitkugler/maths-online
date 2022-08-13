package trivial

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/pass"
	"github.com/benoitkugler/maths-online/prof/teacher"
	tv "github.com/benoitkugler/maths-online/trivial"
	"github.com/benoitkugler/maths-online/utils/testutils"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// Simulate a real world example to exercice the server

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// student client connection
type student struct {
	serverBaseURL string
	gameCode      string

	gameMeta string
	conn     *websocket.Conn
}

func newStudent(serverBaseURL, gameCode string) *student {
	return &student{serverBaseURL: serverBaseURL, gameCode: gameCode}
}

func (st *student) accessGame(deconnect bool) {
	st.setupRequest()
	st.connectRequest(deconnect)
}

func (st *student) setupRequest() {
	u, err := url.Parse(st.serverBaseURL + studentSetup)
	check(err)

	query := make(url.Values)
	query.Set("session-id", st.gameCode)
	query.Set("game-meta", st.gameMeta)
	// anonymous connection
	u.RawQuery = query.Encode()

	resp, err := http.Get(u.String())
	check(err)

	var meta SetupStudentClientOut
	err = json.NewDecoder(resp.Body).Decode(&meta)
	check(err)

	st.gameMeta = meta.GameMeta
}

func (st *student) connectRequest(deconnect bool) {
	u, err := url.Parse(testutils.WebsocketURL(st.serverBaseURL + studentConnect))
	check(err)

	query := make(url.Values)
	query.Set("game-meta", st.gameMeta)
	query.Set("client-pseudo", "Benoit")
	u.RawQuery = query.Encode()

	st.conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	check(err)

	// pump server messages
	for {
		var v tv.StateUpdate
		err := st.conn.ReadJSON(&v)
		check(err)

		// fmt.Println(v)

		// simulate deconnection/reconnection
		if deconnect {
			break
		}
	}

	st.decoReco()
}

func (st *student) decoReco() {
	time.Sleep(time.Second / 10)
	st.conn.WriteControl(websocket.CloseMessage, nil, time.Now().Add(time.Second))
	st.conn.Close()

	time.Sleep(time.Second / 10)

	st.accessGame(false)
}

// teacherC monitor session
type teacherC struct {
	serverBaseURL string
}

func (tc *teacherC) monitorRequest() {
	u, err := url.Parse(testutils.WebsocketURL(tc.serverBaseURL + teacherMonitor))
	check(err)

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	check(err)

	// pump monitor messages
	for {
		var v teacherSocketData
		err := conn.ReadJSON(&v)
		check(err)
	}
}

// endpoints
const (
	studentSetup   = "/setup"
	studentConnect = "/connect"
	teacherMonitor = "/monitor"
)

type server struct {
	e  *echo.Echo
	ct *Controller
}

func (s server) handle(w http.ResponseWriter, r *http.Request) {
	var err error
	context := s.e.NewContext(r, w)
	context.Set("user", &jwt.Token{Claims: &teacher.UserMeta{Teacher: teacher.Teacher{Id: 1}}})

	switch url := r.URL; url.Path {
	case studentSetup:
		err = s.ct.SetupStudentClient(context)
	case studentConnect:
		err = s.ct.ConnectStudentSession(context)
	case teacherMonitor:
		err = s.ct.ConnectTeacherMonitor(context)
	default:
		panic(url)
	}

	if err != nil {
		panic(err)
	}
}

func TestSessionPlay(t *testing.T) {
	tv.ProgressLogger.SetOutput(os.Stdout)

	db, err := testutils.DB.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", testutils.DB, err)
		return
	}
	defer db.Close()

	ct := NewController(db, pass.Encrypter{}, "", teacher.Teacher{Id: 1}) // 1 is the defaut admin

	config, err := Trivial{Questions: demoQuestions, IdTeacher: ct.admin.Id}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	s := server{e: e, ct: ct}

	groupSize := []int{2, 3, 4}
	groups, err := ct.launchConfig(LaunchSessionIn{
		IdConfig: config.Id,
		Groups:   groupSize,
	},
		ct.admin.Id,
	)
	if err != nil {
		t.Fatal(err)
	}

	// listen to "external" requests
	listener := httptest.NewServer(http.HandlerFunc(s.handle))
	defer listener.Close()

	tc1 := teacherC{listener.URL}
	go tc1.monitorRequest()

	tc2 := teacherC{listener.URL}
	go tc2.monitorRequest()

	time.Sleep(50 * time.Millisecond)

	// create the student clients
	var allStudents [][]*student
	for i, roomCode := range groups.GameIDs {
		size := groupSize[i]

		students := make([]*student, size)
		for j := range students {
			st := newStudent(listener.URL, roomCode)

			go st.accessGame(j == 0)

			students[j] = st
		}
		allStudents = append(allStudents, students)
	}

	time.Sleep(time.Second / 2)
}

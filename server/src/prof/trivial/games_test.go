package trivial

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/server/src/pass"
	tcAPI "github.com/benoitkugler/maths-online/server/src/prof/teacher"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	tr "github.com/benoitkugler/maths-online/server/src/sql/trivial"
	tv "github.com/benoitkugler/maths-online/server/src/trivial"
	"github.com/benoitkugler/maths-online/server/src/utils/testutils"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// Simulate a real world example to exercice the server

func check(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatal(err)
	}
}

// student client connection
type student struct {
	t *testing.T

	serverBaseURL string
	gameCode      tv.RoomID

	gameMeta string
	conn     *websocket.Conn
}

func newStudent(t *testing.T, serverBaseURL string, gameCode tv.RoomID) *student {
	return &student{t: t, serverBaseURL: serverBaseURL, gameCode: gameCode}
}

func (st *student) accessGame(deconnect bool) {
	st.setupRequest()
	st.connectRequest(deconnect)
}

func (st *student) setupRequest() {
	u, err := url.Parse(st.serverBaseURL + studentSetup)
	check(st.t, err)

	query := make(url.Values)
	query.Set("session-id", string(st.gameCode))
	query.Set("game-meta", st.gameMeta)
	// anonymous connection
	u.RawQuery = query.Encode()

	resp, err := http.Get(u.String())
	check(st.t, err)

	var meta SetupStudentClientOut
	err = json.NewDecoder(resp.Body).Decode(&meta)
	check(st.t, err)

	st.gameMeta = meta.GameMeta
}

func (st *student) connectRequest(deconnect bool) {
	u, err := url.Parse(testutils.WebsocketURL(st.serverBaseURL + studentConnect))
	check(st.t, err)

	query := make(url.Values)
	query.Set("game-meta", st.gameMeta)
	if rand.Intn(2) == 0 {
		query.Set("client-pseudo", "Benoit")
	}
	u.RawQuery = query.Encode()

	st.conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	check(st.t, err)

	// pump server messages
	for {
		var v tv.StateUpdate
		err := st.conn.ReadJSON(&v)
		check(st.t, err)

		// simulate deconnection/reconnection
		if deconnect {
			break
		}
	}

	st.decoReco()
}

func (st *student) decoReco() {
	time.Sleep(time.Millisecond * 100)
	st.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseGoingAway, "simulating deconnection"), time.Now().Add(time.Second))
	st.conn.Close()

	time.Sleep(time.Millisecond * 100)

	st.accessGame(false)
}

// teacherC monitor session
type teacherC struct {
	t *testing.T

	server *httptest.Server
}

func (tc *teacherC) monitorRequest() {
	u, err := url.Parse(tc.server.URL + teacherMonitor)
	check(tc.t, err)

	// pump monitor messages
	for range [10]int{} {
		resp, err := http.Get(u.String())
		check(tc.t, err)
		var v MonitorOut
		err = json.NewDecoder(resp.Body).Decode(&v)
		check(tc.t, err)

		resp.Body.Close()

		time.Sleep(time.Millisecond * 10)
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
	context.Set("user", &jwt.Token{Claims: &tcAPI.UserMeta{IdTeacher: 1}})

	switch url := r.URL; url.Path {
	case studentSetup:
		err = s.ct.SetupStudentClient(context)
	case studentConnect:
		err = s.ct.ConnectStudentSession(context)
	case teacherMonitor:
		err = s.ct.TrivialTeacherMonitor(context)
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

	config, err := tr.Trivial{Questions: demoQuestions, IdTeacher: ct.admin.Id}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	s := server{e: e, ct: ct}

	groupSize := []int{2, 3, 4}
	groups, err := ct.launchConfig(LaunchSessionIn{
		IdConfig: config.Id,
		Groups:   GroupsStrategyAuto{groupSize},
	},
		ct.admin.Id,
	)
	if err != nil {
		t.Fatal(err)
	}

	// listen to "external" requests
	listener := httptest.NewServer(http.HandlerFunc(s.handle))
	defer listener.Close()

	tc1 := teacherC{t, listener}
	go tc1.monitorRequest()

	tc2 := teacherC{t, listener}
	go tc2.monitorRequest()

	time.Sleep(50 * time.Millisecond)

	// create the student clients
	var allStudents [][]*student
	for i, roomCode := range groups.GameIDs {
		size := groupSize[i]

		students := make([]*student, size)
		for j := range students {
			st := newStudent(t, listener.URL, roomCode)

			go st.accessGame(j == 0)

			students[j] = st
		}
		allStudents = append(allStudents, students)
	}

	time.Sleep(time.Second * 3)
}

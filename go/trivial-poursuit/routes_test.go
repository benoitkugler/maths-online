package trivialpoursuit

import (
	"net/url"
	"os"
	"testing"
	"time"
)

func TestController_buildURL(t *testing.T) {
	localhost, _ := url.Parse("http://localhost:1323")
	deployed, _ := url.Parse("https://www.free.fr")

	tests := []struct {
		host *url.URL
		args string
		want string
	}{
		{localhost, "/trivial/AAA", "http://localhost:1323/trivial/AAA"},
		{deployed, "/trivial/AAA", "https://www.free.fr/trivial/AAA"},
	}
	for _, tt := range tests {
		ct := Controller{
			host: tt.host,
		}
		if got := ct.buildURL(tt.args, false); got != tt.want {
			t.Errorf("Controller.buildURL() = %v, want %v", got, tt.want)
		}
	}
}

func TestGameTimeout(t *testing.T) {
	ProgressLogger.SetOutput(os.Stdout)
	const timeout = time.Second / 10

	ct := NewController("localhost")
	ct.gameTimeout = timeout

	ct.launchGame(LaunchGameIn{NbPlayers: 2})

	if len(ct.stats()) != 1 {
		t.Fatal("expected one game")
	}

	time.Sleep(2 * timeout) // wait for the timeout

	if len(ct.stats()) != 0 {
		t.Fatal("game should have timed out")
	}
}

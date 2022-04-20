package trivialpoursuit

// func TestController_buildURL(t *testing.T) {
// 	localhost, _ := url.Parse("http://localhost:1323")
// 	deployed, _ := url.Parse("https://www.free.fr")

// 	tests := []struct {
// 		host *url.URL
// 		args string
// 		want string
// 	}{
// 		{localhost, "/trivial/AAA", "http://localhost:1323/trivial/AAA"},
// 		{deployed, "/trivial/AAA", "https://www.free.fr/trivial/AAA"},
// 	}
// 	for _, tt := range tests {
// 		ct := Controller{
// 			host: tt.host,
// 		}
// 		if got := ct.buildURL(tt.args, false); got != tt.want {
// 			t.Errorf("Controller.buildURL() = %v, want %v", got, tt.want)
// 		}
// 	}
// }

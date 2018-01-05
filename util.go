package tea

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
	"time"

	"github.com/k4s/tea/message"
)

func mkTimer(deadline time.Duration) <-chan time.Time {

	if deadline == 0 {
		return nil
	}

	return time.After(deadline)
}

var debug = true

func debugf(format string, args ...interface{}) {
	if debug {
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			file = "<?>"
			line = 0
		} else {
			if i := strings.LastIndex(file, "/"); i >= 0 {
				file = file[i+1:]
			}
		}
		fmt.Printf("DEBUG: %s:%d [%s]: %s\n", file, line,
			time.Now().String(), fmt.Sprintf(format, args...))
	}
}

func DrainChannel(ch chan<- *message.Message, expire time.Time) bool {
	var dur = time.Millisecond * 10

	for {
		if len(ch) == 0 {
			return true
		}
		now := time.Now()
		if now.After(expire) {
			return false
		}

		dur = expire.Sub(now)
		if dur > time.Millisecond*10 {
			dur = time.Millisecond * 10
		}
		time.Sleep(dur)
	}
}

func InitConfigFile(filePath string, config interface{}) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, config)
	if err != nil {
		return err
	}
	return nil
}

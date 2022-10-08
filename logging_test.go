package logging_test

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"

	logging "github.com/paweljvm/sgll"
)

func TestLogging(t *testing.T) {

	t.Run("should handle concurrent logging to file", func(t *testing.T) {
		// given
		os.Remove("out.log")
		//ioutil.WriteFile("out.log", []byte{}, fs.FileMode(os.O_WRONLY|os.O_CREATE))
		logging.LogToFile("out.log")
		logging.SetLevel(logging.DEBUG)
		// when
		latch := &sync.WaitGroup{}
		latch.Add(100)
		for i := 0; i < 25; i++ {
			go func() {
				defer latch.Done()
				logging.Debug("Some debug message {}", 123)
			}()
			go func() {
				defer latch.Done()
				logging.Error("Some error message {} {}", "test", "test2")
			}()
			go func() {
				defer latch.Done()
				logging.Info("Some info message {}", nil)
			}()
			go func() {
				defer latch.Done()
				logging.Warn("Some warn message")
			}()
		}
		latch.Wait()
		logging.CloseLogFile()
		// then
		result, err := ioutil.ReadFile("out.log")
		if err != nil {
			t.Errorf("Error during reading log file" + err.Error())
		}
		messages := strings.Split(string(result), "\r\n")
		counts := make(map[string]int)
		for _, message := range messages {
			if len(message) > 0 {
				logParts := strings.Split(message, " ")
				logType := logParts[0]
				counts[logType] = counts[logType] + 1
			}
		}
		if countsLen := len(counts); countsLen != 4 {
			t.Errorf("Expected to get 4 log types but was " + strconv.Itoa(countsLen))
		}
		for k, v := range counts {
			if v != 25 {
				t.Errorf("Expected to have 25 occurences of log type " + k)
			}
		}
	})
}

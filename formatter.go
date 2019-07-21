package fk_logrus_formatter

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	defaultLogOutputFormat = "%time% - %level% - [%package%::%file%::%function%::%line%] - %msg%\n"
	defaultTimestampFormat = "2006-01-02 15:04:05"
)

type FkLogFormatter struct {
	logOutput string
	timeStamp string
}

func (f *FkLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	output := f.logOutput

	if output == "" {
		output = defaultLogOutputFormat
	}

	timestampFormat := f.timeStamp
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	output = strings.Replace(output, "%time%", entry.Time.Format(timestampFormat), 1)

	output = strings.Replace(output, "%msg%", entry.Message, 1)

	level := strings.ToUpper(entry.Level.String())
	output = strings.Replace(output, "%level%", fmt.Sprintf("%-7s", level), 1)

	for k, v := range entry.Data {
		if s, ok := v.(string); ok {
			output = strings.Replace(output, "%"+k+"%", s, 1)
		}
	}

	return []byte(output), nil
}

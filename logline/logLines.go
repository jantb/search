package logline

import (
	"go4.org/intern"
	"strings"
	"time"
)

var Member Void

func ReverseLogline(numbers []LogLine) {
	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
}

type Void struct{}

type LogLine struct {
	Level  *intern.Value
	System *intern.Value
	Time   int64
	Body   []*intern.Value
	Ids    map[*intern.Value]Void
}

func (l LogLine) GetTime() time.Time {
	return time.Unix(0, l.Time*1000000)
}
func (l LogLine) MatchOrNot(query []string, matches map[*intern.Value]Void) (bool, []*intern.Value) {

	var ret []*intern.Value
	if len(query) == 0 {
		return true, ret
	}
	for k := range matches {
		_, ok := l.Ids[k]
		if ok {
			return true, ret
		}
	}

	for _, q := range query {
		found := false
		for value := range l.Ids {
			val := value.Get().(string)
			if strings.Contains(val, q) {
				ret = append(ret, value)
				found = true
				break
			}
		}
		if !found {
			return false, []*intern.Value{}
		}
	}

	return true, ret
}

func (l LogLine) MatchOrNotCount(query []string, matches map[*intern.Value]Void) bool {

	if len(query) == 0 {
		return true
	}

	for k := range matches {
		_, ok := l.Ids[k]
		if ok {
			return true
		}
	}

	return false
}

func (l LogLine) GetBody() string {
	if l.Body == nil {
		return ""
	}

	var sb strings.Builder
	for i, value := range l.Body {
		if i == 0 {
			sb.WriteString(value.Get().(string))
		} else {
			sb.WriteString(" ")
			sb.WriteString(value.Get().(string))
		}
	}
	return sb.String()
}

func (l *LogLine) SetBody(body string) {
	s := strings.Split(body, " ")
	for _, part := range s {
		l.Body = append(l.Body, intern.GetByString(part))
	}
	l.Ids = make(map[*intern.Value]Void)
	l.Ids[l.Level] = Member
	l.Ids[l.System] = Member
	for _, value := range l.Body {
		l.Ids[value] = Member
	}
}

func (l LogLine) GetLevel() string {
	return l.Level.Get().(string)
}

func (l *LogLine) SetLevel(level string) {
	l.Level = intern.GetByString(level)
}

func (l LogLine) GetSystem() string {
	return l.System.Get().(string)
}

func (l *LogLine) SetSystem(s string) {
	l.System = intern.GetByString(s)
}

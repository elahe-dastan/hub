package response

import (
	"fmt"
	"strings"

	"github.com/elahe-dastan/applifier/message"
)

type Response interface {
	Marshal() string
	Unmarshal(string)
}

type Stop struct {
}

type Who struct {
	ID string
}

type List struct {
	ConcatedIds string
}

type Send struct {
	Body string
}

func (s Stop) Marshal() string {
	return fmt.Sprintf("%s", message.STOP)
}

//func (s Stop) Unmarshal() {
//	return fmt.Sprintf("%s", message.STOP)
//}

func (w *Who) Marshal() string {
	return fmt.Sprintf("%s,%s\n", message.WhoAmI, w.ID)
}

func (w *Who) Unmarshal(id string) {
	w.ID = id
}

func (l *List) Marshal() string {
	return fmt.Sprintf("%s,%s\n", message.ListClientIDs, l.ConcatedIds)
}

func (l *List) Unmarshal(ids string) {
	l.ConcatedIds = ids
}

// Body has "\n" itself so there is no need to add it
func (s *Send) Marshal() string {
	return fmt.Sprintf("%s,%s", message.SendMsg, s.Body)
}

func (s *Send) Unmarshal(body string) {
	s.Body = body
}

func Unmarshal(res string) Response {
	arr := strings.Split(res, ",")
	t := arr[0]

	switch t {
	case message.WhoAmI:
		w := &Who{}
		w.Unmarshal(arr[1])
		return w
	case message.ListClientIDs:
		l := &List{}
		l.Unmarshal(arr[1])
		return l
	case message.SendMsg:
		s := &Send{}
		s.Unmarshal(strings.Join(arr[1:], ","))
		return s
	}

	return nil
}

package response

import (
	"fmt"
	"strings"

	"github.com/elahe-dastan/applifier/message"
)

type Response interface {
	MarshalRes() string
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

func (s Stop) MarshalRes() string {
	return fmt.Sprintf("%s", message.STOP)
}

func (w Who) MarshalRes() string {
	return fmt.Sprintf("%s,%s\n", message.WhoAmI, w.ID)
}

func (w Who) Unmarshal(id string) Who {
	w.ID = id
	return w
}

func (l List) MarshalRes() string {
	return fmt.Sprintf("%s,%s\n", message.ListClientIDs, l.ConcatedIds)
}

func (l List) Unmarshal(ids string) List {
	l.ConcatedIds = ids
	return l
}

// Body has "\n" itself so there is no need to add it
func (s Send) MarshalRes() string {
	return fmt.Sprintf("%s,%s", message.SendMsg, s.Body)
}

func (s Send) Unmarshal(body string) Send {
	s.Body = body
	return s
}

func Unmarshal(res string) Response {
	arr := strings.Split(res, ",")
	t := arr[0]

	switch t {
	case message.WhoAmI:
		return Who{}.Unmarshal(arr[1])
	case message.ListClientIDs:
		return List{}.Unmarshal(arr[1])
	case message.SendMsg:
		return Send{}.Unmarshal(strings.Join(arr[1:], ","))
	}

	return nil
}

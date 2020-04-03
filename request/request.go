package request

import (
	"fmt"
	"strings"

	"github.com/elahe-dastan/applifier/message"
)

type Request interface {
	Marshal() string
}

type Stop struct {
}

func (s Stop) Marshal() string {
	panic("implement me")
}

type Who struct {
}

type List struct {
}

type Send struct {
	IDs  []string
	Body string
}

func (w Who) Marshal() string {
	return fmt.Sprintf("%s\n", message.WhoAmI)
}

//func (w *Who) Unmarshal(string) error {
//}

func (l List) Marshal() string {
	return fmt.Sprintf("%s\n", message.ListClientIDs)
}

//func (l *List) Unmarshal(string) error {
//}

// Body has "\n" itself so there is no need to add it
func (s Send) Marshal() string {
	ids := ""

	for _, id := range s.IDs {
		ids = ids + id + "-"
	}

	return fmt.Sprintf("%s,%s,%s", message.SendMsg, ids, s.Body)
}

//func (s *Send) Unmarshal(string) error {
//}

func Unmarshal(req string) Request {
	arr := strings.Split(req, ",")
	t := strings.TrimSpace(arr[0])

	switch t {
	case message.STOP:
		return Stop{}
	case message.WhoAmI:
		return Who{}
	case message.ListClientIDs:
		return List{}
	case message.SendMsg:
		recipientIDs := strings.TrimSuffix(arr[1], "-")
		recipientArr := strings.Split(recipientIDs, "-")

		return Send{
			IDs:  recipientArr,
			Body: strings.Join(arr[2:], ","),
		}
	}

	return nil
}

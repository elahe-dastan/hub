package request

import (
	"fmt"

	"github.com/elahe-dastan/applifier/message"
)

type Who struct {
}

type List struct {
}

type Send struct {
	IDs  []string
	Body string
}

func (w *Who) Marshal() string {
	return fmt.Sprintf("%s\n", message.WhoAmI)
}

//func (w *Who) Unmarshal(string) error {
//}

func (l *List) Marshal() string {
	return fmt.Sprintf("%s\n", message.ListClientIDs)
}

//func (l *List) Unmarshal(string) error {
//}

// Body has "\n" itself so there is no need to add it
func (s *Send) Marshal() string {
	ids := ""

	for _, id := range s.IDs {
		ids = ids + id + "-"
	}

	return fmt.Sprintf("%s,%s,%s", message.SendMsg, ids, s.Body)
}

//func (s *Send) Unmarshal(string) error {
//}

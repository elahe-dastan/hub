package response

import (
	"fmt"

	"github.com/elahe-dastan/applifier/message"
)

type Who struct {
	ID string
}

type List struct {
	IDs []string
}

type Send struct {
	Body string
}

func (w *Who) Marshal() string {
	return fmt.Sprintf("%s,%s\n", message.WhoAmI, w.ID)
}

//func (w *Who) Unmarshal(string) error {
//}

func (l *List) Marshal() string {
	ids := ""

	for _, id := range l.IDs {
		ids = ids + id + "-"
	}
	return fmt.Sprintf("%s,%s\n", message.ListClientIDs, ids)
}

//func (l *List) Unmarshal(string) error {
//}

// Body has "\n" itself so there is no need to add it
func (s *Send) Marshal() string {
	return fmt.Sprintf("%s,%s", message.SendMsg, s.Body)
}

//func (s *Send) Unmarshal(string) error {
//}

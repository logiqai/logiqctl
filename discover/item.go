package discover

import "fmt"

type ListItem struct {
	Name      string
	LastSeen  int64
	FirstSeen int64
}

func (l *ListItem) GetItemForDisplay() string {
	return fmt.Sprintf("%d-%d", l.LastSeen, l.FirstSeen)
}

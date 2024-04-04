package rk9

type Roster struct {
	Players []*EventPlayer
}

type Decklist struct {
	Pokemon  []Card
	Trainers []Card
	Energy   []Card
}

type Card struct {
	Count    int
	Name     string
	Print    string
	Number   int
	Language string
}

func GetRoster(event *Event) (*Roster, error) {
	return nil, nil
}

func GetDecklist(player *EventPlayer) (*Decklist, error) {
	return nil, nil
}

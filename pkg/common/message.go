package common

type Message string

func (m Message) String() string {
	return string(m)
}

package utils

type Ulid func() string

func (u Ulid) Generate() string {
	return u()
}

package serviceParser

type ResponseParser interface {
	ParseDTO(data any) error
}

func Parse(rp ResponseParser, data any) error {

	return rp.ParseDTO(data)
}

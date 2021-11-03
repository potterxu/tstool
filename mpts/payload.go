package mpts

type PayloadType struct {
	data []byte
}

func payload(data []byte) *PayloadType {
	return &PayloadType{
		data: data,
	}
}

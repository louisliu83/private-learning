package worker

const (
	DiffieHellman = "DH"
	OT            = "OT"
)

var (
	protoName2Num = map[string]string{
		"DH": "2",
		"OT": "3",
	}

	protoNum2Name = map[string]string{
		"2": "DH",
		"3": "OT",
	}
)

func protocolNumber(p string) string {
	switch p {
	case DiffieHellman:
		return "2"
	case OT:
		return "3"
	default:
		return "3"
	}
	return "3"
}

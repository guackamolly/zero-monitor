package conn

type msg struct {
	key byte
	id  string
}

func encode(d msg) []byte {
	b := make([]byte, 1+len(d.id))
	b[0] = d.key

	for i := range d.id {
		b[i+1] = d.id[i]
	}

	return b
}

func decode(b []byte) msg {
	var d msg

	l := len(b)
	if l == 0 {
		return d
	}

	d.key = b[0]
	d.id = string(b[1:])

	return d
}

func compose(key byte) msg {
	return msg{id: machineId, key: key}
}

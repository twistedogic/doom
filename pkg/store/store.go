package store

type Getter interface {
	Get(string) ([]byte, error)
}

type Setter interface {
	Set(string, []byte) error
}

type Scanner interface {
	Scan(...string) ([]string, error)
}

type Closer interface {
	Close() error
}

type Store interface {
	Getter
	Setter
	Scanner
	Closer
}

func Transfer(src, dst Store) error {
	keys, err := src.Scan()
	if err != nil {
		return err
	}
	for _, key := range keys {
		b, err := src.Get(key)
		if err != nil {
			return err
		}
		if err := dst.Set(key, b); err != nil {
			return err
		}
	}
	return nil
}

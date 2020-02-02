package store

type Store interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
	Scan(...string) ([]string, error)
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

package jwt

import (
	"fmt"
)

// secstr plain in struct & crypted in json.
type secstr string

func (ss secstr) MarshalJSON() ([]byte, error) {
	b, e := jwtkms.Encrypt("innervcode", "", []byte(ss))
	v := fmt.Sprintf("\"%x\"", b)
	return []byte(v), e
}

func (ss *secstr) UnmarshalJSON(p []byte) error {
	if string(p) == "null" {
		return nil
	}
	var v = []byte{}
	if n, err := fmt.Sscanf(string(p), "\"%x\"", &v); err != nil {
		return fmt.Errorf("not a valid secstr: '%s', err:'%s'", p, err)
	} else if n != 1 {
		return fmt.Errorf("not a valid secstr: '%s', n=%d", p, n)
	}

	b, err := jwtkms.Decrypt("innervcode", "", v)
	if err != nil {
		return fmt.Errorf("invalid vcode:%s, error:%s", string(p), err)
	}
	*ss = secstr(b)
	return nil
}

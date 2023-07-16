package httpauthmw

import "fmt"

var realmCharMap = [4]uint64{ //nolint:gochecknoglobals
	0b_101011111_1111111_11111111_11111011_00000000_00000000_00000000_00000000,
	0b_01000111_11111111_11111111_11111110_10101111_11111111_11111111_11111111,
	0b_00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
	0b_00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
}

func ValidateRialm(realm string) error {
	for _, c := range realm {
		if (realmCharMap[c/64] & (uint64(1) << (c % 64))) == 0 { //nolint:gomnd
			return fmt.Errorf("invalid char %q in realm %q", c, realm)
		}
	}
	return nil
}

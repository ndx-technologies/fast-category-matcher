package distance

// Lev computes Levenshtein distances with up to 256 length.
// For efficiency, temporary data can be stored in provided buffer.
func Lev(s1, s2 string, buf []byte) int {
	if s1 == s2 {
		return 0
	}
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	if len(s1) > len(s2) {
		s1, s2 = s2, s1
	}

	d := buf
	if len(d) < len(s1)+1 {
		d = make([]byte, len(s1)+1)
	} else {
		d = d[:len(s1)+1]
	}

	for j := 1; j <= len(s2); j++ {
		prevDiag := d[0]
		d[0] = byte(j)

		for i := 1; i <= len(s1); i++ {
			tmp := d[i]

			if s1[i-1] == s2[j-1] {
				d[i] = prevDiag
			} else {
				d[i]++

				if d[i-1]+1 < d[i] {
					d[i] = d[i-1] + 1
				}

				if prevDiag+1 < d[i] {
					d[i] = prevDiag + 1
				}
			}
			prevDiag = tmp
		}
	}

	return int(d[len(s1)])
}

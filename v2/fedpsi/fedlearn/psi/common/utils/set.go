package utils

func Union(s1, s2 []string) []string {
	if len(s1) == 0 {
		return s2
	}

	if len(s2) == 0 {
		return s1
	}

	hs := map[string]struct{}{}
	retS := make([]string, 0)

	for _, s := range s1 {
		if _, ok := hs[s]; !ok {
			hs[s] = struct{}{}
			retS = append(retS, s)
		}
	}

	for _, s := range s2 {
		if _, ok := hs[s]; !ok {
			hs[s] = struct{}{}
			retS = append(retS, s)
		}
	}

	return retS
}

func Intersect(s1, s2 []string) []string {

	if len(s1) == 0 || len(s2) == 0 {
		return s1
	}

	hs := map[string]struct{}{}
	retS := make([]string, 0)

	for _, s := range s2 {
		if _, ok := hs[s]; !ok {
			hs[s] = struct{}{}
		}
	}

	for _, s := range s1 {
		if _, ok := hs[s]; !ok {
			retS = append(retS, s)
		}
	}

	return retS
}

package cli

// flatten flattens a nested map. For example:
//   "root": {
//     "sub": {
//       "subone": "val",
//     },
//   }
//
// flattens to: {"root.sub.subone": "val"}
func flatten(in map[string]interface{}, prefix []string, kf KeyFunc, out map[string]interface{}) {
	for k, v := range in {
    path := append(prefix[:], k)

		switch x := v.(type) {
		case map[string]interface{}:
			flatten(x, path, kf, out)
		default:
      key := kf(path)
			out[key] = v
		}
	}
}

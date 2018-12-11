package cli

// flatten flattens a nested map. For example:
//   "root": {
//     "sub": {
//       "subone": "val",
//     },
//   }
//
// flattens to: {"root.sub.subone": "val"}
func flatten(in map[string]interface{}, out map[string]interface{}, prefix []string, kf KeyFunc) {
	for k, v := range in {
		path := append(prefix[:], k)

		switch x := v.(type) {
		case map[string]interface{}:
			flatten(x, out, path, kf)
		default:
			key := kf(path)
			out[key] = v
		}
	}
}

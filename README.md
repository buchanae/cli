
### Why

- Outside of the cli code, codebase deals directly is stucts. i.e. config/arg code
  does not leak into the codebase.

- Defaults are also defined by structs (i.e. pulled from the codebase) instead of
  being cli/arg specific

- Help/docs come from the code

- cli/arg/config code should live *around* your code, not *in* it.

### Done:
- env prefix
- flag/env var naming style ServiceName -> service-name or service_name or service.name
- set defaults from a struct
- auto inspect
- flag generation
- help/docs from comments
- read from flag, env, yaml, json
- support time.Duration in yaml, json, env
- support byte size
- support SI suffix (K, G, M, etc)
- alias/link/source field value from another field
- ignore/hide fields
- define short fields
- report unknown fields
- dump yaml
  - printing config, but only non-defaults

- validation interface
- choices + validation

### TODO:
- manage editing config file
- pluggable sources
- slice of choices
- improve stringSlice.String() format
- dump json, env, flags
- handle map[string]string via "key=value" flag value
- explore "storage.local.allowed_dirs.append"
- pull fieldname from json tag
- recognize misspelled env var
- case sensitivity
- inspect roger.Vals via reflection

### Complex:
- reloading
- multiple config files with merging

### Questions
- how to handle pointers? cycles?
- how are slices handled in env vars?
- how are slices of structs handled in flags?
- how to handle unknown type wrappers, e.g. type Foo int64


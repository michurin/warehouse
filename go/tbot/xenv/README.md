# It mimics systemd v253 EnvironmentFile option

The parser is borrowed from `systemd` `v253` as is despite the original parser slightly oversimplify and allows to do weird things, see tests.

- `[parse_env_file_internal](https://github.com/systemd/systemd/blob/v253/src/basic/env-file.c#L22)` — `systemd` implementation
- Useful constants: [[1](https://github.com/systemd/systemd/blob/v253/src/basic/string-util.h#L13)], [[2](https://github.com/systemd/systemd/blob/v253/src/basic/escape.h#L15)]
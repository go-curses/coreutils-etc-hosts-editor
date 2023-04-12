# etc-hosts-editor

Utility for managing the OS /etc/hosts file.

## INSTALLATION

``` shell
> go install github.com/go-curses/coreutils-etc-hosts-editor/cmd/eheditor@latest
```

## DOCUMENTATION

``` shell
> eheditor --help
NAME:
   eheditor - etc hosts editor

USAGE:
   eheditor [options] [/etc/hosts]

VERSION:
   v0.1.0 (trunk)

DESCRIPTION:
   command line utility for managing the OS /etc/hosts file

GLOBAL OPTIONS:
   --help, -h, --usage  display command-line usage information (default: false)
   --read-only, -r      do not write any changes to the etc hosts file (default: false)
   --version, -v        display the version (default: false)
```


## LICENSE

```
Copyright 2023  The Go-Curses Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use file except in compliance with the License.
You may obtain a copy of the license at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

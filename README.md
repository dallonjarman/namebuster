# Namebuster

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Tests](https://github.com/benbusby/namebuster/actions/workflows/tests.yml/badge.svg)](https://github.com/benbusby/namebuster/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/benbusby/namebuster)](https://goreportcard.com/report/github.com/benbusby/namebuster)

Generates a list of possible common username permutations given a list of names, a url, or a file.

## Install
Go: `go install github.com/dallonjarman/namebuster@latest`

## Usage
### Command Line
```bash
kali@kali:~$ namebuster                                            
                                                        
  Usage: namebuster <text|url|file>

  Example (names): namebuster "John Doe" > usernames.txt
  Example (single): namebuster "admin" > usernames.txt
  Example (url):    namebuster https://example.com > usernames.txt
  Example (file):   namebuster employees.txt > usernames.txt
```

For each discovered name, namebuster will generate ~200 possible usernames. You can then use this list with a tool like [kerbrute](https://github.com/ropnop/kerbrute), for example (originally used for the [Sauna](https://app.hackthebox.com/machines/Sauna) machine on [HackTheBox](https://hackthebox.com)):

```bash
[ kali : ~/test ]
$ namebuster https://sauna.htb > usernames.txt
[ kali : ~/test ]
$ ./kerbrute_linux_amd64 userenum ./usernames.txt -d DOMAIN.LOCAL --dc sauna.htb

    __             __               __
   / /_____  _____/ /_  _______  __/ /____
  / //_/ _ \/ ___/ __ \/ ___/ / / / __/ _ \
 / ,< /  __/ /  / /_/ / /  / /_/ / /_/  __/
/_/|_|\___/_/  /_.___/_/   \__,_/\__/\___/

Version: v1.0.3 (9dad6e1) - 02/18/20 - Ronnie Flathers @ropnop

2020/02/18 23:47:59 >  Using KDC(s):
2020/02/18 23:47:59 >  	domain.com:88

2020/02/18 23:47:59 >  [+] VALID USERNAME:	 fsmith@DOMAIN.LOCAL
2020/02/18 23:47:59 >  Done! Tested 125 usernames (1 valid) in 1.585 seconds
```

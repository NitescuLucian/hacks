
# edirb

## Why? 

Economic Directory Buster (edirb) for when you need to diretory bruteforce more then 10k hosts but without making too many requests on each host.

```
sc = status code
al = aproximated content length
```

### Helper

```
Usage of edirb
  -t int
        Number of threads for requests (default 100)
  -w string
        Wordlist file path for bruteforce
```

* You can `anew` and `grep` for what you need and it is easy to use while piping.

## Install

You can install using go:

```
go install -v github.com/NitescuLucian/hacks/edirb@latest
```

Contact me at [@LucianNitescu](https://twitter.com/LucianNitescu)


# urlop

## Why? 

Whenever you need every query parameter taken one by one for any given URL, you should use urlop.

It stands for: URL One Parameter

## Usage Example

```
echo "https://exmaple.com/?a=x&b=y" | urlop 
```

Output:

```
https://exmaple.com/?a=x
https://exmaple.com/?b=y
```

## Install

You can install using go:

```
go install -v github.com/NitescuLucian/hacks/urlop@latest
```

Contact me at [@LucianNitescu](https://twitter.com/LucianNitescu)

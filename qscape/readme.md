
# qscape

## Why? 

It stands for QueryUnescape function in golang.

## Usage Example

```
echo "http://example.com/?a=a" | qsreplace 'he<script>' | qscape 
```

Output:

```
http://example.com/?a=he<script>
```

## Install

You can install using go:

```
go install -v github.com/NitescuLucian/hacks/qscape@latest
```

Contact me at [@LucianNitescu](https://twitter.com/LucianNitescu)

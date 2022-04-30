# unsimi

Very similar to https://github.com/tomnomnom/anew but with a "lucky" twist, it appends lines from stdin to a file, but only if they don't already appear in the file and if the line is not similar to the previouse line.

Outputs new lines to `stdout` too, making it a bit like a `tee -a` that removes duplicates.

## Why? 

It is usefule to have a tool that can be used to append to a file, but only if the line is not already in the file and if it is not similar to the previous line.

It is usefule from times to times to avoid similar targets in a file. It is not perfect, but you might be lucky enough to filter out similar targets and reduce time of execution.

## Usage Example

Here, a file called `things.txt` contains a list of numbers. `newthings.txt` contains a second
list of numbers, some of which appear in `things.txt` and some of which do not. `unsimi` is used
to append the latter to `things.txt`.


```
▶ cat things.txt
example.com
cars.example.com
cars2.example.com
lol.example.com
example.example.com
api.example.com

▶ cat newthings.txt
example.com
cars.example.com
cars2.example.com
cars3.example.com
cars4.example.com
cars5.example.com
cars6.example.com
lol.example.com
example.example.com
api.example.com
abuse.exampe.com
example.org
example.net
example.com

▶ cat newthings.txt | unsimi
example.com
cars.example.com
lol.example.com
example.example.com
api.example.com
example.org

▶ cat newthings.txt | unsimi things.txt
cars3.example.com
example.org

▶ cat things.txt
example.com
cars.example.com
cars2.example.com
lol.example.com
example.example.com
api.example.com
cars3.example.com
example.org

```

Note that the new lines added to `things.txt` are also sent to `stdout`, this allows for them to
be redirected to another file:

```
▶ cat newthings.txt | unsimi things.txt > added-lines.txt
▶ cat added-lines.txt
Three
Four
```

## Flags

- To view the output in stdout, but not append to the file, use the dry-run option `-d`.
- To append to the file, but not print anything to stdout, use quiet mode `-q`.

## Install

You can either install using go:

```
go install -v github.com/NitescuLucian/hacks/unsimi@latest
```

# Thanks to

- https://github.com/tomnomnom/anew
- https://github.com/adrg/strutil#hamming
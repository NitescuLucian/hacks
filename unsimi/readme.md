# unsimi

Very similar to https://github.com/tomnomnom/anew but with a "lucky" twist, it appends lines from stdin to a file, but only if they don't already appear in the file and if they are not similar to the previouse line.

Outputs new lines to `stdout` too, making it a bit like a `tee -a` that removes duplicates.

## Why? 

It is usefule to have a tool that can be used to append to a file, but only if the line is not already in the file and if it is not similar to the previous line.

It is usefule from times to times to avoid similar targets in a file.

## Usage Example

Here, a file called `things.txt` contains a list of numbers. `newthings.txt` contains a second
list of numbers, some of which appear in `things.txt` and some of which do not. `unsimi` is used
to append the latter to `things.txt`.


```
▶ cat things.txt
Zero
One
Two

▶ cat newthings.txt
One
Two
Three
Four

▶ cat newthings.txt | unsimi things.txt
Three
Four

▶ cat things.txt
Zero
One
Two
Three
Four

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

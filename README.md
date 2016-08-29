# DirWalker

Walk specific directory to calculate SHA1 for all files.

## Install
`go get github.com/imeoer/dirwalker`

## Usage

```
dirwalker [path] [to/ignore/path]...
  -o Specific output file path
```

## Example

`dirwalker -o=result.txt /x /x/y/* /x/y.txt`

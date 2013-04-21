# Example godocdown (strings)

This markdown was generated with the help of custom template file ([example.template](http://github.com/robertkrimen/godocdown/blob/master/example.template)). To add custom
markdown to your documentation, you can do something like:

    godocdown -template=godocdown.tmpl ...

The template format is the standard Go text/template: http://golang.org/pkg/text/template

Along with the standard template functionality, the starting data argument has the following interface:

    {{ .Emit }}
    // Emit the standard documentation (what godocdown would emit without a template)

    {{ .EmitHeader }}
    // Emit the package name and an import line (if one is present/needed)

    {{ .EmitSynopsis }}
    // Emit the package declaration

    {{ .EmitUsage }}
    // Emit package usage, which includes a constants section, a variables section,
    // a functions section, and a types section. In addition, each type may have its own constant,
    // variable, and/or function/method listing.

    {{ if .IsCommand  }} ... {{ end }}
    // A boolean indicating whether the given package is a command or a plain package

    {{ .Name }}
    // The name of the package/command (string)

    {{ .ImportPath }}
    // The import path for the package (string)
    // (This field will be the empty string if godocdown is unable to guess it)

godocdown for the http://golang.org/pkg/strings package:

--

# strings
--
    import "strings"

Package strings implements simple functions to manipulate strings.

## Usage

#### func  Contains

```go
func Contains(s, substr string) bool
```
Contains returns true if substr is within s.

#### func  ContainsAny

```go
func ContainsAny(s, chars string) bool
```
ContainsAny returns true if any Unicode code points in chars are within s.

#### func  ContainsRune

```go
func ContainsRune(s string, r rune) bool
```
ContainsRune returns true if the Unicode code point r is within s.

#### func  Count

```go
func Count(s, sep string) int
```
Count counts the number of non-overlapping instances of sep in s.

#### func  EqualFold

```go
func EqualFold(s, t string) bool
```
EqualFold reports whether s and t, interpreted as UTF-8 strings, are equal under
Unicode case-folding.

#### func  Fields

```go
func Fields(s string) []string
```
Fields splits the string s around each instance of one or more consecutive white
space characters, as defined by unicode.IsSpace, returning an array of
substrings of s or an empty list if s contains only white space.

#### func  FieldsFunc

```go
func FieldsFunc(s string, f func(rune) bool) []string
```
FieldsFunc splits the string s at each run of Unicode code points c satisfying
f(c) and returns an array of slices of s. If all code points in s satisfy f(c)
or the string is empty, an empty slice is returned.

#### func  HasPrefix

```go
func HasPrefix(s, prefix string) bool
```
HasPrefix tests whether the string s begins with prefix.

#### func  HasSuffix

```go
func HasSuffix(s, suffix string) bool
```
HasSuffix tests whether the string s ends with suffix.

#### func  Index

```go
func Index(s, sep string) int
```
Index returns the index of the first instance of sep in s, or -1 if sep is not
present in s.

#### func  IndexAny

```go
func IndexAny(s, chars string) int
```
IndexAny returns the index of the first instance of any Unicode code point from
chars in s, or -1 if no Unicode code point from chars is present in s.

#### func  IndexFunc

```go
func IndexFunc(s string, f func(rune) bool) int
```
IndexFunc returns the index into s of the first Unicode code point satisfying
f(c), or -1 if none do.

#### func  IndexRune

```go
func IndexRune(s string, r rune) int
```
IndexRune returns the index of the first instance of the Unicode code point r,
or -1 if rune is not present in s.

#### func  Join

```go
func Join(a []string, sep string) string
```
Join concatenates the elements of a to create a single string. The separator
string sep is placed between elements in the resulting string.

#### func  LastIndex

```go
func LastIndex(s, sep string) int
```
LastIndex returns the index of the last instance of sep in s, or -1 if sep is
not present in s.

#### func  LastIndexAny

```go
func LastIndexAny(s, chars string) int
```
LastIndexAny returns the index of the last instance of any Unicode code point
from chars in s, or -1 if no Unicode code point from chars is present in s.

#### func  LastIndexFunc

```go
func LastIndexFunc(s string, f func(rune) bool) int
```
LastIndexFunc returns the index into s of the last Unicode code point satisfying
f(c), or -1 if none do.

#### func  Map

```go
func Map(mapping func(rune) rune, s string) string
```
Map returns a copy of the string s with all its characters modified according to
the mapping function. If mapping returns a negative value, the character is
dropped from the string with no replacement.

#### func  Repeat

```go
func Repeat(s string, count int) string
```
Repeat returns a new string consisting of count copies of the string s.

#### func  Replace

```go
func Replace(s, old, new string, n int) string
```
Replace returns a copy of the string s with the first n non-overlapping
instances of old replaced by new. If n < 0, there is no limit on the number of
replacements.

#### func  Split

```go
func Split(s, sep string) []string
```
Split slices s into all substrings separated by sep and returns a slice of the
substrings between those separators. If sep is empty, Split splits after each
UTF-8 sequence. It is equivalent to SplitN with a count of -1.

#### func  SplitAfter

```go
func SplitAfter(s, sep string) []string
```
SplitAfter slices s into all substrings after each instance of sep and returns a
slice of those substrings. If sep is empty, SplitAfter splits after each UTF-8
sequence. It is equivalent to SplitAfterN with a count of -1.

#### func  SplitAfterN

```go
func SplitAfterN(s, sep string, n int) []string
```
SplitAfterN slices s into substrings after each instance of sep and returns a
slice of those substrings. If sep is empty, SplitAfterN splits after each UTF-8
sequence. The count determines the number of substrings to return:

    n > 0: at most n substrings; the last substring will be the unsplit remainder.
    n == 0: the result is nil (zero substrings)
    n < 0: all substrings

#### func  SplitN

```go
func SplitN(s, sep string, n int) []string
```
SplitN slices s into substrings separated by sep and returns a slice of the
substrings between those separators. If sep is empty, SplitN splits after each
UTF-8 sequence. The count determines the number of substrings to return:

    n > 0: at most n substrings; the last substring will be the unsplit remainder.
    n == 0: the result is nil (zero substrings)
    n < 0: all substrings

#### func  Title

```go
func Title(s string) string
```
Title returns a copy of the string s with all Unicode letters that begin words
mapped to their title case.

BUG: The rule Title uses for word boundaries does not handle Unicode punctuation
properly.

#### func  ToLower

```go
func ToLower(s string) string
```
ToLower returns a copy of the string s with all Unicode letters mapped to their
lower case.

#### func  ToLowerSpecial

```go
func ToLowerSpecial(_case unicode.SpecialCase, s string) string
```
ToLowerSpecial returns a copy of the string s with all Unicode letters mapped to
their lower case, giving priority to the special casing rules.

#### func  ToTitle

```go
func ToTitle(s string) string
```
ToTitle returns a copy of the string s with all Unicode letters mapped to their
title case.

#### func  ToTitleSpecial

```go
func ToTitleSpecial(_case unicode.SpecialCase, s string) string
```
ToTitleSpecial returns a copy of the string s with all Unicode letters mapped to
their title case, giving priority to the special casing rules.

#### func  ToUpper

```go
func ToUpper(s string) string
```
ToUpper returns a copy of the string s with all Unicode letters mapped to their
upper case.

#### func  ToUpperSpecial

```go
func ToUpperSpecial(_case unicode.SpecialCase, s string) string
```
ToUpperSpecial returns a copy of the string s with all Unicode letters mapped to
their upper case, giving priority to the special casing rules.

#### func  Trim

```go
func Trim(s string, cutset string) string
```
Trim returns a slice of the string s with all leading and trailing Unicode code
points contained in cutset removed.

#### func  TrimFunc

```go
func TrimFunc(s string, f func(rune) bool) string
```
TrimFunc returns a slice of the string s with all leading and trailing Unicode
code points c satisfying f(c) removed.

#### func  TrimLeft

```go
func TrimLeft(s string, cutset string) string
```
TrimLeft returns a slice of the string s with all leading Unicode code points
contained in cutset removed.

#### func  TrimLeftFunc

```go
func TrimLeftFunc(s string, f func(rune) bool) string
```
TrimLeftFunc returns a slice of the string s with all leading Unicode code
points c satisfying f(c) removed.

#### func  TrimPrefix

```go
func TrimPrefix(s, prefix string) string
```
TrimPrefix returns s without the provided leading prefix string. If s doesn't
start with prefix, s is returned unchanged.

#### func  TrimRight

```go
func TrimRight(s string, cutset string) string
```
TrimRight returns a slice of the string s, with all trailing Unicode code points
contained in cutset removed.

#### func  TrimRightFunc

```go
func TrimRightFunc(s string, f func(rune) bool) string
```
TrimRightFunc returns a slice of the string s with all trailing Unicode code
points c satisfying f(c) removed.

#### func  TrimSpace

```go
func TrimSpace(s string) string
```
TrimSpace returns a slice of the string s, with all leading and trailing white
space removed, as defined by Unicode.

#### func  TrimSuffix

```go
func TrimSuffix(s, suffix string) string
```
TrimSuffix returns s without the provided trailing suffix string. If s doesn't
end with suffix, s is returned unchanged.

#### type Reader

```go
type Reader struct {
}
```

A Reader implements the io.Reader, io.ReaderAt, io.Seeker, io.WriterTo,
io.ByteScanner, and io.RuneScanner interfaces by reading from a string.

#### func  NewReader

```go
func NewReader(s string) *Reader
```
NewReader returns a new Reader reading from s. It is similar to
bytes.NewBufferString but more efficient and read-only.

#### func (*Reader) Len

```go
func (r *Reader) Len() int
```
Len returns the number of bytes of the unread portion of the string.

#### func (*Reader) Read

```go
func (r *Reader) Read(b []byte) (n int, err error)
```

#### func (*Reader) ReadAt

```go
func (r *Reader) ReadAt(b []byte, off int64) (n int, err error)
```

#### func (*Reader) ReadByte

```go
func (r *Reader) ReadByte() (b byte, err error)
```

#### func (*Reader) ReadRune

```go
func (r *Reader) ReadRune() (ch rune, size int, err error)
```

#### func (*Reader) Seek

```go
func (r *Reader) Seek(offset int64, whence int) (int64, error)
```
Seek implements the io.Seeker interface.

#### func (*Reader) UnreadByte

```go
func (r *Reader) UnreadByte() error
```

#### func (*Reader) UnreadRune

```go
func (r *Reader) UnreadRune() error
```

#### func (*Reader) WriteTo

```go
func (r *Reader) WriteTo(w io.Writer) (n int64, err error)
```
WriteTo implements the io.WriterTo interface.

#### type Replacer

```go
type Replacer struct {
}
```

A Replacer replaces a list of strings with replacements.

#### func  NewReplacer

```go
func NewReplacer(oldnew ...string) *Replacer
```
NewReplacer returns a new Replacer from a list of old, new string pairs.
Replacements are performed in order, without overlapping matches.

#### func (*Replacer) Replace

```go
func (r *Replacer) Replace(s string) string
```
Replace returns a copy of s with all replacements performed.

#### func (*Replacer) WriteString

```go
func (r *Replacer) WriteString(w io.Writer, s string) (n int, err error)
```
WriteString writes s to w with all replacements performed.

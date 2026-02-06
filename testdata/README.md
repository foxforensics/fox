# Fox Standard Test Pattern
Author: Christian Uhsat <christian@uhsat.de>

## Abstract
The Fox Standard Test Pattern (*FSTP*) was developed to provide a distinct pattern for debugging in text and hex mode. It has characterful technical and visual properties, combined with a small size of 64 bytes, that is still compressible.

## Definition
1. The marker `FOX` followed by a `0x0A` linebreak.
2. The pattern `0123` followed by a `0x20` space.
3. The character `U` repeated 47 times followed by a `0x20` space.
4. The count `47` followed by an `!` exclamation mark, then a `0x0A` linebreak.
5. The marker `XOF` marking the end of file.

## Annex I: Example as regular expression
```
FOX\n0123 U{47} 47!\nXOF
```

## Annex II: Example in ASCII text format
```
FOX
0123 UUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUU 47!
XOF
```

## Annex III: Example in canonical hex format
```
00000000  46 4f 58 0a 30 31 32 33  20 55 55 55 55 55 55 55  |FOX.0123 UUUUUUU|
00000010  55 55 55 55 55 55 55 55  55 55 55 55 55 55 55 55  |UUUUUUUUUUUUUUUU|
00000020  55 55 55 55 55 55 55 55  55 55 55 55 55 55 55 55  |UUUUUUUUUUUUUUUU|
00000030  55 55 55 55 55 55 55 55  20 34 37 21 0a 58 4f 46  |UUUUUUUU 47!.XOF|
```
---
*December 2025*
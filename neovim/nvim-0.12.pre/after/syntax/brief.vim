scriptencoding utf-8

if exists("b:current_syntax")
  finish
endif

let b:current_syntax = 'brief'

syntax keyword briefType int string bool contained
syntax keyword briefTypeKeyword idempotent rpc message
syntax keyword briefConst service const
syntax match briefComment /`[^`]\+`/
syntax match briefString /"[^"]\+"/ contained
syntax match briefMethod /\<[a-zA-Z]\+\(\s*(\)\@=/
syntax match briefType /\<[a-zA-Z]\+\(\s*{\)\@=/
syntax region briefDef start=/{/ end=/}/ fold contains=briefType,briefString

highlight default link briefTypeKeyword Keyword
highlight default link briefConst Constant
highlight default link briefComment Comment
highlight default link briefString String
highlight default link briefMethod Function
highlight default link briefArg Type
highlight default link briefType Type

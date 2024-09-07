function SetupBrief()
  if exists('b:brief_loaded')
    finish
  endif
  let b:brief_loaded = 1

  setlocal filetype=brief
  "syntax keyword briefKeyword rpc nextgroup=briefMethodName skipwhite
  syntax match briefKeyword "\<rpc\>" contains=@NoSpell skipwhite nextgroup=briefMethodName
  syntax match briefKeyword "\<bool\>" contains=@NoSpell
  syntax keyword briefKeyword message nextgroup=briefMessageName skipwhite
  syntax keyword briefKeyword const nextgroup=briefMessageName skipwhite
  syntax keyword briefKeyword service nextgroup=briefServiceName skipwhite
  syntax keyword briefKeyword int string
  syntax keyword briefKeyword idempotent
  syntax match briefServiceName "\"\S\+\""hs=s+1,he=e-1
  syntax match briefMessageName "\S\+" contained
  syntax match briefMethodName "\S\+" contained
  syntax region briefComment start="`"hs=e+1 end="`"he=s-1
  syntax match briefComment "//.*$"

  highlight briefKeyword ctermfg=Blue
  highlight briefComment ctermfg=Green cterm=italic
  highlight briefMethodName ctermfg=Yellow
  highlight briefMessageName ctermfg=Cyan
  highlight briefServiceName ctermfg=Red cterm=bold,italic
endfunction

autocmd BufRead,BufNewFile *.brief call SetupBrief()

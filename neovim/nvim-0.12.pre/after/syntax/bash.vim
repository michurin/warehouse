syntax region xBashJsonRegin
      \ start=/{/
      \ end=/}/
      \ contained
      \ containedin=shSingleQuote
      \ contains=xBashJsonKey,xBashJsonPunctuation
      \ keepend

syntax match xBashJsonKey /"\%([^"\\]\|\\.\)*"\ze:/ contained
syntax match xBashJsonPunctuation /[{},:\[\]]/ contained

syntax match xBashSqlComment /--.*/ containedin=shSingleQuote
syntax region xBashSqlCommentBlock start=/\/\*/ end=/\*\// containedin=shSingleQuote
syntax match xBashCppComment /\/\/.*/ containedin=shSingleQuote
syntax match xBashBashComment /#.*/ containedin=shSingleQuote

highlight xBashJsonKey ctermfg=green guifg=green
highlight xBashJsonPunctuation ctermfg=lightgreen guifg=lightgreen
highlight xBashSqlComment ctermfg=red guifg=red
highlight xBashSqlCommentBlock ctermfg=red guifg=red
highlight xBashCppComment ctermfg=lightgreen guifg=lightgreen
highlight xBashBashComment ctermfg=lightgreen guifg=lightgreen

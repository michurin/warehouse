" IDEAS: set up QuickFixLine (:h hl-QuickFixLine)

" contained force rule to match only after nextgroup trigger
syntax match MyQfFile /^[!-~]\+/ nextgroup=MyQfSpace
syntax match MyQfSpace /·\+/ contained nextgroup=MyQfLineNr skipwhite
syntax match MyQfLineNr /\d\+/ contained nextgroup=MyQfExample
syntax match MyQfExample /.\+/ contained

syntax match MyQfFileTest /_\zstest\ze.go/ contained containedin=MyQfFile
syntax match MyQfStatement /\<func\>/ contained containedin=MyQfExample

highlight link MyQfFile Comment
highlight link MyQfFileTest Statement
highlight link MyQfStatement Statement
" highlight link MyQfExample QuickFixLine " debugging only
highlight link MyQfSpace Whitespace
highlight link MyQfLineNr LineNr
" highlight MyTest guifg=#ff0000 gui=bold " custom

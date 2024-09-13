# to be included to bashrc

cd() {
    if test "$#" -eq 0
    then
        s="$(ls -1a | fzf --height 40% --reverse)"
        if test -d "$s"
        then
            builtin cd "$s"
        elif test -f "$s"
        then
            vim "$s"
        fi
    else
        builtin cd "$@"
    fi
}

#!/bin/bash

LOG=logs/log.log # /dev/null

FROM="$tg_message_from_id"

API() {
    API_STDOUT "$@" >>"$LOG"
}

API_STDOUT() {
    url="http://localhost$tg_x_ctrl_addr/$1"
    shift
    echo "====== curl $url $@" >>"$LOG"
    curl -qs "$url" "$@" 2>>"$LOG"
    echo >>"$LOG"
    echo >>"$LOG"
}

(
    echo '==================='
    echo "Args: $@"
    echo "Environment:"
    env | grep tg_ | sort
    echo '...................'
) >>"$LOG"

case "$1" in
    debug)
        echo '%!PRE'
        echo "Args: $@"
        echo "Environment:"
        env | grep tg_ | sort
        echo "FROM=$FROM"
        echo "LOG=$LOG"
        ;;
    about)
        echo '%!PRE'
        API_STDOUT getMe | jq
        ;;
    two)
        API "?to=$FROM" -d 'OK ONE!'
        API "?to=$FROM" -d 'OK TWO!!'
        echo 'OK NATIVE'
        ;;
    buttons)
        bGoogle='{"text":"Google","url":"https://www.google.com/"}'
        bDuck='{"text":"DuckDuckGo","url":"https://duckduckgo.com/"}'
        API sendMessage \
            -F chat_id=$FROM \
            -F text='Select search engine' \
            -F reply_markup='{"inline_keyboard":[['"$bGoogle,$bDuck"']]}'
        ;;
    image)
        curl -qs https://github.com/fluidicon.png
        ;;
    invert)
        wm=0
        fid=''
        for x in $tg_message_photo # finding the biggest image but ignoring too big ones
        do
            v=${x}_file_size
            s=${!v} # trick: getting variable name from variable; we need bash for it
            if test $s -gt 102400; then continue; fi # skipping too big files
            v=${x}_width
            w=${!v}
            v=${x}_file_id
            f=${!v}
            if test $w -gt $wm; then wm=$w; fid=$f; fi
        done
        if test -n "$fid"
        then
            API_STDOUT '' -G --data-urlencode "file_id=$fid" -o - | mogrify -flip -flop -format png -
        else
            echo "attache not found (maybe it was skipped due to enormous size)"
        fi
        ;;
    reaction)
        API setMessageReaction \
            -F chat_id=$FROM \
            -F message_id=$tg_message_message_id \
            -F reaction='[{"type":"emoji","emoji":"👾"}]'
        echo 'Bot reacted to your message☝️'
        ;;
    madrid)
        API sendLocation \
            -F chat_id="$FROM" \
            -F latitude='40.423467' \
            -F longitude='-3.712184'
        ;;
    menu)
        mShowEnv='{"text":"show environment","callback_data":"menu-debug"}'
        mShowNotification='{"text":"show notification","callback_data":"menu-notification"}'
        mShowAlert='{"text":"show alert","callback_data":"menu-alert"}'
        mLikeIt='{"text":"like it","callback_data":"menu-like"}'
        mUnlikeIt='{"text":"unlike it","callback_data":"menu-unlike"}'
        mDelete='{"text":"delete this message","callback_data":"menu-delete"}'
        mLayout="[[$mShowEnv],[$mShowAlert,$mShowNotification],[$mLikeIt,$mUnlikeIt],[$mDelete]]"
        API sendMessage \
            -F chat_id=$FROM \
            -F text='Actions' \
            -F reply_markup='{"inline_keyboard":'"$mLayout"'}'
        ;;
    run)
        API "?to=$FROM&a=reactions&a=$tg_message_message_id" -X RUN
        echo "Let me show you long run ☝️"
        ;;
    edit)
        API "?to=$FROM&a=editing" -X RUN
        ;;
    progress)
        API "?to=$FROM&a=progress" -X RUN
        ;;
    id)
        echo '%!PRE'
        id 2>&1
        ;;
    caps)
        echo '%!PRE'
        getpcaps --verbose --iab $$
        ;;
    hostname)
        echo '%!PRE'
        hostname 2>&1
        ;;
    help)
        API sendMessage -F chat_id=$FROM -F parse_mode=Markdown -F text='
Known commands:

- `debug` — show args, environment and vars
- `about` — reslut of getMe
- `two` — one request, two responses
- `buttons` — message with buttons
- `image` — show image
- `invert` (as capture to image) — returns flipped flopped image
- `reaction` — show reaction
- `madrid` — show location
- `menu` — scripted buttons
- `run` — long-run example (long sequence of reactions)
- `edit` — long-run example (editing)
- `progress` — one more long-run example (editing)
- `id` — check user who script runs from
- `caps` — check current capabilities (`getpcaps $$`)
- `hostname` — check hostname where script runs
- `help` — show this message
- `privacy` — mandatory privacy information
- `start` — just very first greeting message
'
        ;;
    start)
        API sendMessage -F chat_id=$FROM -F parse_mode=Markdown -F text='
Hi there!👋
It is demo bot to show an example of usage [cnbot](https://github.com/michurin/cnbot) bot engine.
You can use `help` command to see all available commands.'
        ;;
    privacy) # https://telegram.org/tos/bot-developers#4-privacy
        echo "This bot does not collect or share any personal information."
        ;;
    *)
        if test -n "$tg_callback_query_data"
        then
            case "$1" in
                menu-debug)
                    API answerCallbackQuery -F callback_query_id="$tg_callback_query_id"
                    echo '%!PRE'
                    echo "Environment:"
                    env | grep tg_ | sort
                    ;;
                menu-like)
                    API answerCallbackQuery -F callback_query_id="$tg_callback_query_id" -F "text=Like it"
                    API setMessageReaction -F chat_id=$tg_callback_query_message_chat_id \
                        -F message_id=$tg_callback_query_message_message_id \
                        -F reaction='[{"type":"emoji","emoji":"👾"}]'
                    ;;
                menu-unlike)
                    API answerCallbackQuery -F callback_query_id="$tg_callback_query_id" -F "text=Don't like it"
                    API setMessageReaction -F chat_id=$tg_callback_query_message_chat_id \
                        -F message_id=$tg_callback_query_message_message_id \
                        -F reaction='[]'
                    ;;
                menu-delete)
                    API answerCallbackQuery -F callback_query_id="$tg_callback_query_id"
                    API deleteMessage -F chat_id=$tg_callback_query_message_chat_id \
                        -F message_id=$tg_callback_query_message_message_id
                    ;;
                menu-notification)
                    API answerCallbackQuery -F callback_query_id="$tg_callback_query_id" -F text="Notification text (200 chars maximum)"
                    ;;
                menu-alert)
                    API answerCallbackQuery -F callback_query_id="$tg_callback_query_id" -F text="Notification text shown as alert" -F show_alert=true
                    ;;
            esac
        else
            API sendMessage -F chat_id=$FROM -F text='Invalid command. Say `help`.' -F parse_mode=Markdown
        fi
        ;;
esac

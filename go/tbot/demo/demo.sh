#!/bin/sh

# It is demo UNIX-shell script.
# You are free to use any language to express your ideas.

# Structure of file. [TODO]

# -----------------------------------------------------------------------
# Preamble. You can skip it and continue from PART I.
# -----------------------------------------------------------------------

# -- initialization -----------------------------------------------------

# It is good place to set up $PATH, $CURL_HOME etc

# -- current working directory ------------------------------------------
# you have to do it before any dowloading

x_script="$(realpath -- "$0")" # save real path for help generating
cd "$(dirname "$0")" || exit 1
mkdir -p tmp || exit 1
cd tmp || exit 1

# -- check binaries -----------------------------------------------------

for cmd in 'curl' 'date' 'dirname' 'echo' 'env' 'sort' 'uptime'
do
    if ! command -v "$cmd" >/dev/null
    then
        echo "This script ('$0') require '$cmd' command"
        exit
    fi
done

# -- prepare images -----------------------------------------------------

img_file_block=tmp_img_1.png
img_file_blue=tmp_img_2.png
img_file_fuchsia=tmp_img_3.png

[ -e "$img_file_block" ] || curl -qso "$img_file_block" https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_Black.png >&2
[ -e "$img_file_blue" ] || curl -qso "$img_file_blue" https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_Blue.png >&2
[ -e "$img_file_fuchsia" ] || curl -qso "$img_file_fuchsia" https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_Fuchsia.png >&2


# -----------------------------------------------------------------------
# PART I: What we can do without API
# -----------------------------------------------------------------------

# -- simplest: show standard output -------------------------------------

# Try to send to bot messages
# echo
# /echo
# ECHO ME...
# You will receive "Echo!"
if [ "$1" = 'echo' ]
then
    echo "Echo!"
    exit
fi

# -- simple formatted text ----------------------------------------------

# Try command "env" to see all environment variables
# You cal also send an image with caption "env" (literally)
# to see corresponding environment
if [ "$1" = 'env' -o "$tg_message_caption" = 'env' ]
then
    echo '%!PRE'
    env | sort
    exit
fi

# Try
# args
# args one two tree
if [ "$1" = 'args' ]
then
    echo '%!PRE'
    i=0
    for x in "$@"
    do
        i=$(($i + 1))
        echo "arg$i=$x"
    done
    exit
fi

# -- simple images ------------------------------------------------------

# You can just cat any image, audio, video, document (like PDF). The bot
# detects format and send content in reply in proper way.
if [ "$1" == 'img' ]
then
    cat "$img_file_fuchsia"
    exit
fi

# -----------------------------------------------------------------------
# PART II: What we can do with simple API
# -----------------------------------------------------------------------

# -- simplest call for control url --------------------------------------

CTRL="http://localhost$tg_x_ctrl_addr"

# sending without API engagement
# TODO note: url prefix
# TODO note: redirect output
# TODO note: beware telegram do not like such frequent messaging
if [ "$1" = 'send' ]
then
    curl -qs "$CTRL/x/?to=$tg_message_from_id" -d 'Text message' >&2
    curl -qs "$CTRL/x/?to=$tg_message_from_id" --data-binary @"$img_file_block" >&2
    exit
fi

# -----------------------------------------------------------------------
# PART III: What we can do with full Telegram API
# -----------------------------------------------------------------------

# -- rich formatting ----------------------------------------------------

# It is important to suppress output of curl
if [ "$1" = 'md' ]
then
    curl -qs "$CTRL/x/sendMessage" \
    -F chat_id="$tg_message_from_id" \
    -F text='Markdown: ||secret||, `pre`, *bold*, _italic_, ~cross~, __underline__, [inline link](https://core.telegram.org/bots/api#markdownv2-style)' \
    -F parse_mode='MarkdownV2' \
    -F disable_web_page_preview=true >&2
    exit
fi

# -- replaying ----------------------------------------------------------

# Try
# re one two
if [ "$1" = 're' ]
then
    curl -qs "$CTRL/x/sendMessage" -F chat_id="$tg_message_from_id" -F reply_to_message_id="$tg_message_message_id" -F text="My reply to '$*'" >&2
    exit
fi

# -----------------------------------------------------------------------
# PART IV: What we can do with API, advanced examples
# -----------------------------------------------------------------------

# -- inline keyboard ----------------------------------------------------

if [ "$1" = 'kbd' ]
then
    curl -qs "$CTRL/x/sendMessage" \
    -F chat_id="$tg_message_from_id" \
    -F text='Keyboard example \(`DELETE` will delete this message and show modal popup\)' \
    -F parse_mode='MarkdownV2' \
    -F reply_markup='{"inline_keyboard":[[{"text":"DATE", "callback_data":"kbd_date"}, {"text":"UPTIME", "callback_data":"kbd_uptime"}], [{"text":"DELITE", "callback_data":"kbd_delete"}]]}' >&2
    exit
fi

if [ "$tg_callback_query_data" = 'kbd_date' ]
then
    curl -qs "$CTRL/x/answerCallbackQuery" -F callback_query_id="$tg_callback_query_id" -F text="I'm going to show you current date" >&2
    date
    exit
fi

if [ "$tg_callback_query_data" = 'kbd_uptime' ]
then
    curl -qs "$CTRL/x/answerCallbackQuery" -F callback_query_id="$tg_callback_query_id" -F text="I'm going to show you uptime" >&2
    uptime
    exit
fi

# TODO comment: how to delete message
# TODO comment: where is userID
if [ "$tg_callback_query_data" = 'kbd_delete' ]
then
    curl -qs "$CTRL/x/answerCallbackQuery" -F callback_query_id="$tg_callback_query_id" -F text="I'm deleting message" -F show_alert=true >&2
    curl -qs "$CTRL/x/deleteMessage" \
    -F chat_id="$tg_callback_query_from_id" \
    -F message_id="$tg_callback_query_message_message_id" >&2
    exit
fi

# -- rich images --------------------------------------------------------

# rich caption. brackets escaped
if [ "$1" = 'photo' ]
then
    curl -qs "$CTRL/x/sendPhoto" \
    -F chat_id="$tg_message_from_id" \
    -F photo=@"$img_file_block" \
    -F caption='the image \(caption with ||rich format||\)' \
    -F parse_mode='MarkdownV2' >&2
    exit
fi

# -- group of media -----------------------------------------------------

if [ "$1" = 'group' ]
then
    curl -qs "$CTRL/x/sendMediaGroup" \
    -F chat_id="$tg_message_from_id" \
    -F media='[{"type":"photo","media":"attach://file_1"}, {"type":"photo","media":"attach://file_2"}]' \
    -F file_1=@"$img_file_block" \
    -F file_2=@"$img_file_blue" >&2
    exit
fi

# -- download (obtain from user), process and send back images ----------

if [ -n "$tg_message_photo" ] # list of images is not empty
then
    # TODO check convert
    filename_in="$tg_message_photo_0_file_id.img"
    filename_out="$tg_message_photo_0_file_id.jpg"
    curl -qs "$CTRL/x" -G --data-urlencode "file_id=$tg_message_photo_0_file_id" -o "$filename_in" >&2
    convert "$filename_in" -channel RGB -negate "png:$filename_out"
    cat "$filename_out"
    exit
fi

# -- send location ------------------------------------------------------

if [ "$1" = 'london' ]
then
    curl -qs "$CTRL/x/sendLocation" -F chat_id="$tg_message_from_id" -F latitude='51.5' -F longitude='-0.125' >&2
    exit
fi

if [ -n "$tg_message_location_longitude" ] # someone shares location
then
    # TODO check perl
    curl -qs "$CTRL/x/sendLocation" \
        -F chat_id="$tg_message_from_id" \
        -F latitude="$(perl -e 'print(- '"$tg_message_location_latitude"')')" \
        -F longitude="$(perl -e 'print(- '"$tg_message_location_longitude"')')" >&2
    echo 'The opposite corner of the Earth'
    exit
fi

# -- manage bot itself --------------------------------------------------

# note: use text='<-' to vanish filename-part of Content-Disposition header, Telegram API likes that
if [ "$1" = 'menu' ]
then
    # sending warning
    (
        echo '```'
        echo '**************'
        echo '*  CAUTION!  *'
        echo '**************'
        echo '```'
        echo 'You need to restart some Telegram clients to see this menu update'
    ) | curl -qs "$CTRL/x/sendMessage" -F chat_id="$tg_message_from_id" -F text='<-' -F parse_mode='MarkdownV2' >&2
    # check if menu present
    s="$(curl -qs "$CTRL/x/getMyCommands")" # {"ok":true,"result":[]}
    # toggle menu
    if [ -z ${s##*[]*} ] # if $s contains '[]'
    then
        b="$(perl -ne 'if (/"\$1"\s*=\s*'"'"'([^'"'"']+)'"'"'/) {print(qq|${sep}{"command":"/$1","description":"$1"}|); $sep=", ";}' "$x_script")"
        curl -qs "$CTRL/x/setMyCommands" -F commands='['"$b"']' >&2
        echo "Menu added"
    else
        curl -qs "$CTRL/x/deleteMyCommands" >&2
        echo "Menu deleted"
    fi
    exit
fi

# -----------------------------------------------------------------------
# PART V: Long-running routines
# -----------------------------------------------------------------------

if [ "$1" = 'long' ]
then
    curl -qs -X RUN "$CTRL/x/?to=$tg_message_from_id"
    exit
fi

if [ "$1" = 'cd' ]
then
    curl -qs -X RUN "$CTRL/x/?to=$tg_message_from_id&a=countdown"
    exit
fi

# -----------------------------------------------------------------------
# Appendix: Helpers
# -----------------------------------------------------------------------

# -- fallback: unrecognized command -------------------------------------

if [ "$1" = 'help' ]
then
    if command -v perl >/dev/null
    then
        echo 'Available commands:'
        perl -ne 'if (/"\$1"\s*=\s*'"'"'([^'"'"']+)'"'"'/) {print("$sep$1"); $sep=", ";}' <"$x_script"
    else
        echo 'This command require perl'
    fi
    exit
fi

echo "Unknown command '$1'; say 'help' to see all available commands"

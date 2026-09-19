#!/usr/bin/env python3

'''
echo 'cookie-value' | ./cookie-value-to-cookie-jar-file.py --domain one.com
$ echo 'HSID=ArB; SSID=Aws; APISID=55F' | ./cookie-value-to-cookie-jar-file.py --domain one.com
# Netscape HTTP Cookie File
# https://curl.se/docs/http-cookies.html
one.com FALSE   /       FALSE   0       HSID    ArB
one.com FALSE   /       FALSE   0       SSID    Aws
one.com FALSE   /       FALSE   0       APISID  55F

useful command:
$ pbpaste | ./cookie-value-to-cookie-jar-file.py --domain one.com
'''


import argparse
import sys


def main():
    parser = argparse.ArgumentParser(
        description="Convert a Cookie header value to a Netscape cookie jar."
    )
    parser.add_argument(
        "--domain",
        required=True,
        help="Cookie domain, e.g. example.com",
    )
    args = parser.parse_args()

    cookie_header = sys.stdin.read().strip()

    if not cookie_header:
        print("Error: empty Cookie header", file=sys.stderr)
        sys.exit(1)

    domain = args.domain.lstrip(".")

    if not domain:
        print("Error: empty domain", file=sys.stderr)
        sys.exit(1)

    print("# Netscape HTTP Cookie File")
    print("# https://curl.se/docs/http-cookies.html")

    count = 0

    for item in cookie_header.split(";"):
        item = item.strip()

        if not item:
            continue

        if "=" not in item:
            print(f"Warning: skipping invalid cookie: {item!r}", file=sys.stderr)
            continue

        name, value = item.split("=", 1)
        name = name.strip()
        value = value.strip()

        if not name:
            print("Warning: skipping cookie with empty name", file=sys.stderr)
            continue

        print(f"{domain}\tFALSE\t/\tFALSE\t0\t{name}\t{value}")
        count += 1

    if count == 0:
        print("Error: no valid cookies found", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    try:
        main()
    except BrokenPipeError:
        sys.exit(0)
    except Exception as e:
        print(f"error: {e}", file=sys.stderr)
        sys.exit(1)

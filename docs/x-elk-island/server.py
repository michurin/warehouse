#!/usr/bin/env python3

import sys
from http import server

class MyHTTPRequestHandler(server.SimpleHTTPRequestHandler):
    def end_headers(self):
        self.send_header("Cache-Control", "no-cache, no-store, must-revalidate")
        self.send_header("Pragma", "no-cache")
        self.send_header("Expires", "0")
        server.SimpleHTTPRequestHandler.end_headers(self)

if __name__ == '__main__':
    if len(sys.argv) == 2:
        port = int(sys.argv[1])
    else:
        port = 9999
    print(f'Starting at :{port}')
    server.HTTPServer(('', port), MyHTTPRequestHandler).serve_forever()

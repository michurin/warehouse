#!/usr/bin/env python3
from http import server

class MyHTTPRequestHandler(server.SimpleHTTPRequestHandler):
    def end_headers(self):
        self.send_header("Cache-Control", "no-cache, no-store, must-revalidate")
        self.send_header("Pragma", "no-cache")
        self.send_header("Expires", "0")
        server.SimpleHTTPRequestHandler.end_headers(self)

if __name__ == '__main__':
    server.HTTPServer(('', 9999), MyHTTPRequestHandler).serve_forever()

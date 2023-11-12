#!/usr/bin/env python3

import sys
from http import server

class MyHTTPRequestHandler(server.SimpleHTTPRequestHandler):
    def end_headers(self):
        self.send_header('Cache-Control', 'no-cache, no-store, must-revalidate')
        self.send_header('Pragma', 'no-cache')
        self.send_header('Expires', '0')
        super().end_headers()
    def log_message(self, format, *args):
        message = format % args
        sys.stderr.write('%s - %s [%s] %s\n' % (
            self.address_string(),
            self.headers.get('user-agent', '-'),
            self.log_date_time_string(),
            message.translate(self._control_char_table)))

if __name__ == '__main__':
    if len(sys.argv) == 2:
        port = int(sys.argv[1])
    else:
        port = 9999
    print(f'Starting at :{port}')
    try:
        server.HTTPServer(('', port), MyHTTPRequestHandler).serve_forever()
    except KeyboardInterrupt:
        print('') # New line after ^C

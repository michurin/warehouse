#!/bin/bash

export _nginxver=$(nginx -v 2>&1 | sed 's-.*/--')
read -p "nginx version=$_nginxver. is it ok? [y/n] " yn
case $yn in
	[Yy]*) ;;
	*) exit;;
esac
echo 'build...'

set -e
curl -o ngx.tar.gz https://nginx.org/download/nginx-$_nginxver.tar.gz
tar xzf ngx.tar.gz
cd nginx-$_nginxver
./configure --with-http_xslt_module=dynamic --with-compat
make modules
cd ..

echo 'RESULT:'
find nginx-$_nginxver/objs -type f -name '*.so'
echo '"cp" or "ln -s" it to . (/etc/nginx)'
echo ''
echo 'sudo systemctl restart nginx && echo OK && sudo systemctl status nginx'

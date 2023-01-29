#!/bin/bash

export _nginxver=$(nginx -v 2>&1 | sed 's-.*/--')
read -p "nginx version=$_nginxver. is it ok? " yn
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

echo 'RESULT:'
find "$(pwd)" -type f -name 'ngx_http_xslt_filter_module.so'
echo 'cp it to /usr/lib/nginx/modules/'
echo 'or       /etc/nginx'

#!/bin/sh

sed -i "s|%API_BASE_URL%|$REACT_API_BASE_URL|g" /usr/share/nginx/html/index.html

exec nginx -g "daemon off;"

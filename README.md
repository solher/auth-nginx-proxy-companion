# nginx-proxy-auth-companion
A simple auth companion for Nginx proxy.

# Basic usage
The simplest way to use it is by running it along my [customized Nginx proxy](https://github.com/solher/nginx-proxy) that has built-in support.

By default, it denies access to everybody. To can enable the development mode by setting the env variable `GRANT_ALL` to `true`.
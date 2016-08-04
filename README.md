# nginx-proxy-auth-companion
A simple auth companion for Nginx proxy.

# Basic usage
The simplest way to use it is by running it along my [customized Nginx proxy](https://github.com/solher/nginx-proxy) that has built-in support.

An example full Nginx proxy setup can be found [here](https://github.com/solher/compose-nginx-proxy).

By default, it denies all access to everybody. To can enable the development mode by setting the env variable `GRANT_ALL` to `true`.
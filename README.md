# loginsrv

**NOTE**: This fork does not aim to be up-to-date with upstream anymore. For more details, see below.

loginsrv is a standalone minimalistic login server providing a [JWT](https://jwt.io/) login for multiple login backends.

[![Docker](https://img.shields.io/docker/pulls/kernle32dll/loginsrv.svg)](https://hub.docker.com/r/kernle32dll/loginsrv/)
[![Build Status](https://github.com/kernle32dll/loginsrv/workflows/test/badge.svg)](https://github.com/kernle32dll/loginsrv/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/kernle32dll/loginsrv)](https://goreportcard.com/report/github.com/kernle32dll/loginsrv)
[![Coverage Status](https://coveralls.io/repos/github/kernle32dll/loginsrv/badge.svg?branch=master)](https://coveralls.io/github/kernle32dll/loginsrv?branch=master)

## Abstract

Loginsrv provides a minimal endpoint for authentication. The login is performed against the providers and returned as a JSON Web Token (JWT).
It can be used as:

* Standalone microservice
* Docker container

![](.screenshot.png)

## Supported Provider Backends
The following providers (login backends) are supported.

* [Htpasswd](#htpasswd)
* [Simple](#simple) (user/password pairs by configuration)
* [Httpupstream](#httpupstream)
* [OAuth2](#oauth2)
  * GitHub login
  * Google login
  * Bitbucket login
  * Facebook login
  * Gitlab login

## Difference to tarent/loginsrv

This fork is a re-start of the project, tailored to the authors needs. This fork is in no way supported by tarent!

Version 1.4.0 is a trimmed down version of the latest master of tarent/loginsrv at the time of forking. Trimmed down means:

- Removed support for Caddy (only v1 was supported, and it was breaking with dependency updates)
- Removed support for OSIAM (unmaintained since 2019)
- Removed support for correlation id
- Replacement of [dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go) with [golang-jwt/jwt](https://github.com/golang-jwt/jwt)
- Go 1.18 baseline
- Up-to-date dependencies

### Future

There are at least two releases after 1.4.0 planned:

#### 1.5.0 A maintenance release

Fixing code problems, adding missing tests, and other quality-of-life improvements. Logging will be migrated from
[sirupsen/logrus](https://github.com/sirupsen/logrus) to [rs/zerolog](https://github.com/rs/zerolog). Existing
backends will be evaluated and tested to see if they are actually still working.

#### 1.6.0 A feature release

Planned features include more quality of life improvements, such as adding support for tracing via
[open-telemetry](https://github.com/open-telemetry/opentelemetry-go) (as a replacement for the removed correlation id),
as well as adding new features which did not make it into upstream loginsrv.

## Configuration and Startup
### Config Options

| Parameter                 | Type        | Default      | Description                                                                                           |
|---------------------------|-------------|--------------|-------------------------------------------------------------------------------------------------------|
| -cookie-domain            | string      |              | Optional domain parameter for the cookie                                                              |
| -cookie-expiry            | string      | session      | Expiry duration for the cookie, e.g. 2h or 3h30m                                                      |
| -cookie-http-only         | boolean     | true         | Set the cookie with the HTTP only flag                                                                |
| -cookie-name              | string      | "jwt_token"  | Name of the JWT cookie                                                                                |
| -cookie-secure            | boolean     | true         | Set the secure flag on the JWT cookie. (Set this to false for plain HTTP support)                     |
| -github                   | value       |              | OAuth config in the form: client_id=..,client_secret=..[,scope=..][,redirect_uri=..]                  |
| -google                   | value       |              | OAuth config in the form: client_id=..,client_secret=..[,scope=..][,redirect_uri=..]                  |
| -bitbucket                | value       |              | OAuth config in the form: client_id=..,client_secret=..[,scope=..][,redirect_uri=..]                  |
| -facebook                 | value       |              | OAuth config in the form: client_id=..,client_secret=..[,scope=..][,redirect_uri=..]                  |
| -gitlab                   | value       |              | OAuth config in the form: client_id=..,client_secret=..[,scope=..,][redirect_uri=..]                  |
| -host                     | string      | "localhost"  | Host to listen on                                                                                     |
| -htpasswd                 | value       |              | Htpasswd login backend opts: file=/path/to/pwdfile                                                    |
| -jwt-expiry               | go duration | 24h          | Expiry duration for the JWT token, e.g. 2h or 3h30m                                                   |
| -jwt-secret               | string      | "random key" | Secret used to sign the JWT token.                                                                    |
| -jwt-secret-file          | string      |              | File to load the jwt-secret from, e.g. `/run/secrets/some.key`. **Takes precedence over jwt-secret!** |
| -jwt-algo                 | string      | "HS512"      | Signing algorithm to use (ES256, ES384, ES512, RS256, RS384, RS512, HS256, HS384, HS512)              |
| -log-level                | string      | "info"       | Log level                                                                                             |
| -login-path               | string      | "/login"     | Path of the login resource                                                                            |
| -logout-url               | string      |              | URL or path to redirect to after logout                                                               |
| -port                     | string      | "6789"       | Port to listen on                                                                                     |
| -redirect                 | boolean     | true         | Allow dynamic overwriting of the the success by query parameter                                       |
| -redirect-query-parameter | string      | "backTo"     | URL parameter for the redirect target                                                                 |
| -redirect-check-referer   | boolean     | true         | Check the referer header to ensure it matches the host header on dynamic redirects                    |
| -redirect-host-file       | string      | ""           | A file containing a list of domains that redirects are allowed to, one domain per line                |
| -simple                   | value       |              | Simple login backend opts: user1=password,user2=password,..                                           |
| -success-url              | string      | "/"          | URL to redirect to after login                                                                        |
| -template                 | string      |              | An alternative template for the login form                                                            |
| -text-logging             | boolean     | true         | Log in text format instead of JSON                                                                    |
| -jwt-refreshes            | int         | 0            | The maximum number of JWT refreshes                                                                   |
| -grace-period             | go duration | 5s           | Duration to wait after SIGINT/SIGTERM for existing requests. No new requests are accepted.            |
| -user-file                | string      |              | A YAML file with user specific data for the tokens. (see below for an example)                        |
| -user-endpoint            | string      |              | URL of an endpoint providing user specific data for the tokens. (see below for an example)            |
| -user-endpoint-token      | string      |              | Authentication token used when communicating with the user endpoint                                   |
| -user-endpoint-timeout    | go duration | 5s           | Timeout used when communicating with the user endpoint                                                |

### Environment Variables
All of the above Config Options can also be applied as environment variables by using variables named this way: `LOGINSRV_OPTION_NAME`.
So e.g. `jwt-secret` can be set by environment variable `LOGINSRV_JWT_SECRET`.

### Startup Examples
The simplest way to use loginsrv is by the provided docker container.
E.g. configured with the simple provider:
```sh
$ docker run -d -p 8080:8080 kernle32dll/loginsrv -cookie-secure=false -jwt-secret my_secret -simple bob=secret

$ curl --data "username=bob&password=secret" 127.0.0.1:8080/login
eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJib2IifQ.uWoJkSXTLA_RvfLKe12pb4CyxQNxe5_Ovw-N5wfQwkzXz2enbhA9JZf8MmTp9n-TTDcWdY3Fd1SA72_M20G9lQ
```

The same configuration could be written with environment variables this way:
```sh
$ docker run -d -p 8080:8080 -E COOKIE_SECURE=false -e LOGINSRV_JWT_SECRET=my_secret -e LOGINSRV_BACKEND=provider=simple,bob=secret kernle32dll/loginsrv
```

## API

### GET /login

Per default, it returns a simple bootstrap styled login form for unauthenticated requests and a page with user info for authenticated requests.
When the call accepts a JSON output, the json content of the token is returned to authenticated requests.

| Parameter-Type | Parameter                | Description                                                  |         | 
|----------------|--------------------------|--------------------------------------------------------------|---------|
| Http-Header    | Accept: text/html        | Return the login form or user html.                          | default |
| Http-Header    | Accept: application/json | Return the user Object as json, or 403 if not authenticated. |         |

### GET /login/<provider>

Starts the OAuth Web Flow with the configured provider. E.g. `GET /login/github` redirects to the GitHub login form.

### POST /login

Performs the login and returns the JWT. Depending on the content-type and parameters, a classical JSON-Rest or a redirect can be performed.

#### Runtime Parameters

| Parameter-Type | Parameter                                       | Description                                                       |              | 
|----------------|-------------------------------------------------|-------------------------------------------------------------------|--------------|
| Http-Header    | Accept: text/html                               | Set the JWT as a cookie named 'jwt_token'                         | default      |
| Http-Header    | Accept: application/jwt                         | Returns the JWT within the body. No cookie is set                 |              |
| Http-Header    | Content-Type: application/x-www-form-urlencoded | Expect the credentials as form encoded parameters                 | default      |
| Http-Header    | Content-Type: application/json                  | Take the credentials from the provided JSON object                |              |
| Post-Parameter | username                                        | The username                                                      |              |
| Post-Parameter | password                                        | The password                                                      |              |
| Get or Post    | backTo                                          | Dynamic redirect target after login (see (Redirects)[#redirects]) | -success-url |

#### Possible Return Codes

| Code | Meaning               | Description                                                                                                               |
|------|-----------------------|---------------------------------------------------------------------------------------------------------------------------|
| 200  | OK                    | Successfully authenticated                                                                                                |
| 403  | Forbidden             | The credentials are wrong                                                                                                 |
| 400  | Bad Request           | Missing parameters                                                                                                        |
| 500  | Internal Server Error | Internal error, e.g. the login provider is not available or failed                                                        |
| 303  | See Other             | Sets the JWT as a cookie, if the login succeeds and redirect to the URLs provided in `redirectSuccess` or `redirectError` |

Hint: The status `401 Unauthorized` is not used as a return code to not conflict with an HTTP Basic authentication.

#### JWT-Refresh

If the POST-Parameters for username and password are missing and a valid JWT-Cookie is part of the request, then the JWT-Cookie is refreshed.
This only happens if the jwt-refreshes config option is set to a value greater than 0. 

### DELETE /login

Deletes the JWT cookie.

For simple usage in web applications, this can also be called by `GET|POST /login?logout=true`

### API Examples

#### Example:
Default is to return the token as Content-Type application/jwt within the body.
```sh
curl -i --data "username=bob&password=secret" http://127.0.0.1:6789/login
HTTP/1.1 200 OK
Content-Type: application/jwt
Date: Mon, 14 Nov 2016 21:35:42 GMT
Content-Length: 100

eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJib2IifQ.-51G5JQmpJleARHp8rIljBczPFanWT93d_N_7LQGUXU
```

#### Example: Credentials as JSON
The credentials can also be sent JSON encoded.
```sh
curl -i -H 'Content-Type: application/json'  --data '{"username": "bob", "password": "secret"}' http://127.0.0.1:6789/login
HTTP/1.1 200 OK
Content-Type: application/jwt
Date: Mon, 14 Nov 2016 21:35:42 GMT
Content-Length: 100

eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJib2IifQ.-51G5JQmpJleARHp8rIljBczPFanWT93d_N_7LQGUXU
```

#### Example: web based flow with 'Accept: text/html'
Sets the JWT as a cookie and redirects to a web page.
```sh
curl -i -H 'Accept: text/html' --data "username=bob&password=secret" http://127.0.0.1:6789/login
HTTP/1.1 303 See Other
Location: /
Set-Cookie: jwt_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJib2IifQ.-51G5JQmpJleARHp8rIljBczPFanWT93d_N_7LQGUXU; HttpOnly
```

#### Example: AJAX call with JQuery to fetch a JWT token and create a cookie from it
Creates a cookie from a successful API call to login.
```js
$.ajax({
	url: "http://localhost:8080/login",
	type: 'POST',
	dataType: 'text',
	contentType: 'application/json',
	data: JSON.stringify( { 
		'username': 'demo', 
		'password': 'demo'
	}),
	success: function(data) {
		document.cookie = "jwt_token=" + data + ";path=/";
	},
	error: function (xhr, ajaxOptions, thrownError) {
	}
});
```
Make sure your main page has JQuery:
```html
<script src="https://code.jquery.com/jquery-3.3.1.min.js"></script>
```

### Redirects

The API has support for a redirect query parameter, e.g. `?backTo=/dynamic/return/path`. For security reasons, the default behaviour is very restrictive:

* Only local redirects (same host) are allowed.
* The `Referer` header is checked to ensure that the call to the login page came from the same page.

These restrictions are there, to prevent you from unchecked redirect attacks, e.g. phishing or login attacks.
If you know, what you are doing, you can disable the `Referer` check with `--redirect-check-referer=false` and provide a whitelist file
for allowed external domains with `--redirect-host-file=/some/domains.txt`.

## The JWT Token
Depending on the provider, the token may look as follows:
```json
{
  "sub": "smancke",
  "picture": "https://avatars2.githubusercontent.com/u/4291379?v=3",
  "name": "Sebastian Mancke",
  "email": "s.mancke@kernle32dll.de",
  "origin": "github"
}
```

## Provider Backends

### Htpasswd
Authentication against htpasswd file. MD5, SHA1 and Bcrypt are supported. But we recommend to only use Bcrypt for security reasons (e.g. `htpasswd -B -C 15`).

Parameters for the provider:

| Parameter-Name | Description                                                                        |
|----------------|------------------------------------------------------------------------------------|
| file           | Path to the password file (multiple files can be used by separating them with ';') |

Example:
```sh
loginsrv -htpasswd file=users
```

### Httpupstream
Authentication against an upstream HTTP server by performing a HTTP Basic authentication request and checking the response for a HTTP 200 OK status code. Anything other than a 200 OK status code will result in a failure to authenticate.

Parameters for the provider:

| Parameter-Name | Description                                                               |
|----------------|---------------------------------------------------------------------------|
| upstream       | HTTP/HTTPS URL to call                                                    |
| skipverify     | True to ignore TLS errors (optional, false by default)                    |
| timeout        | Request timeout (optional 1m by default, go duration syntax is supported) |

Example:
```sh
loginsrv -httpupstream upstream=https://google.com,timeout=1s
```

### Simple
Simple is a demo provider for testing only. It holds a user/password table in memory.

Example
```sh
loginsrv -simple bob=secret
```

## OAuth2

The OAuth Web Flow (aka 3-legged-OAuth flow) is also supported.
Currently the following OAuth provider is supported:

* GitHub
* Google
* Bitbucket
* Facebook
* Gitlab

An OAuth provider supports the following parameters:

| Parameter-Name | Description                           |
|----------------|---------------------------------------|
| client_id      | OAuth Client ID                       |
| client_secret  | OAuth Client Secret                   |
| scope          | Space separated scope List (optional) |
| redirect_uri   | Alternative Redirect URI (optional)   |

When configuring the OAuth parameters at your external OAuth provider, a redirect URI has to be supplied. This redirect URI has to point to the path `/login/<provider>`.
If not supplied, the OAuth redirect URI is calculated out of the current URL. This should work in most cases and should even work
if loginsrv is routed through a reverse proxy, if the headers `X-Forwarded-Host` and `X-Forwarded-Proto` are set correctly.

### GitHub Startup Example
```sh
$ docker run -p 80:80 kernle32dll/loginsrv -github client_id=xxx,client_secret=yyy
```

## Templating

A custom template can be supplied by the parameter `template`. 
You can find the original template in [login/login_form.go](https://github.com/kernle32dll/loginsrv/blob/master/login/login_form.go).

The templating uses the Golang template package. A short intro can be found [here](https://astaxie.gitbooks.io/build-web-application-with-golang/en/07.4.html).

When you specify a custom template, only the layout of the original template is replaced. The partials of the original are still loaded into the template context and can be used by your template. So a minimal unstyled login template could look like this:

```html
<!DOCTYPE html>
<html>
  <head>
      <!-- your styles -->
  <head>
  <body>
      <!-- your header -->

      {{ if .Error}}
        <div class="alert alert-danger" role="alert">
          <strong>Internal Error. </strong> Please try again later.
        </div>
      {{end}}

      {{if .Authenticated}}

         {{template "userInfo" . }}

      {{else}}

        {{template "login" . }}

      {{end}}

      <!-- your footer -->
  </body>
</html>
```

## Custom claims

To customize the content of the JWT token either a file wich contains
user data or an endpoint providing claims can be provided.

### User file

A user file is a YAML file which contains additional information which
is encoded in the token. After successful authentication against a
backend system, the user is searched within the file and the content
of the claims parameter is used to enhance the user JWT claim
parameters.

To match an entry, the user file is searched in linear order and all attributes has to match
the data of the authentication backend. The first matching entry will be used and all parameters
below the claim attribute are written into the token. The following attributes can be used for matching:
* `sub` - the username (all backends)
* `origin` - the provider or backend name (all backends)
* `email` - the mail address (the OAuth provider)
* `domain` - the domain (Google only)
* `groups` - the full path string of user groups enclosed in an array (Gitlab only)

Example:
* The user bob will become the `"role": "superAdmin"`, when authenticating with htpasswd file
* The user admin@example.org will become `"role": "admin"` and `"projects": ["example"]`, when authenticating with Google OAuth
* All other Google users with the domain example will become `"role": "user"` and `"projects": ["example"]`
* All other Gitlab users with group `example/subgroup` and `othergroup` will become `"role": "admin"`.
* All others will become `"role": "unknown"`, independent of the authentication provider

```yaml
- sub: bob
  origin: htpasswd
  claims:
    role: superAdmin

- email: admin@example.org
  origin: Google
  claims:
    role: admin
    projects:
      - example

- domain: example.org
  origin: Google
  claims:
    role: user
    projects:
      - example

- groups:
    - example/subgroup
    - othergroup
  origin: gitlab
  claims:
    role: admin

- claims:
    role: unknown
```

### User endpoint

A user endpoint is a http endpoint which provides additional
information on an authenticated user. After successful authentication
against a backend system, the endpoint gets called and the provided
information is used to enhance the user JWT claim parameters.

loginsrv passes these parameters to the endpoint:
* `sub` - the username (all backends)
* `origin` - the provider or backend name (all backends)
* `email` - the mail address (the OAuth provider)
* `domain` - the domain (Google only)
* `groups` - the full path string of user groups enclosed in an array (Gitlab only)

An interaction looks like this

```http
GET /claims?origin=google&sub=test@example.com&email=test@example.com HTTP/1.1
Host: localhost:8080
Accept: */*
Authorization: Bearer token

HTTP/1.1 200 OK
Content-Type: application/json

{
  "sub":"test@example.com",
  "uid":"113",
  "origin":"google",
  "permissions": ["read", "write"]
}
```

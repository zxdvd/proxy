## proxy
simple proxy support tcp, udp and http connect

### Configuration
very easy toml style configuration, e.g.
```
[tcp]
  [tcp.t1]
  from = "localhost:9001"
  to = "www.example.com:443"

[udp]
  [udp.u1]
  from = "localhost:9101"
  to = "localhost:9102"

[http]
  [http.h1]
  listen = "localhost:9201"
    [[http.h1.basic_auth]]
    user = "user1"
    password = "password1"
    [[http.h1.basic_auth]]
    # can give password without user (user can be anything)
    password = "password2"
```

### Run
`go build && ./proxy`

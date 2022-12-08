# Red Team Tools

For e2e tests

```
./elasticsearch/e2e.sh
./redis/e2e.sh
```

# features to add

* add flag for number of workers
* keep track of stopping point of ranges scanned, so process can be resumed


### ideas for subdomain takeover
query cert.sh
get the values from matching identities
remove base domain
add to knock.py wordlist - https://github.com/guelfoweb/knock/blob/master/knockpy/wordlist.txt
add virus total results: https://github.com/guelfoweb/knock/blob/master/knockpy/knockpy.py#L100
do a dns lookup for each variant

```
❯ nslookup jobs.google.com
Server:		192.168.1.1
Address:	192.168.1.1#53

Non-authoritative answer:
jobs.google.com	canonical name = www3.l.google.com.
Name:	www3.l.google.com
Address: 142.250.191.206
```

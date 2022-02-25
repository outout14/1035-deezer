
#  1035-Deezer
  
_(This project is useless.)_

1035-Deezer serves as a DNS server which returns the current listening of a Deezer, a music streaming service, user.


## How it works 
Due to limitations from the Deezer API, only the **last played music** will be returned, not the current one.  
Once logged, the Oauth user token will be stored in a Redis Database later access.
When a DNS TXT query will be asked to the server, the software will lookup the Redis server, get the user token and perform an API call to the history endpoint of the Deezer API. 
A DNS response will be returned containing the last played title and it's performer.  

##  Utilization
In theses examples, ``.dz.bb0.nl`` is the base domain and ``399552895`` is the user Deezer ID we are looking for. 

### Login
If the user access token isn't registered in the Redis database, an Oauth authorization URL will be returned 

Example query : 
```
❯ dig 399552895.dz.bb0.nl  TXT +short
"Can't get this user playing song."
"User may not exist."
"If that's you, connect the app to your Deezer account :"
"https://connect.deezer.com/oauth/auth.php?app_id=529162&redirect_uri=https://1035.bb0.nl/callback&perms=listening_history,offline_access"
```

### Request user 
Example query : 
```
❯ dig 399552895.dz.bb0.nl  TXT +short
"Last played song : Kiminosei"
"Author : the peggies"
```
## Installation 
DEB and RPM package are generated on each release. 
The configuration file is stored in ``/etc/1035-deezer/config.json``. 
A systemd service is provided as ``1035-deezer.service``. 
Proxying the HTTP endpoint is recommanded for TLS support. 

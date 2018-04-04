# listen-to
a dinky commandline tool for recommending music on a shared server

## why
idk. music sharing and recommendations is so fancy now. there are industries built around it. there are companies that make a lot of cash while allowing you to recommend and be recommended music via, often, algorithmically opaque machine learning models. if you want to get some music recommendations which will quench your thirst for new tunes, like, asap and with minimal effort then you are definitely in the wrong place.

`listen-to` isn't fancy at all. its meant to be used on a server shared by you and your friends. its very lo-fi-tech; something i imagine someone would have made several decades ago.

## how
### make recommendations

``` shell
λ listen-to brigadier jerry - jamaica jamaica
```

this will append your recommendation into a centralized log file on the server.

``` shell
2018-04-03T20:43:00-04:00 c brigadier jerry - jamaica jamaica
```

### listen for recommendations

``` shell
λ listen-to subscribe
λ listen-to s  # short-hand version
```
this will start a blocking filesystem listener and print recommendations everytime someone makes a recommendation.

``` shell
c says you should listen to brigadier jerry - jamaica jamaica
```

also note that you can run the listener in the background via `listen-to s &` and it will not block but will print out in your terminal still.

### list all recommendations

``` shell
λ listen-to list
```
this will print out a list of all recommendations.

### epilogue: init
the previous instructions assumed that a central music log directory had already been setup on the server you are running `listen-to` on. By default this directory is set to `/etc/listen-to/` with a default music log file of `music.log`. However, this is configurable via an environment variable which can be set system-wide in `/etc/environment`,

``` shell
LISTEN_TO_LOG_FILE=/path/to/your/custom/music.log
```
note that you need the sysadmin to set this up for you if it hasn't been setup already.

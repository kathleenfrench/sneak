## sneak

`sneak` is a tool written for playing `hack the box`, though it's certainly adaptable to other platforms as well. 

### local db

`sneak` uses `bolthold` (which wraps `boltdb`) to manage data locally beyond the user's custom configurations at the application level. if you want to interact with the database, there's a hidden command `sneak db` that will enable you to view the bucket(s), reset the database, as well as back it up


#### running it in docker

```
	@docker run \
		-v $(CWD)/shared/:/home/$(USERNAME)/tools \
		-v $(CWD)/boxes/:/home/$(USERNAME)/boxes \
		--privileged \
		--sysctl net.ipv6.conf.all.disable_ipv6=0 \
		--env LOCAL_NETWORK=192.168.1.151 \
		-p 8118:8118 \
		-it kathleenfrench/honcho \
		/bin/bash
```

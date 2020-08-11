## sneak

`sneak` is a tool written for playing `hack the box`, though it's certainly adaptable to other platforms as well. 

### local db

`sneak` uses `bolthold` (which wraps `boltdb`) to manage data locally beyond the user's custom configurations at the application level. if you want to interact with the database, there's a hidden command (`sneak db`) that will enable you to view the bucket(s), reset the database, as well as back it up

#### running sneak in docker

create an `env.mk` file with your `hack the box` username:

```
## env.mk
HTB_USERNAME := yourusername
```

this value will be imported when you run `sneaker`. running `sneak`'s containerized environment, `sneaker`, is very simple:

```
	@docker run \
		--privileged \
		--sysctl net.ipv6.conf.all.disable_ipv6=0 \
		--env LOCAL_NETWORK=$(local_network) \
		--cap-add=NET_ADMIN \
		-p 8118:8118 \
		-it sneaker \
		 /bin/sh
```

##### mounting your local data

if you want to persist/mount data from your local installation of `sneak`, add volume flags to the above `docker run` command, i.e.:

```
	@docker run \
		--privileged \
		--sysctl net.ipv6.conf.all.disable_ipv6=0 \
		--env LOCAL_NETWORK=$(local_network) \
		--cap-add=NET_ADMIN \
		-v $(HOME)/.sneak/:/home/$(HTB_USERNAME)/.sneak \
		-v $(CWD)/build/sneak:/go/bin/sneak \
		-p 8118:8118 \
		-it sneaker \
		 /bin/sh
```

once you're in the container and want to run `sneak` just append the `--mount` (`-m`) flag a single time to update your config files to the correct path.

```
sneak --mount (-m)
```

if you want to switch back to running `sneak` outside of docker with the same persisted files, run:

```
sneak --unmount (-u)
```
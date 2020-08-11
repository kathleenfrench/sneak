## sneak

`sneak` is a command line tool written for playing `hack the box`, but the repository also includes an image for running sneak in a pre-baked containerized environment for `sneak` that enables the user to connect to the `htb` VPN and run a reverse proxy (via `privoxy`) in their `sneaker` container to access `htb` sites in their local browser without needing to run a full virtual machine.

### local db

`sneak` uses `bolthold` (which wraps `boltdb`) to manage data locally beyond the user's custom configurations at the application level. if you want to interact with the database, there's a hidden command (`sneak db`) that will enable you to view the bucket(s), reset the database, as well as back it up

#### running sneak in docker

running `sneak`'s containerized environment, `sneaker`, is very simple:

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

##### running the container as a custom user

if you want to set your `hack the box` username as the container user, create an `env.mk` file with the following:

```
## env.mk
HTB_USERNAME := yourusername
```

this value will be imported when you run build and run the `sneaker` image, which is simplest to do via the `Makefile`'s `up` command:

```
make up
```

which will handle building the image with your custom user info and starting the container


##### mounting your local data

if you want to persist/mount data from your local installation of `sneak`, add volume flags to the above `docker run` command, i.e.:

```
	@docker run \
		--privileged \
		--sysctl net.ipv6.conf.all.disable_ipv6=0 \
		--env LOCAL_NETWORK=$(local_network) \
		--cap-add=NET_ADMIN \
		-v (SEE BELOW)
		-v $(CWD)/build/sneak:/go/bin/sneak \
		-p 8118:8118 \
		-it sneaker \
		 /bin/sh
```

if you are using the **default** `sneaker` image (which is the `sneak` user), use the following `--volume` flag:

```
-v $(HOME)/.sneak/:/home/sneak/.sneak
```

if you are running the `sneaker` image as a **custom user** (with your `hack the box` username), use the following `--volume` flag:

```
-v $(HOME)/.sneak/:/home/$(HTB_USERNAME)/.sneak
```

once you're in the container and want to run `sneak` just append the `--mount` (`-m`) flag a single time to update your config files to the correct path.

```
sneak --mount (-m)
```

if you want to switch back to running `sneak` outside of docker with the same persisted files, run:

```
sneak --unmount (-u)
```

### connecting to the VPN

after your configs have been set, run:

```
sneak vpn setup
```

this will prompt you for your generated `.ovpn` file from `hack the box' and create your `privoxy` config file so you can connect to the container locally via reverse proxy.

to actually connect, run:

```
sneak vpn connect
```

this will start the `openvpn` client. you can always verify your connection with:

```
sneak vpn test
```

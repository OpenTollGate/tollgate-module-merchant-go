# Tollgate Module - merchant (go)

This Tollgate module is responsible for:

- Handling payments from customers
- Initiating session after payment
- Evaluating LAN sightings from crowsnest and making an economical decision on whether to connect to those spotted Tollgates.

# Compile for ATH79 (GL-AR300 NOR)

```bash
cd ./src
env GOOS=linux GOARCH=mips GOMIPS=softfloat go build -o merchant -trimpath -ldflags="-s -w"

# Hint: copy to connected router 
scp -O merchant root@192.168.1.1:/root/merchant
```

# Compile for GL-MT3000

## Build

```bash
cd ./src
env GOOS=linux GOARCH=arm64 go build -o merchant -trimpath -ldflags="-s -w"

# Hint: copy to connected router 
scp -O merchant root@192.168.1.1:/root/merchant # X.X == Router IP
```

## Required Firewall rules

First, test if the merchant is up by going to your router's ip on port `2122`. You should get a JSON response with your IP and mac address.

Add to `/etc/config/firewall`:
```uci
config rule
	option name 'Allow-merchant-In'
	option src 'lan'
	option proto 'tcp'
	option dest_port '2122' # merchant port
	option target 'ACCEPT'

config redirect
	option name 'TollGate - Nostr merchant DNAT'
	option src 'lan'
	option dest 'lan'
	option proto 'tcp'
	option src_dip '192.168.21.21'
	option src_dport '2121'
	option dest_ip '192.168.X.X' # Router IP
	option dest_port '2122' # merchant port
	option target 'DNAT'

config redirect
        option name 'TollGate - Nostr merchant DNAT port'
        option src 'lan'
        option dest 'lan'
        option proto 'tcp'
        option src_dip '192.168.X.X' # Router IP
        option src_dport '2121'
        option dest_ip '192.168.X.X' # Router IP
        option dest_port '2122' # merchant port
        option target 'DNAT'
```

Run `service firewall restart` to make changes go into effect.

To test the firewall rule, go to `192.168.21.21:2122`. You should be greeted with the same JSON.


## OpenNDS rules
**Prerequisite: OpenNDS is installed**

To allow unauthenticated clients to reach the merchant, we need to explicitly allow access.

Add to `/etc/config/opennds` under `config opennds`:
```uci
config opennds
    list users_to_router 'allow tcp port 2122' # merchant port
    list preauthenticated_users 'allow tcp port 2122 to 192.168.21.21'
```

Run `service opennds restart` to make changes go into effect.

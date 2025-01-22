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
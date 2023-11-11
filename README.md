# namecheap-ddns-go
A namecheap dynamic DNS service written in Go.

## Why
There are many namecheap-based ddns updater services that can be easily found on github and other repo-services. I made this because I have multiple domains and I wanted something that would support an arbitrary number of different domains and hosts each with different ddns tokens.

## Thanks
I lifted logic/ideas/inspiration/code from the following repositories. Both of these tools work perfectly. I wanted to be able to update multiple domains at once.

- https://github.com/nickjer/namecheap-ddns (rust)
- https://github.com/Henelik/ddns-go (golang)


## Config
A yaml file (defaults to: `/etc/namecheap-ddns-go.yaml`) that looks like:

```yaml
configs: # an array, allows updating multiple domains' subdomains' ip addresses
  - domain: "example0.com" # the domain to update
    subdomains: ["*", "@", "www"] # array of hosts to update
    token: "abc123" # namecheap ddns token. See: Namecheap -> Domain list -> <choose a domain> -> Advanced DNS -> Dynamic DNS -> Dynamic DNS Password
    # ip: "1.1.1.1" # OPTIONAL ip address to set the subdomain(s).domain. Defaults to the ip making this request
  - domain: "example0.net"
    subdomains: ["host1", "host2"]
    token: "xyz321"
     # ip: "1.1.1.1" # OPTIONAL ip address to set the subdomain(s).domain. Defaults to the ip making this request
```

## Linux - systemd

If you want to set this up as a service you will need to create a service file
and corresponding timer.

An `example_namecheap-ddns-go.yaml` file has been included. Make a copy of `example_namecheap-ddns-go.yaml` eg: `namecheap-ddns-go.yaml` and change it to suit your domain setup. Then, copy the file to `/etc/namecheap-ddns-go.yaml`.

Be sure to `chmod 600 /etc/namecheap-ddns-go.yaml` as it contains your namecheap ddns token(s).

**NOTE: This repo's `.gitignore` is set to ignore `namecheap-ddns-go.yaml` to prevent accidentally storing your namecheap ddns token in source control.**

1. Create the service itself that updates your subdomains:

   ```desktop
   # /etc/systemd/system/namecheap-ddns-go.service

   [Unit]
   Description=Update DDNS records for Namecheap
   After=network-online.target

   [Service]
   Type=simple
   Environment=NAMECHEAP_DDNS_CONFIG=/etc/namecheap-ddns-go.yaml
   ExecStart=/path/to/namecheap-ddns-go
   User=<USER>

   [Install]
   WantedBy=default.target
   ```

   Be sure to fill in the correct path to your binary as well as the
   environment variables.

2. Note that the super secret token is in the `/etc/namecheap-ddns-go.yaml` file, so we should set
   restrictive permissions:

   ```shell
   sudo chmod 600 /etc/namecheap-ddns-go.yaml
   ```

3. Create the timer that runs this service:

   ```desktop
   # /etc/systemd/system/namecheap-ddns-go.timer

   [Unit]
   Description=Run DDNS update every 15 minutes
   Requires=namecheap-ddns-go.service

   [Timer]
   Unit=namecheap-ddns-go.service
   OnUnitInactiveSec=15m
   AccuracySec=1s

   [Install]
   WantedBy=timers.target
   ```

4. Now we reload the daemon with the new services and start them:

   ```shell
   sudo systemctl daemon-reload
   sudo systemctl start namecheap-ddns-go.service namecheap-ddns-go.timer
   ```

You can view the logs from the service with the following command:

```shell
sudo journalctl -u namecheap-ddns-go.service
```

There is an example [namecheap-ddns-go.service](namecheap-ddns-go.service) and [namecheap-ddns-go.timer](namecheap-ddns-go.timer) file in this repo. Note you must change the `USER` in `namecheap-ddns-go.service`.
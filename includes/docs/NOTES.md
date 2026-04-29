# Inquisitor — Notes & Caveats

## Demo flow

```sh
make up
make info        # grab MACs of ftp-server (10.10.10.10) and ftp-victim (10.10.10.20)
make attack      # → inside attacker container:
inquisitor 10.10.10.20 <victim-mac> 10.10.10.10 <server-mac>
# in another terminal:
make victim      # → inside victim:
ftp 10.10.10.10  # log in user/pass, then `get somefile` / `put somefile`
```

You should see `[FTP] RETR -> somefile` in the attacker terminal, and on
`Ctrl+C` the ARP-restore lines, then exit.

## Caveats

- Docker bridges enable a kernel feature (`hairpin` / MAC-learning) that can
  interfere with ARP spoofing in some kernels. If forwarding doesn't work,
  run the attacker with `--privileged` or test on a libvirt/VirtualBox
  network — the spoofing itself works regardless, but full MITM forwarding
  is what's bridge-sensitive.
- For real MITM (forwarding poisoned traffic onward) you'd also want
  `sysctl -w net.ipv4.ip_forward=1` inside the attacker. The spec only
  requires *poisoning* + filename sniffing, both of which work without
  forwarding.

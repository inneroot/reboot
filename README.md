# reboot over http

## Security config

youruser ALL=(root) NOPASSWD: /sbin/reboot, /sbin/reboot -f

## Without authentication

```bash
curl -X POST http://your-vps-ip:8080/api/reboot \
  -H "Content-Type: application/json" \
  -d '{"delay": 10}'
```

## With authentication

```bash
curl -X POST http://your-vps-ip:8080/api/reboot \
  -u "admin:your-auth-token" \
  -H "Content-Type: application/json" \
  -d '{"delay": 10, "force": false}'
```

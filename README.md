# reboot over http

## Security config

youruser ALL=(root) NOPASSWD: /sbin/reboot, /sbin/reboot -f

## Health Check

```bash
curl -X GET http://your-vps-ip:8080/health
```

## Without authentication

```bash
curl -X POST http://your-vps-ip:8080/reboot \
  -H "Content-Type: application/json" \
  -d '{"delay": 10}'
```

## With authentication

```bash
curl -X POST http://your-vps-ip:8080/reboot \
  -H "Content-Type: application/json" \
  -d '{"delay": 10, "force": false, "token": "your-auth-token"}'
```

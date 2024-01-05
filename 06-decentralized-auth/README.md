# Episode 6: Decentralized Auth

```bash
nsc add operator --generate-signing-key --sys --name my_org

nsc edit operator --require-signing-keys --account-jwt-server-url "nats://0.0.0.0:4222"

nsc add account TEAM_A
nsc edit account TEAM_A --sk generate
nsc add user --account TEAM_A user_a

nsc add account TEAM_B
nsc edit account TEAM_B --sk generate
nsc add user --account TEAM_B math

nsc env

nsc list keys -A

nsc generate config --nats-resolver --sys-account SYS > resolver.conf

nats-server -c resolver conf

nsc push -A

nats context save my_org_sys --nsc "nsc://my_org/SYS/sys"
nats context save my_org_user_a --nsc "nsc://my_org/TEAM_A/user_a"
nats context save my_org_math --nsc "nsc://my_org/TEAM_B/math"

nats generate creds -n math > math.creds

go run math_service.go

# Try running with math
# Try running with user a, show it isn't imported

nsc add export --account TEAM_B --name maths --subject math.* --service

nsc add import -i

# test with team a, now it can use the math service
```

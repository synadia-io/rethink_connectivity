# Rethink Connectivity: Episode 8

```sh
# After you've loaded your account from NGS into nsc, you can look at your account details
nsc describe account

# To configure dynamic permissions, we will need a signing key. Signing keys in NATS do 2 things.
# First, they are good security measures. If your root key is compromised, your are out of luck, this is thr primary reason why you should create signing keys.
# Secondly, you can create a scoped signing key, which you can think of as a user role. Any time you sign a user with a scoped signing key, permissions for that user will be dynamically managed via that signing key.
# Let's make our first signing key
nsc edit account --sk generate

# And another one
nsc edit account --sk generate

# Let's verify we have signing keys on our account
nsc describe account

# Now that we have our two signing keys, let's edit one of them to add a scope of user permissions
nsc edit signing-key -h
nsc edit signing-key --sk "YOUR_SIGNING_KEY" --role chat_user

# Now when we describe our account we can see the scoped signing key
nsc describe account

# Let's add some permissions
# When adding permissions to a scoped signing key, we can also use template functions like `name()`, `account-name()` or even pull a `tag()` off the user or account. This can be super useful for further customizing the process.
# Here, we are creating a user that is allowed to publish a message in chat and listen for others messages for whichever org is tagged in the user.
nsc edit signing-key --sk "YOUR_SIGNING_KEY" \
    --allow-pub "chat.post.{{tag(org)}}.{{name()}}" \
    --allow-sub "chat.post.{{tag(org)}}.*" \
    --allow-pub-response
    
# Now let's create some users, using our chat_user scoped signing key
nsc add user jeremy -K chat_user --tag org:synadia
nsc add user liz -K chat_user --tag org:synadia

# And inspect them
nsc describe user jeremy
nsc describe user liz

# Now let's test out a basic chat construct with the nats CLI, we'll add contexts for our two users
nats context save jeremy --select --nsc nsc://synadia/codegangsta_chat/jeremy
nats context save liz --select --nsc nsc://synadia/codegangsta_chat/liz

# And we'll attempt to subscribe as jeremy
nats sub "chat.post.synadia.*" --context jeremy

# Let's post a message as jeremy
nats pub "chat.post.synadia.jeremy" "Hello?" --context jeremy

# If I try to post a message with the liz user, on the jeremy subject
nats pub "chat.post.synadia.jeremy" "Hello?" --context liz

# but I can post as liz to the liz subject
nats pub "chat.post.synadia.liz" "Hey there!" --context liz

# So we can see that this chat system is working and secure, but let's imagine for a sec that we want to roll out a new feature, if we needed to add new permissions, we would have to generate new user credentials for jeremy and liz, and redistribute it to them. Not a great solution.
# Since we used scope signing keys though, we can actually update the signing key and push it to the server, and jeremy and liz's users will automatically be updated. Pretty cool, huh?
# Let's add a direct message feature
nsc edit signing-key --sk "YOUR_SIGNING_KEY" \
    --allow-pub "chat.dm.{{tag(org)}}.{{name()}}.*" \
    --allow-sub "chat.dm.{{tag(org)}}.*.{{name()}}"

# If we inspect jeremy, we see his permissions have automatically been updated without having to change his credentials
nsc describe user jeremy

# Now let's try our new dm feature, we'll subscribe to the subject that's from liz to jeremy
nats sub "chat.dm.synadia.liz.jeremy"

# And liz can publish to this subject
nats pub "chat.post.synadia.liz.jeremy" "Hey" --context liz

# With just a couple commands, we created a set of dynamic permissions that don't require any updates to the services or people using those credentials. Scoped signing keys are a little known feature, but can be incredibly powerful for streamlining the process of minting new user jwts

# In the next episode, we will iterate on this concept, and I will show you how to build a service that can mint user jwts on the fly that can secure clients in environments like mobile and web browsers. But that's all for today

```

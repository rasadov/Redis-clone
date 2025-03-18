# Redis Clone

## Functionality
Notes:
* In the example below I am using `redis-cli` and regular strings to simply explain logic. In the app itself the RESP protocol is being used for communication
* Redis command names are case-insensitive, so ECHO, echo and EcHo are all valid commands.

1. `PING` - for checking health. Supposed to return "PONG" as response
```bash
  $ redis-cli PING
  PONG
```
2. `ECHO` - returns same message that was sent
```bash
  $ redis-cli echo message
  message
```
* Note: if echo receives more than one word it will return an error
```bash
  $ redis-cli echo some message 
  Wrong number of arguments for 'ECHO' command
```
3. `SET` - command used to set a key to a value. Can accept px argument which is a keyword for expiry. The value is passed of expire is passed in milliseconds
```bash
  $ redis-cli set foo bar 
```
Or set expiration to be 100 milliseconds
```bash
  $ redis-cli set foo bar px 100 
```
4. `GET` - command used for getting a value with key. Returns null bulk string in case if it's not set or expired 
```bash
  $ redis-cli set foo bar
  $ redis-cli get foo
  bar
```
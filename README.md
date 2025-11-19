## Loom
In-memory pub/sub that:
- Allows clients to subscribe and publish messages over WebSockets.
- Supports multiple topics.
- Handles concurrent publishers/subscribers safely.

## Event format

Published event must follow this exact byte format:

~~~
+--------+----------+--------+------------+-------------+
| ACTION | T_LEN    | TOPIC  | P_LEN      |  PAYLOAD    |
| (1B)   | (1–2B)   | bytes  | (4B)       |  bytes      |
+--------+----------+--------+------------+-------------+
~~~

Received event by a subscriber looks like this:

~~~
+-------------+
|  PAYLOAD    |
|   bytes     |
+-------------+
~~~

See `cmd/encode/encode.go` example code.

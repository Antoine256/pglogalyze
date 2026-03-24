import random
from datetime import datetime, timedelta

OUTPUT_FILE = "pg_logs_test.log"
NUM_LINES = 1000000  # nombre de lignes à générer

LEVELS = ["DEBUG", "LOG", "INFO", "NOTICE", "WARNING", "ERROR", "FATAL", "PANIC"]
USERS = ["alice", "bob", "carol", "dave", "eve", "frank"]
DBS = ["appdb", "users", "orders", "inventory", "analytics"]

MESSAGES = [
    "duration: 12.345 ms  statement: SELECT * FROM users WHERE id = 42;",
    "duration: 1.234 ms  statement: INSERT INTO orders VALUES (1, 'item');",
    "duration: 45.678 ms  statement: UPDATE products SET price = price * 1.1;",
    "could not open file \"/var/lib/postgresql/data/base/12345/67890\": No such file or directory",
    "relation \"orders\" does not exist",
    "syntax error at or near \"FROM\"",
    "permission denied for relation users",
    "duplicate key value violates unique constraint \"users_pkey\"",
    "connection received: host=127.0.0.1 port=5432",
    "connection authorized: user=postgres database=postgres",
    "terminating connection due to idle-in-transaction timeout",
    "checkpoint starting: time, redo location 0/12345678",
    "autovacuum launcher started",
    "invalid input syntax for integer: \"abc\"",
    "could not serialize access due to concurrent update",
    "could not connect to server: Connection refused",
    "unexpected EOF on client connection",
    "starting background worker process",
    "replication terminated by primary server",
    "stats collector process started"
]

# Étape 1 : générer toutes les lignes avec timestamp
lines = []
for _ in range(NUM_LINES):
    ts = datetime.now() - timedelta(
        days=random.randint(0, 365),
        hours=random.randint(0, 23),
        minutes=random.randint(0, 59),
        seconds=random.randint(0, 59),
        milliseconds=random.randint(0, 999)
    )
    level = random.choice(LEVELS)
    pid = random.randint(1000, 9999)
    message = random.choice(MESSAGES)
    
    if random.random() < 0.5:
        user = random.choice(USERS)
        db = random.choice(DBS)
        line = (ts, f"{ts.strftime('%Y-%m-%d %H:%M:%S.%f')[:-3]} [{pid}] {level}: [{user}@{db}] {message}\n")
    else:
        line = (ts, f"{ts.strftime('%Y-%m-%d %H:%M:%S.%f')[:-3]} [{pid}] {level}: {message}\n")
    
    lines.append(line)

# Étape 2 : trier par timestamp
lines.sort(key=lambda x: x[0])

# Étape 3 : écrire dans le fichier
with open(OUTPUT_FILE, "w", buffering=1) as f:
    for _, line in lines:
        f.write(line)
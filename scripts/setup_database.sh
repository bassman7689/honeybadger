psql -U root -h database -p 26257 <<EOF
CREATE DATABASE IF NOT EXISTS honeybadger;
EOF

/migrate -path=/migrations/ -database 'cockroachdb://root@database:26257/honeybadger?sslmode=disable' up

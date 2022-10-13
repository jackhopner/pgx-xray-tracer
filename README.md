# pgx-xray-tracer

Adds [AWS X-Ray](https://github.com/aws/aws-xray-sdk-go) tracing support to [pgx](https://github.com/jackc/pgx)

## Usage

```
connConfig, err := pgxpool.ParseConfig("postgres://test:password@localhost:5432/test")
if err != nil {
    return errors.Wrap(err, "failed to connect to db")
}
connConfig.ConnConfig.Tracer = pgxxray.NewPGXTracer()

conn, err := pgxpool.NewWithConfig(context.Background(), connConfig)
if err != nil {
    return errors.Wrap(err, "failed to connect to db")
}

```

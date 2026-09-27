# LogisticsWBHackaton

Go backend prototype for automatic warehouse transport dispatch.

## What is implemented

- `api-gateway` (`:8080` in container, `:18080` on host by default)
- `data-aggregator` (`:8081` in container, `:18081` on host by default)
- `dispatch-logic` (`:8082` in container, `:18082` on host by default)
- `transport-management` (`:8083` in container, `:18083` on host by default)

All services follow clean architecture layers:
- `controller`
- `usecase`
- `domain`
- `adapter`
- `dto`

Real adapters are connected:
- `PostgreSQL` for raw records and transport requests
- `Redis` for active requests cache and latest dispatch plan storage
- `RabbitMQ` for transport events publishing

`dispatch-logic` can run with:
- real Python ML backend (`ML_BACKEND_URL`)
- built-in mock mode (`DISPATCH_USE_MOCK_ML=true`)

## Repository structure

```text
.
├── backend/
│   ├── api/
│   ├── aggregator/
│   ├── dispatch_logic/
│   ├── transport_management/
│   └── platform/
├── deploy/
│   ├── env/
│   └── postgres/init/
├── training/
│   └── solution.ipynb
├── docker-compose.yml
├── Makefile
└── .env.example
```

## One-command run

```bash
make up
```

This command builds and starts:
- Go services
- `postgres`
- `redis`
- `rabbitmq`

Useful commands:

```bash
make logs
make ps
make down
make clean
```

Default host ports are already shifted to avoid conflicts. You can still override them in `.env`:

```bash
API_HOST_PORT=18080
AGGREGATOR_HOST_PORT=18081
DISPATCH_HOST_PORT=18082
TRANSPORT_HOST_PORT=18083
POSTGRES_HOST_PORT=55432
REDIS_HOST_PORT=56379
RABBITMQ_HOST_PORT=55672
RABBITMQ_MANAGEMENT_HOST_PORT=55673
```

## Config files

Docker compose uses ready env files:
- `deploy/env/api.env`
- `deploy/env/aggregator.env`
- `deploy/env/dispatch.env`
- `deploy/env/transport.env`

Template for local custom run:
- `.env.example`

## Database init

PostgreSQL schema is initialized automatically from:
- `deploy/postgres/init/001_init.sql`

Created tables:
- `raw_records`
- `transport_requests`

`raw_records` now includes:
- `route_id`
- `office_from_id`
- `timestamp`
- `status_1..status_8`
- `target_2h`

If you already initialized an older DB volume (without `target_2h`), apply manual patch:
- `deploy/postgres/manual/001_add_target_2h.sql`
- `make db-patch-target2h`

## API endpoints

`api-gateway`:
- `POST /v1/warehouse/records`
- `POST /v1/dispatch/run`
- `GET /v1/dashboard/overview`

`data-aggregator`:
- `POST /v1/records`
- `GET /v1/records/window?limit=14000`
- `GET /v1/stats`

`dispatch-logic`:
- `POST /v1/dispatch/run`
- `GET /v1/dispatch/plans/latest`

`transport-management`:
- `POST /v1/requests/bulk`
- `GET /v1/requests/active`
- `PATCH /v1/requests/{id}/status`

## Local Go tests

```bash
make test
```

## ML model training

The training and prediction pipeline is available in
`training/solution.ipynb`. It trains a `CatBoostRegressor` to forecast
`target_2h` using route, calendar, cyclic time, lag and rolling-window
features.

Before running the notebook, install the Python dependencies:

```bash
python -m pip install numpy pandas seaborn matplotlib catboost scikit-learn pyarrow jupyter
```

Place the competition datasets in the following locations relative to the
`training` directory:

```text
training/
├── data/
│   ├── train_team_track.parquet
│   └── test_team_track.parquet
└── solution.ipynb
```

Start Jupyter from `training` so that the relative paths in the notebook are
resolved correctly:

```bash
cd training
jupyter notebook solution.ipynb
```

Run the cells in order. The training cell writes `catboost_model.cbm`; the
final cell loads this model, predicts the test set sequentially (using prior
predictions as history), and writes `submission.csv`. Both files are created
in `training/`.

## Training dataset importer

Training parquet file is expected at:
- `datasets/train_team_track.parquet`

Run dry validation (no DB writes):

```bash
make import-train-dry
```

Import into `raw_records`:

```bash
make import-train
```

Importer supports extra flags:
- `--truncate` to clean `raw_records` before import
- `--batch-size` to tune import throughput
- `--limit` to import only first N rows
- `--database-url` to override PostgreSQL DSN

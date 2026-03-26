package repo

import (
    "context"
    "database/sql"
    "time"

    "starflow-scheduler/internal/models"
)

type ExecutorRepository struct {
    db *sql.DB
}

func NewExecutorRepository(db *sql.DB) *ExecutorRepository {
    return &ExecutorRepository{db: db}
}

func (r *ExecutorRepository) Upsert(ctx context.Context, exec *models.Executor) (int64, error) {
    res, err := r.db.ExecContext(ctx, `
        INSERT INTO executors (
            executor_code, host, ip, grpc_addr, http_addr, tags, capacity,
            current_load, version, status, last_heartbeat_at, metadata
        ) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)
        ON DUPLICATE KEY UPDATE
            host = VALUES(host),
            ip = VALUES(ip),
            grpc_addr = VALUES(grpc_addr),
            http_addr = VALUES(http_addr),
            tags = VALUES(tags),
            capacity = VALUES(capacity),
            current_load = VALUES(current_load),
            version = VALUES(version),
            status = VALUES(status),
            last_heartbeat_at = VALUES(last_heartbeat_at),
            metadata = VALUES(metadata)`,
        exec.ExecutorCode, exec.Host, exec.IP, exec.GrpcAddr, exec.HttpAddr, exec.Tags,
        exec.Capacity, exec.CurrentLoad, exec.Version, exec.Status, exec.LastHeartbeat,
        exec.Metadata,
    )
    if err != nil {
        return 0, err
    }
    id, err := res.LastInsertId()
    if err != nil {
        return 0, nil
    }
    return id, nil
}

func (r *ExecutorRepository) UpdateHeartbeat(ctx context.Context, executorCode string, currentLoad int) error {
    _, err := r.db.ExecContext(ctx, `
        UPDATE executors
        SET current_load = ?, last_heartbeat_at = ?, status = 'online'
        WHERE executor_code = ?`, currentLoad, time.Now(), executorCode)
    return err
}

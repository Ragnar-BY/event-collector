
CREATE TABLE IF NOT EXISTS events (
ClientTime Timestamp,
DeviceID UUID,
DeviceOS String,
Session String,
Sequence UInt64,
Event String,
ParamInt UInt64,
ParamStr String,
ClientIP IPv4,
ServerTime Timestamp
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(ServerTime)
ORDER BY (ServerTime)

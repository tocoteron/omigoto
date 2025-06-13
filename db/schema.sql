CREATE EXTENSION vector;

CREATE TABLE youtube_channels (
    channel_id TEXT PRIMARY KEY, -- UC1cnByKe24JjTv38tH_7BYw
    handle TEXT NOT NULL UNIQUE -- @izuho_omi
);

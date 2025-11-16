INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status) VALUES ($1, $2, $3, $4)
RETURNING pull_request_id, status, author_id, pull_request_name, merged_at, created_at;

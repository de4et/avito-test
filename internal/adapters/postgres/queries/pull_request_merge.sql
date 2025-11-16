UPDATE pull_requests
SET status = $1,
merged_at = COALESCE(merged_at, NOW())
WHERE pull_request_id = $2
RETURNING pull_request_id, status, author_id, pull_request_name, merged_at, created_at;

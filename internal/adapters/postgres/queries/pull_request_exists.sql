SELECT 1
FROM pull_requests
WHERE pull_request_id = $1
LIMIT 1

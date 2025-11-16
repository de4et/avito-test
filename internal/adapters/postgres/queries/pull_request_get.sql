select
  pull_request_id,
  status,
  author_id,
  pull_request_name,
  merged_at,
  created_at
from
  pull_requests
where pull_request_id = $1;

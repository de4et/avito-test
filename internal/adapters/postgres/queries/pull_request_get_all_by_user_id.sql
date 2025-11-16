SELECT
  pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status
FROM
  pull_requests pr
  JOIN pull_request_reviewers prr ON pr.pull_request_id = prr.pull_request_id
WHERE
  prr.reviewer_id = $1;

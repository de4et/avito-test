UPDATE pull_request_reviewers
SET reviewer_id = $1
WHERE reviewer_id = $2 AND pull_request_id=$3;

#!/bin/bash

set -euo pipefail

BASE_URL="http://localhost:8080"

curl -s "$BASE_URL/api/healthz"
echo ""
curl -s "$BASE_URL/app/"
echo ""
curl -s "$BASE_URL/admin/metrics"
echo ""
curl -s -X POST "$BASE_URL/admin/reset"
echo ""
curl -s "$BASE_URL/admin/metrics"
echo ""
curl -s -X POST --data '""' "$BASE_URL/api/validate_chirp"
echo ""
curl -s -X POST --data '{"Body":"this is a chirp"}' "$BASE_URL/api/validate_chirp"
echo ""
curl -s -X POST --data '{"Body":"this is a loooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooong chirp"}' "$BASE_URL/api/validate_chirp"

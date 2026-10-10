#!/usr/bin/env bash
# 轻量冒烟：示例环境不含真实密钥，公开金句接口可读。
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

fail() {
  echo "smoke failed: $*" >&2
  exit 1
}

frontend_env="book-quote-chat-frontend/.env.example"
backend_env="book-quote-chat-backend/.env.example"
dev_env="book-quote-chat-frontend/.env.dev"

for f in "$frontend_env" "$backend_env" README.md; do
  [[ -f "$f" ]] || fail "missing $f"
done

grep -q 'book-quote-chat-frontend' README.md || fail "README missing frontend dir"
grep -q 'book-quote-chat-backend' README.md || fail "README missing backend dir"
grep -q 'VITE_API_BASE_URL' README.md || fail "README missing frontend env keys"
grep -q 'ALLOWED_ORIGIN=' README.md || fail "README missing backend ALLOWED_ORIGIN"
grep -q 'UPLOAD_ACCEPT=' README.md || fail "README missing UPLOAD_ACCEPT"
grep -q 'JWT_SECRET=' README.md || fail "README missing JWT_SECRET"

if grep -q 'ALLOWED_ORIGINS\|UPLOAD_EXTS\|UPLOAD_MAX_SIZE_MB' README.md; then
  fail "README still documents removed env names"
fi

secret_re='dz1991|43\.100\.91\.149|007eJx|c59f9a9fc7a04fb8a8a9291f7f6fc21c|192\.168\.|ALLOWED_ORIGIN=\*'
for f in "$frontend_env" "$backend_env" "$dev_env" README.md; do
  [[ -f "$f" ]] || continue
  if grep -E -q "$secret_re" "$f"; then
    fail "$f contains a real address or credential"
  fi
done

grep -q 'VITE_API_BASE_URL=' "$frontend_env" || fail "frontend example missing VITE_API_BASE_URL"
grep -q 'VITE_WS_URL=' "$frontend_env" || fail "frontend example missing VITE_WS_URL"
grep -q 'JWT_SECRET=' "$backend_env" || fail "backend example missing JWT_SECRET"
grep -q 'ALLOWED_ORIGIN=' "$backend_env" || fail "backend example missing ALLOWED_ORIGIN"
grep -q 'UPLOAD_ACCEPT=' "$backend_env" || fail "backend example missing UPLOAD_ACCEPT"

cd "$ROOT/book-quote-chat-backend"
# 只编译冒烟用例，避开仓库里原有、与本次无关的测试编译错误。
(
  cd api
  srcs=()
  for f in *.go; do
    [[ "$f" == *_test.go && "$f" != smoke_test.go ]] && continue
    srcs+=("$f")
  done
  go test -count=1 -run TestSmokePublicQuoteEndpoints "${srcs[@]}"
)
(
  cd service
  srcs=()
  for f in *.go; do
    [[ "$f" == *_test.go && "$f" != jwt_secret_test.go ]] && continue
    srcs+=("$f")
  done
  go test -count=1 -run TestJWTSecretFromEnv "${srcs[@]}"
)

echo "smoke ok"

# 获取 zcache 最新的 commits
zcache_commit_hash=$(curl -s https://api.github.com/repos/Zondax/golem/commits | grep "sha" | head -n 1 | cut -d '"' -f 4)
# 使用提取的 commit hash 通过 go get 获取 zcache-go
go get github.com/zondax/golem@$zcache_commit_hash

# 获取 dns 最新的 commits
dns_commit_hash=$(curl -s https://api.github.com/repos/miekg/dns/commits | grep "sha" | head -n 1 | cut -d '"' -f 4)
# 使用提取的 commit hash 通过 go get 获取 dns
go get github.com/miekg/dns@$dns_commit_hash

# 获取 gopsutil 最新的 commits
gopsutil_commit_hash=$(curl -s https://api.github.com/repos/shirou/gopsutil/commits | grep "sha" | head -n 1 | cut -d '"' -f 4)
# 使用提取的 commit hash 通过 go get 获取 gopsutil
go get github.com/shirou/gopsutil@$gopsutil_commit_hash

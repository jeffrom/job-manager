package beredis

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"

	"github.com/jeffrom/job-manager/pkg/resource"
)

const idxKey = "mjob:idx"

func indexKey(parts ...string) string {
	return redisKey(idxKey, strings.Join(parts, "-"))
}

// to index a job lexically, basically can have a sorted set for each axis to query
// STATUS:ID
// CLAIMSET:ID - this would be multiple members with a sorted combo of claims
// could do composite index too
// STATUS:CLAIMSET:ID
// but maybe for now, just filter the results after by only including results
// that matched all queries

var lexicalSeparator = []byte{0x0, 0x0}

func lexicalKey(parts ...string) string {
	var b strings.Builder
	for i, part := range parts {
		if i > 0 {
			b.Write(lexicalSeparator)
		}
		b.WriteString(part)
	}

	return b.String()
}

func indexJobMembers(jb *resource.Job, queueName string) map[string][]string {
	status := strconv.FormatInt(int64(jb.Status), 10)
	// fmt.Println("status #", status, jb.Status.String())
	members := map[string][]string{
		// job queue / status index
		indexKey("status"):            []string{lexicalKey(status, jb.ID)},
		indexKey("status", queueName): []string{lexicalKey(status, jb.ID)},
	}

	if jb.Data != nil && len(jb.Data.Claims) > 0 {
		claims := jb.Data.Claims.Format()
		claimMembers := make([]string, len(claims))
		for i, part := range claims {
			claimMembers[i] = lexicalKey(part, jb.ID)
		}
		members[indexKey("claims")] = claimMembers
	}

	return members
}

func (be *RedisBackend) indexJob(ctx context.Context, pipe redis.Pipeliner, queueName string, jb, prevJob *resource.Job) error {
	var prevIdx map[string][]string
	if prevJob != nil {
		prevIdx = indexJobMembers(prevJob, queueName)
	}

	for key, members := range indexJobMembers(jb, queueName) {
		if prevMembers, ok := prevIdx[key]; ok {
			for _, prevMember := range prevMembers {
				pipe.ZRem(ctx, key, &redis.Z{Score: 0, Member: prevMember})
			}
		}
		for _, member := range members {
			// fmt.Printf("indexJob: %q, member: %q\n", key, member)
			pipe.ZAdd(ctx, key, &redis.Z{Score: 0, Member: member})
		}
	}
	return nil
}

func (be *RedisBackend) indexLookup(ctx context.Context, limit int64, opts *resource.JobListParams) ([]string, error) {
	var queueStatus []*redis.StringSliceCmd
	_, err := be.rds.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		if len(opts.Names) > 0 {
			if len(opts.Statuses) > 0 {
				for _, queueName := range opts.Names {
					key := indexKey("status", queueName)
					for _, st := range opts.Statuses {
						lex := lexicalKey(strconv.FormatInt(int64(st), 10), "")
						minlex := `[` + lex
						maxlex := `(` + lex + string(0xff)
						// fmt.Printf("lex: %q\n", lex)
						queueStatus = append(queueStatus, pipe.ZRangeByLex(ctx, key, &redis.ZRangeBy{
							Min:   minlex,
							Max:   maxlex,
							Count: limit,
						}))
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var allIds []string
	for _, qst := range queueStatus {
		// fmt.Printf("qst args: %q, res: %+v, err: %+v\n", qst.Args(), qst.Val(), qst.Err())
		allIds = append(allIds, qst.Val()...)
	}

	seen := make(map[string]bool)
	var resIds []string
	for _, val := range allIds {
		parts := strings.Split(val, string(lexicalSeparator))
		id := parts[len(parts)-1]
		if ok := seen[id]; ok {
			continue
		}
		resIds = append(resIds, id)
		seen[id] = true

		if int64(len(resIds)) >= limit {
			break
		}
	}
	return resIds, nil
}

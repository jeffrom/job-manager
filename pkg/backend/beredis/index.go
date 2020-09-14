package beredis

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"

	"github.com/jeffrom/job-manager/pkg/label"
	"github.com/jeffrom/job-manager/pkg/resource"
)

const idxKey = "mjob:idx"
const checkpointKey = "mjob:chk"

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

func buildLexicalIndex(jb *resource.Job, queueName string) map[string][]string {
	status := strconv.FormatInt(int64(jb.Status), 10)
	// fmt.Println("status #", status, jb.Status.String())
	members := map[string][]string{
		// job queue / status index
		indexKey("status"):            []string{lexicalKey(status, jb.ID)},
		indexKey("status", queueName): []string{lexicalKey(status, jb.ID)},
	}

	// if jb.Data != nil && len(jb.Data.Claims) > 0 {
	// 	claims := jb.Data.Claims.Format()
	// 	claimMembers := make([]string, len(claims))
	// 	for i, part := range claims {
	// 		claimMembers[i] = lexicalKey(part, jb.ID)
	// 	}
	// 	members[indexKey("claims")] = claimMembers
	// }

	return members
}

func buildSetIndex(jb *resource.Job, labels label.Labels) map[string]string {
	m := make(map[string]string)

	if jb.Data != nil && len(jb.Data.Claims) > 0 {
		claims := jb.Data.Claims.Format()
		for _, claim := range claims {
			// fmt.Println("indexJobMembers claim:", claim)
			m[indexKey("claim", claim)] = jb.ID
		}
	}

	if len(labels) > 0 {
		for k, v := range labels {
			m[indexKey("label", k)] = jb.ID
			m[indexKey("label", k+"="+v)] = jb.ID
		}
	}

	return m
}

func (be *RedisBackend) indexJob(ctx context.Context, pipe redis.Pipeliner, queueName string, labels label.Labels, jb, prevJob *resource.Job) error {
	var prevLexIdx map[string][]string
	if prevJob != nil {
		prevLexIdx = buildLexicalIndex(prevJob, queueName)
	}

	for key, members := range buildLexicalIndex(jb, queueName) {
		if prevMembers, ok := prevLexIdx[key]; ok {
			for _, prevMember := range prevMembers {
				pipe.ZRem(ctx, key, prevMember)
			}
		}
		for _, member := range members {
			// fmt.Printf("indexJob: %q, member: %q\n", key, member)
			pipe.ZAdd(ctx, key, &redis.Z{Score: 0, Member: member})
		}
	}

	for key, id := range buildSetIndex(jb, labels) {
		pipe.SAdd(ctx, key, id)
	}
	return nil
}

func (be *RedisBackend) checkpointJob(ctx context.Context, pipe redis.Pipeliner, jobID, logID string) error {
	pipe.HSet(ctx, checkpointKey, jobID, logID)
	return nil
}

func (be *RedisBackend) indexLookup(ctx context.Context, limit int64, opts *resource.JobListParams) ([]string, error) {
	// fmt.Printf("AAAA IDX LOOKUP OPTS: %+v\n", opts)
	var claims []*redis.StringSliceCmd
	var queueStatus []*redis.StringSliceCmd
	var status []*redis.StringSliceCmd
	_, err := be.rds.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		claims = be.indexLookupClaims(ctx, pipe, limit, opts)
		queueStatus = be.indexLookupQueueStatus(ctx, pipe, limit, opts)
		status = be.indexLookupStatus(ctx, pipe, limit, opts)
		return nil
	})
	if err != nil {
		return nil, err
	}

	var allIds []string
	for _, qst := range claims {
		allIds = append(allIds, qst.Val()...)
	}
	for _, qst := range queueStatus {
		// fmt.Printf("qst args: %q, res: %+v, err: %+v\n", qst.Args(), qst.Val(), qst.Err())
		allIds = append(allIds, qst.Val()...)
	}
	for _, qst := range status {
		allIds = append(allIds, qst.Val()...)
	}
	if len(allIds) == 0 {
		return nil, nil
	}

	seen := make(map[string]bool)
	var filteredIds []string
	for _, val := range allIds {
		parts := strings.Split(val, string(lexicalSeparator))
		id := parts[len(parts)-1]
		if ok := seen[id]; ok {
			continue
		}
		filteredIds = append(filteredIds, id)
		seen[id] = true

		if int64(len(filteredIds)) >= limit {
			break
		}
	}

	resCmd, err := be.rds.HMGet(ctx, checkpointKey, filteredIds...).Result()
	if err != nil {
		return nil, err
	}

	resIds := make([]string, len(filteredIds))
	for i, iid := range resCmd {
		resIds[i] = iid.(string)
	}
	// fmt.Println("indexLookup resIds:", resIds)
	return resIds, nil
}

func (be *RedisBackend) indexLookupQueueStatus(ctx context.Context, pipe redis.Pipeliner, limit int64, opts *resource.JobListParams) []*redis.StringSliceCmd {
	var queueStatus []*redis.StringSliceCmd
	if len(opts.Names) == 0 || len(opts.Statuses) == 0 {
		return nil
	}
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
	return queueStatus
}

func (be *RedisBackend) indexLookupStatus(ctx context.Context, pipe redis.Pipeliner, limit int64, opts *resource.JobListParams) []*redis.StringSliceCmd {
	key := indexKey("status")
	var status []*redis.StringSliceCmd
	for _, st := range opts.Statuses {
		lex := lexicalKey(strconv.FormatInt(int64(st), 10), "")
		minlex := `[` + lex
		maxlex := `(` + lex + string(0xff)
		status = append(status, pipe.ZRangeByLex(ctx, key, &redis.ZRangeBy{
			Min:   minlex,
			Max:   maxlex,
			Count: limit,
		}))
	}
	return status
}

func (be *RedisBackend) indexLookupClaims(ctx context.Context, pipe redis.Pipeliner, limit int64, opts *resource.JobListParams) []*redis.StringSliceCmd {
	if len(opts.Claims) == 0 {
		return nil
	}

	claims := opts.Claims.Format()
	keys := make([]string, len(claims))
	for i, claim := range claims {
		keys[i] = indexKey("claim", claim)
	}

	cmd := pipe.SUnion(ctx, keys...)
	return []*redis.StringSliceCmd{cmd}
}

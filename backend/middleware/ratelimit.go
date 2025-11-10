package middlewares

import (
	"net/http"
	"social-network/app/helper"
	"time"
)

type RateLimit struct {
	count        int
	FirstTime    time.Time
	BlockedUntil time.Time
	UserId       string
}

type ErrorStruct struct {
	Type string
	Text string
}

var (
	CommentRateLimits        = make(map[string]*RateLimit)
	PostRateLimits           = make(map[string]*RateLimit)
	LikesRateLimits          = make(map[string]*RateLimit)
	FollowUnfollowRateLimits = make(map[string]*RateLimit)
	StoryRateLimits          = make(map[string]*RateLimit)
	GroupJoinRateLimits      = make(map[string]*RateLimit)
	GroupInviteRateLimits    = make(map[string]*RateLimit)
	EventsRateLimits         = make(map[string]*RateLimit)
	EventActionRateLimit     = make(map[string]*RateLimit)
	CreateGroupRateLimit     = make(map[string]*RateLimit)
)

func CheckRateLimit(ratelimit *RateLimit, window time.Duration, maxAttempts int) bool {
	if time.Now().Before(ratelimit.BlockedUntil) {
		return false
	}
	if time.Now().After(ratelimit.BlockedUntil) && ratelimit.count > maxAttempts {
		ratelimit.FirstTime = time.Now()
		ratelimit.BlockedUntil = time.Time{}
		ratelimit.count = 0
	}
	ratelimit.count++
	if ratelimit.count > maxAttempts {
		ratelimit.BlockedUntil = time.Now().Add(window)
		return false
	}
	return true
}

func UserInfos(r *http.Request) (*RateLimit, bool) {
	rateLimit := &RateLimit{
		count:        0,
		FirstTime:    time.Now(),
		BlockedUntil: time.Time{},
		UserId:       "",
	}
	_, userID := helper.IsLoggedIn(r)
	rateLimit.UserId = userID
	return rateLimit, true
}

func RatelimitMiddleware(next http.HandlerFunc, rateLimitType string, maxAttempts int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		theMap := make(map[string]*RateLimit)
		if rateLimitType == "posts" {
			theMap = PostRateLimits
		} else if rateLimitType == "comments" {
			theMap = CommentRateLimits
		} else if rateLimitType == "follow" || rateLimitType == "unfollow" {
			theMap = FollowUnfollowRateLimits
		} else if rateLimitType == "likes" {
			theMap = LikesRateLimits
		} else if rateLimitType == "story" {
			theMap = StoryRateLimits
		} else if rateLimitType == "joinGroup" {
			theMap = GroupJoinRateLimits
		} else if rateLimitType == "groupInvite" {
			theMap = GroupInviteRateLimits
		} else if rateLimitType == "events" {
			theMap = EventsRateLimits
		}else if rateLimitType == "eventAction"{
			theMap = EventActionRateLimit
		}else if rateLimitType == "createGroup"{
			theMap = CreateGroupRateLimit
		}
		userRateLimit, ok := UserInfos(r)

		if !ok {
			helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			// errorr := ErrorStruct{
			// 	Type: "error",
			// 	Text: "Unauthorized",
			// }

			// w.Header().Set("Content-Type", "application/json")
			// w.WriteHeader(http.StatusUnauthorized)
			// json.NewEncoder(w).Encode(errorr)
			return
		}

		ratelimit, exists := theMap[userRateLimit.UserId]
		if !exists {
			AddUserToTheMap(userRateLimit, theMap)
			ratelimit = userRateLimit
		}

		if !CheckRateLimit(ratelimit, 1*time.Minute, maxAttempts) {
			helper.RespondWithError(w, http.StatusTooManyRequests, "Too many requests")
			return
		}
		next.ServeHTTP(w, r)
	}
}

func AddUserToTheMap(ratelimit *RateLimit, theMap map[string]*RateLimit) {
	theMap[ratelimit.UserId] = ratelimit
}

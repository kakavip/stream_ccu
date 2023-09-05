package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

type LimitCCUDto struct {
	UserId  string
	Session string
}

type LimitCCUReq struct {
	UserId  string `json:"user_id"`
	Session string `json:"session"`
}

func (req LimitCCUReq) toDto() *LimitCCUDto {
	return &LimitCCUDto{
		UserId:  req.UserId,
		Session: req.Session,
	}
}

type LimitCCURes struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

const (
	PING_ACTION          = "ping"
	REFRESH_ACTION       = "refresh"
	END_ACTION           = "end"
	TOKEN_CACHE_PREFIX   = "token_limit_ccu:"
	DRM_LIMIT_CCU_PREFIX = "drm_limit_ccu:"
	DRM_LIMIT_CCU_EXP    = 2
)

func limitCCUHandler(w http.ResponseWriter, r *http.Request) {
	accountId := chi.URLParam(r, "account_id")
	tenantSlug := chi.URLParam(r, "tenant_slug")
	action := chi.URLParam(r, "action")

	var session, token string

	var requestDto *LimitCCUDto = nil
	if r.Method == "GET" {
		session = r.URL.Query().Get("session")
		if session != "" {
			requestDto = &LimitCCUDto{
				UserId:  accountId,
				Session: r.URL.Query().Get("session"),
			}
		}
		token = r.URL.Query().Get("token")
		if token != "" {
			val, _ := rdb.Get(rCtx, TOKEN_CACHE_PREFIX+token).Result()
			requestDto = &LimitCCUDto{
				UserId:  accountId,
				Session: val,
			}
		}
	} else if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var t LimitCCUReq
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}

		requestDto = t.toDto()
	}

	if requestDto != nil {
		handleLimitCCU(tenantSlug, accountId, *requestDto, action)
	}

	if r.Method == "GET" {
		switch action {
		case PING_ACTION:
			if token != "" {
				response.success(w, LimitCCURes{})
				return
			}
			if session != "" {
				newToken := strings.ToLower(randSeq(64))
				rdb.Set(rCtx, TOKEN_CACHE_PREFIX+newToken, session, time.Hour*24)
				response.success(w, LimitCCURes{
					Token: newToken,
				})
				return
			}

			response.success(w, LimitCCURes{})

		case REFRESH_ACTION:
			newToken := strings.ToLower(randSeq(64))
			rdb.Set(rCtx, TOKEN_CACHE_PREFIX+newToken, session, time.Hour*24)
			response.success(w, LimitCCURes{
				Token: newToken,
			})
		case END_ACTION:
			response.success(w, LimitCCURes{})
		}
	} else {
		response.success(w, LimitCCURes{Status: "Success"})
	}
}

func getConcurrentStreamingLimit(tenantSlug string) int {
	return 5
}

func _getUserTimeBox(userId string, next int) string {
	now := time.Now()
	sec := now.Unix()
	return DRM_LIMIT_CCU_PREFIX + fmt.Sprintf("%s_%d", userId, int(sec/60)+next)
}

func _getExpiredTime(expiredMinutes int) int {
	now := time.Now()
	sec := now.Unix()

	return int(sec/60)*60 + expiredMinutes*60 - int(sec)
}

func checkUserTimeBox(userId string, token string, maxCcu int) {
	var crtTimeBox string = _getUserTimeBox(userId, 0)
	var preTimeBox string = _getUserTimeBox(userId, -1)

	var crtTokens, preTokens []string
	crtTokensString, e := rdb.Get(rCtx, crtTimeBox).Result()
	if e != nil {
		crtTokens = []string{}
	} else {
		crtTokens = strings.Split(crtTokensString, "_")
	}
	preTokensString, e := rdb.Get(rCtx, preTimeBox).Result()
	if e != nil {
		preTokens = []string{}
	} else {
		preTokens = strings.Split(preTokensString, "_")
	}
	if slices.Contains(crtTokens, token) {
		return
	}

	if maxCcu == 0 || len(uniqueArray(append(append([]string{}, crtTokens...), preTokens...))) < maxCcu || slices.Contains(preTokens, token) {
		crtTokens = append(crtTokens, token)
		expTime := int(time.Second) * _getExpiredTime(DRM_LIMIT_CCU_EXP)
		rdb.Set(rCtx, crtTimeBox, strings.Join(crtTokens, "_"), time.Duration(expTime))
		return

	} else {
		// raise BaseApiException(
		// 	"Number of devices were limited.", error_code="max_devices_ccu_limited"
		// )
		panic("Number of devices were limited.")
	}

}

func endUserSession(userId string, token string) {
	crtTimeBox := _getUserTimeBox(userId, 0)
	preTimeBox := _getUserTimeBox(userId, -1)

	var crtTokens, preTokens []string
	crtTokenString, e := rdb.Get(rCtx, crtTimeBox).Result()
	if e != nil {
		crtTokens = []string{}
	} else {
		crtTokens = strings.Split(crtTokenString, "_")
	}

	preTokenString, e := rdb.Get(rCtx, preTimeBox).Result()
	if e != nil {
		preTokens = []string{}
	} else {
		preTokens = strings.Split(preTokenString, "_")
	}

	if slices.Contains(crtTokens, token) {
		crtTokens = uniqueArray(crtTokens)
		crtTokens = RemoveItem(crtTokens, token)
		expTime := int(time.Second) * _getExpiredTime(DRM_LIMIT_CCU_EXP)
		rdb.Set(rCtx, crtTimeBox, crtTokens, time.Duration(expTime))
	}

	if slices.Contains[[]string, string](preTokens, token) {
		preTokens = uniqueArray(preTokens)
		preTokens = RemoveItem(preTokens, token)
		expTime := int(time.Second) * _getExpiredTime(DRM_LIMIT_CCU_EXP)
		rdb.Set(rCtx, preTimeBox, preTokens, time.Duration(expTime))
	}
}

func handleLimitCCU(tenantSlug string, accountId string, dto LimitCCUDto, action string) {
	maxCcu := getConcurrentStreamingLimit(tenantSlug)

	switch action {
	case PING_ACTION:
		checkUserTimeBox(dto.UserId, dto.Session, maxCcu)
	case END_ACTION:
		endUserSession(dto.UserId, dto.Session)
	}
}

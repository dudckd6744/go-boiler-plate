package user

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type CommonResposne struct {
	Data   interface{} `json:"data"`
	Status int         `json:"status"`
	Error  interface{} `json:"error"`
}

func Response(w http.ResponseWriter, data interface{}, status int, error error, recoverError ...interface{}) {
	var res CommonResposne

	if status == http.StatusOK {
		res.Data = data
	} else {
		if error != nil {
			res.Error = error.Error()
		} else {
			res.Error = recoverError
		}
	}
	res.Status = status

	w.Header().Set("Content-Type", "application/json;")
	json.NewEncoder(w).Encode(res)
}

func UserController(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// 서버 오류날때 서버다운되지않게
	defer func() {
		if r := recover(); r != nil {
			fmt.Print(r)
			Response(w, nil, http.StatusBadRequest, nil, r)
		}
	}()

	Service.InitService(db)

	switch r.Method {
	// POST 생성
	case http.MethodPost:
		var body struct {
			Email string
			Name  string
		}

		err := json.NewDecoder(r.Body).Decode(&body)

		if err != nil {
			Response(w, nil, http.StatusInternalServerError, err)
			return
		}

		_, error := Service.CreateUser(body)

		if error != nil {
			Response(w, nil, http.StatusBadRequest, err)
			return
		}

		Response(w, "OK", http.StatusOK, nil)
	// GET 조회
	case http.MethodGet:
		// 별도 라우팅 패키지를 사용하지않으면 query Param은 직접 URL 구문 분석해야됨

		// string의 zeroValue는 ""
		stringTypeId := strings.TrimPrefix(r.URL.Path, "/api/v1/user/")

		// int의 zeroValue는 0
		conversionNumberTypeId, _ := strconv.Atoi(stringTypeId)

		// param이 없다는건 findAllUser함수를 호출하는것
		if conversionNumberTypeId == 0 {
			result, err := Service.FindAllUser()

			if err != nil {
				Response(w, nil, http.StatusBadRequest, err)
				return
			}

			Response(w, result, http.StatusOK, nil)

			// param 으로 넘오면 id로 특정 User 찾기
		} else {

			result, err := Service.FindDetailUser(conversionNumberTypeId)

			if err != nil {
				switch err.Error() {
				case "NOT FOUND":
					Response(w, nil, http.StatusBadRequest, errors.New("user not found"))
				default:
					Response(w, nil, http.StatusBadRequest, err)
				}
				return
			}
			Response(w, result, http.StatusOK, nil)
		}
		/**
		* patch API 작업하다가 포인터 변수로 함수전달하는데 메모리가 각각 함수마다 달라지길래 원인 파악
		* bear://x-callback-url/open-note?id=06FCABCE-2750-4C4F-9A8A-831AFAB6CBA4-1458-000005C23D966CD4
		* patch API가 현재 수정이 안됨. 쿼리단에 문제가 있는듯...
		**/
	case http.MethodPatch:
		// string의 zeroValue는 ""
		stringTypeId := strings.TrimPrefix(r.URL.Path, "/api/v1/user/")

		// int의 zeroValue는 0
		conversionNumberTypeId, _ := strconv.Atoi(stringTypeId)

		var body struct {
			Name string
		}

		err := json.NewDecoder(r.Body).Decode(&body)

		if err != nil {
			Response(w, nil, http.StatusInternalServerError, err)
			return
		}

		result, err := Service.PatchUserName(&conversionNumberTypeId, &body)

		if err != nil {
			Response(w, nil, http.StatusBadRequest, err)
			return
		}

		affected, _ := result.RowsAffected()

		Response(w, affected, http.StatusOK, nil)
	case http.MethodDelete:
		// string의 zeroValue는 ""
		stringTypeId := strings.TrimPrefix(r.URL.Path, "/api/v1/user/")

		// int의 zeroValue는 0
		conversionNumberTypeId, _ := strconv.Atoi(stringTypeId)

		_, err := Service.DeleteUserById(&conversionNumberTypeId)

		if err != nil {
			Response(w, nil, http.StatusBadRequest, err)
			return
		}

		Response(w, "OK", http.StatusOK, nil)

	}

}

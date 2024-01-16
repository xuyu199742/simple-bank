package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	mockDb "simple-bank/db/mock"
	db "simple-bank/db/sqlc"
	"simple-bank/util"
	"testing"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testCase := []struct {
		name          string
		accountID     int64
		buildSubs     func(store *mockDb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "ok",
			accountID: account.ID,
			buildSubs: func(store *mockDb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account.ID).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "not found",
			accountID: account.ID,
			buildSubs: func(store *mockDb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account.ID).
					Times(1).
					Return(account, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				//requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "server error",
			accountID: account.ID,
			buildSubs: func(store *mockDb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account.ID).
					Times(1).
					Return(account, pgx.ErrTooManyRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				//requireBodyMatchAccount(t, recorder.Body, account)
			},
		},

		{
			name:      "request error",
			accountID: 0,
			buildSubs: func(store *mockDb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		//todo add other case..
	}

	for i := range testCase {
		tc := testCase[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			// build mock store.
			store := mockDb.NewMockStore(ctrl)
			tc.buildSubs(store)

			// start server and send request.
			server := NewServer(store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/account/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			// check response status
			tc.checkResponse(t, recorder)
		})
	}

}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var res Response
	var goAccount db.Account
	err = json.Unmarshal(data, &res)

	jsonData, err := json.Marshal(res.Data)
	err = json.Unmarshal(jsonData, &goAccount)

	require.NoError(t, err)
	require.Equal(t, account, goAccount)
}

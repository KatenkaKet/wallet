package handler

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/KatenkaKet/wallet"
	"github.com/KatenkaKet/wallet/pkg/service"
	mock_service "github.com/KatenkaKet/wallet/pkg/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/magiconair/properties/assert"
)

func uuidFromString(s string) uuid.UUID {
	var id uuid.UUID
	if err := id.DecodeText(nil, []byte(s)); err != nil {
		log.Fatal(err)
	}
	return id
}

// Тут тесты для двух запросов
// TestHandler_getWalletBalance
// TestHandler_createWalletTransaction

func TestHandler_getWalletBalance(t *testing.T) {
	type mockBehavior func(s *mock_service.MockWallet, wallet wallet.Wallet)

	testTable := []struct {
		name         string
		inputWallet  wallet.Wallet
		mockBehavior mockBehavior
		expectedCode int
		expectedBody string
	}{
		{
			name: "success",
			inputWallet: wallet.Wallet{
				ValletId: uuidFromString("11111111-1111-1111-1111-111111111111"),
				Balance:  100.5,
			},
			mockBehavior: func(s *mock_service.MockWallet, w wallet.Wallet) {
				s.EXPECT().GetBalance(gomock.Any(), w.ValletId).Return(w.Balance, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"balance":100.5}`,
		},
		{
			name: "service error",
			inputWallet: wallet.Wallet{
				ValletId: uuidFromString("22222222-2222-2222-2222-222222222222"),
			},
			mockBehavior: func(s *mock_service.MockWallet, w wallet.Wallet) {
				s.EXPECT().GetBalance(gomock.Any(), w.ValletId).Return(float64(0), errors.New("wallet not found"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error":"wallet not found"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			// Создаём контроллер моков
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Создаём мок-сервис
			mockWallet := mock_service.NewMockWallet(ctrl)
			test.mockBehavior(mockWallet, test.inputWallet)

			// Создаём handler
			srv := &service.Service{Wallet: mockWallet}
			h := NewHandler(srv)

			// Инициализируем роутер
			r := gin.New()
			r.GET("/api/v1/wallets/:id", h.getWalletBalance)

			// Создаём GET-запрос
			url := "/api/v1/wallets/" + test.inputWallet.ValletId.UUID.String()
			//fmt.Println(url)

			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			// Выполняем запрос
			r.ServeHTTP(w, req)
			//fmt.Println(w.Code)
			//fmt.Println(string(w.Body.String()))

			// Проверяем результаты
			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
		})
	}
}

func TestHandler_createWalletTransaction(t *testing.T) {
	type mockBehavior func(s *mock_service.MockWallet, WT wallet.WalletTransactions)

	testTable := []struct {
		name         string
		inputBody    string
		inputWT      wallet.WalletTransactions
		mockBehavior mockBehavior
		expectedCode int
		expectedBody string
	}{
		{
			name:      "success",
			inputBody: `{"valletId":"11111111-1111-1111-1111-111111111111","operationType":"DEPOSIT","amount":100.5}`,
			inputWT: wallet.WalletTransactions{
				ValletId:      uuidFromString("11111111-1111-1111-1111-111111111111"),
				OperationType: "DEPOSIT",
				Amount:        100.5,
			},
			mockBehavior: func(s *mock_service.MockWallet, WT wallet.WalletTransactions) {
				s.EXPECT().UpdateBalance(gomock.Any(), WT).Return(nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"status":"success"}`,
		},
		{
			name:      "service error",
			inputBody: `{"valletId":"22222222-2222-2222-2222-222222222222","operationType":"DEPOSIT","amount":50}`,
			inputWT: wallet.WalletTransactions{
				ValletId:      uuidFromString("22222222-2222-2222-2222-222222222222"),
				OperationType: "DEPOSIT",
				Amount:        50,
			},
			mockBehavior: func(s *mock_service.MockWallet, WT wallet.WalletTransactions) {
				s.EXPECT().UpdateBalance(gomock.Any(), WT).Return(errors.New("update failed"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error":"update failed"}`,
		},
		{
			name:         "invalid JSON",
			inputBody:    `{"valletId":111,"amount":"abc"}`,
			mockBehavior: func(s *mock_service.MockWallet, WT wallet.WalletTransactions) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:      "negative amount",
			inputBody: `{"valletId":"33333333-3333-3333-3333-333333333333","operationType":"DEPOSIT","amount":-10}`,
			inputWT: wallet.WalletTransactions{
				ValletId:      uuidFromString("33333333-3333-3333-3333-333333333333"),
				OperationType: "DEPOSIT",
				Amount:        -10,
			},
			mockBehavior: func(s *mock_service.MockWallet, WT wallet.WalletTransactions) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"amount must be positive"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockWallet := mock_service.NewMockWallet(ctrl)
			if test.mockBehavior != nil {
				test.mockBehavior(mockWallet, test.inputWT)
			}

			srv := &service.Service{Wallet: mockWallet}
			h := NewHandler(srv)

			r := gin.New()
			r.POST("/wallet", h.createWalletTransaction)

			req := httptest.NewRequest("POST", "/wallet", strings.NewReader(test.inputBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			// Проверяем код и тело
			assert.Equal(t, test.expectedCode, w.Code)
			if test.expectedBody != "" {
				// Сравниваем JSON как строки
				assert.Equal(t, test.expectedBody, w.Body.String())
			}
		})
	}
}

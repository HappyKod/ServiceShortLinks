// Package midleware работа с cooke пользователя.
package midleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"HappyKod/ServiceShortLinks/internal/app/middleware"
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// WorkCooke работа с cooke пользователя.
func WorkCooke(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	userID, userCooke := generateCookie()
	md, _ := metadata.FromIncomingContext(ctx)
	values := md.Get(constans.CookeSessionName)
	if len(values) != 0 {
		_userID, b, err := validCookie(values[0])
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, middleware.ErrorValidCooke)
		}
		if b {
			userID = _userID
			userCooke = values[0]
		}
	}
	md.Delete(constans.CookeSessionName)
	md.Delete(constans.CookeUserIDName)
	md.Set(constans.CookeUserIDName, userID)
	md.Set(constans.CookeSessionName, userCooke)
	ctxNew := metadata.NewIncomingContext(ctx, md)
	return handler(ctxNew, req)
}

// validCookie проверка cooke.
func validCookie(cooke string) (string, bool, error) {
	data, err := hex.DecodeString(cooke)
	if err != nil {
		return "", false, err
	}
	if len(data[:]) < constans.CookeUserIDLen {
		return "", false, errors.New("длина cooke не соответствует требованиям")
	}
	h := hmac.New(sha256.New, constans.GlobalContainer.Get("secret-key").([]byte))
	h.Write(data[:constans.CookeUserIDLen])
	sign := h.Sum(nil)
	if hmac.Equal(sign, data[constans.CookeUserIDLen:]) {
		return string(data[:constans.CookeUserIDLen]), true, nil
	}
	return "", false, nil
}

// generateCookie генерируем новую cooke.
func generateCookie() (string, string) {
	h := hmac.New(sha256.New, constans.GlobalContainer.Get("secret-key").([]byte))
	userID := []byte(utils.GeneratorStringUUID()[:constans.CookeUserIDLen])
	h.Write(userID)
	dst := h.Sum(nil)
	var cooke []byte
	cooke = append(cooke, userID[:]...)
	cooke = append(cooke, dst...)
	return string(userID), fmt.Sprintf("%x", cooke)
}

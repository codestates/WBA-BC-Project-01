package utils

import (
	"WBA/models"
	"context"
	"encoding/hex"
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// 유저 비밀키 가져오기
// 이메일과 패스워드 받아서 개인키를 반환합니다.
func GetPrivateKey(walletColl *mongo.Collection, email string, password string) (string, error) {
	var keyjson *models.Keystores
	filter := bson.M{"email": email}
	if err := walletColl.FindOne(context.TODO(), filter).Decode(&keyjson); err != nil {
		return "", err
	}

	keyjsonToByte, err := json.Marshal(keyjson.KeyStore)
	if err != nil {
		return "", err
	}
	key, err := keystore.DecryptKey(keyjsonToByte, string(password))
	if err != nil {
		return "", err
	}
	//개인키 : hex.EncodeToString(key.PrivateKey.D.Bytes()))
	return hex.EncodeToString(key.PrivateKey.D.Bytes()), nil
}

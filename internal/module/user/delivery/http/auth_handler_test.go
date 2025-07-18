package userhttp

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/algorithm9/flash-deal/internal/config"
	"github.com/algorithm9/flash-deal/internal/module/user/repository"
	"github.com/algorithm9/flash-deal/internal/shared/entx"
	"github.com/algorithm9/flash-deal/internal/shared/idgen"
	"github.com/algorithm9/flash-deal/internal/shared/middleware"
)

func TestMockData(t *testing.T) {
	cfgPath := "../../../../../conf.toml"
	configConfig := config.LoadConfig(cfgPath)
	databaseConfig := config.ProvideDB(configConfig)
	client, cleanup, err := entx.NewEntClient(databaseConfig)
	if err != nil {
		t.Fatalf("这里不应该出错%v", err)
	}
	machine := config.ProvideMachine(configConfig)
	snowflakeIDGen, err := idgen.NewSnowflakeIDGen(machine)
	if err != nil {
		cleanup()
		t.Fatalf("这里不应该出错%v", err)
	}
	userRepository := repository.NewUserRepo(client, snowflakeIDGen)
	jwt := config.ProvideJWT(configConfig)
	authJWT, err := middleware.NewAuthJWTMiddleware(jwt, userRepository)
	if err != nil {
		cleanup()
	}

	// 创建文件
	file, err := os.Create("tokens.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	//// 写入标题
	//writer.Write([]string{"user_id", "token"})

	// 生成并写入 1000 个 token

	phone := "+8618686986910"
	for i := 0; i < 10000; i++ {
		phone = NextPhoneNumber(phone)
		hash, err := bcrypt.GenerateFromPassword([]byte("123456789"), bcrypt.DefaultCost)
		if err != nil {

		}

		id, err := snowflakeIDGen.NextID()

		client.User.
			Create().
			SetID(id).
			SetPhone(phone).
			SetPasswordHash(string(hash)).
			Save(context.Background())

		token, _, err := authJWT.TokenGenerator(&middleware.User{
			UserID: strconv.FormatUint(id, 10),
			Phone:  phone,
		})
		writer.Write([]string{token})
	}

	fmt.Println("✅ 生成完毕，tokens.csv 文件已保存")
}

func NextPhoneNumber(phone string) string {
	numberPart := phone[3:] // 去掉前缀+86
	num, err := strconv.ParseUint(numberPart, 10, 64)
	if err != nil {
		return ""
	}
	num += 1
	newNumber := fmt.Sprintf("+86%d", num)
	return newNumber
}

package authentication

import (
	context "context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"wartech-studio.com/monster-reacher/libraries/config"
	"wartech-studio.com/monster-reacher/libraries/database"
)

const ACCESS_TOKEN_LIFETIME = 5 * time.Minute
const ACCESS_TOKEN_EXTEND_AFTER = 1 * time.Minute
const ACCESS_TOKEN_EXTEND_MAX = 6

const NAME_DATABASE = "auth"
const NAME_TABLE = "token"

const HASH_DATA = "kmaADfas45asd16"

type authenticationServer struct{}

type authenticationDBSchema struct {
	ID              string `bson:"id"`
	IP              string `bson:"ip"`
	Platform        string `bson:"platform"`
	CreateTimestamp int64  `bson:"create_timestamp"`
	AccessToken     string `bson:"access_token"`
	ExtendCount     uint8  `bson:"extend_count"`
}

func NewAuthenticationServer() AuthenticationServer {
	return &authenticationServer{}
}

func (*authenticationServer) SignUp(ctx context.Context, req *SignUpRequest) (*SignUpResponse, error) {
	dbDriver := getDriver()
	defer dbDriver.Close()

	createTimestamp := time.Now().UTC().Unix()
	secretKey := getSecretKey(req.Id, req.Ip, req.Platform, createTimestamp)
	accessToken := genAccessToken(secretKey)

	_, err := dbDriver.PushOne(ctx, authenticationDBSchema{
		ID:              req.Id,
		IP:              req.Ip,
		Platform:        req.Platform,
		CreateTimestamp: createTimestamp,
		AccessToken:     accessToken,
		ExtendCount:     0,
	})

	if err != nil {
		return nil, err
	}

	log.Printf("generate new access token from id %s ip %s paltform %s token %s date %s \n", req.Id, req.Ip, req.Platform, accessToken, time.Unix(createTimestamp, 0).Format(time.UnixDate))

	return &SignUpResponse{
		AccessToken: accessToken,
	}, nil
}

func (*authenticationServer) SignIn(ctx context.Context, req *SignInRequest) (*SignInResponse, error) {
	dbDriver := getDriver()
	defer dbDriver.Close()

	result := dbDriver.SelectOne(ctx, database.MongoDBSelectOneQueryFilterOne("access_token", req.AccessToken))

	if err := database.MongoDBSelectOneResultGetError(result); err != nil {
		return &SignInResponse{IsValid: false}, nil
	}

	authData := &authenticationDBSchema{}

	if err := database.MongoDBDecodeResultToStruct(result, authData); err != nil {
		return &SignInResponse{IsValid: false}, err
	}

	timeout := authData.CreateTimestamp + int64(ACCESS_TOKEN_LIFETIME.Seconds()*float64(authData.ExtendCount))

	if time.Now().UTC().Unix()-timeout <= 0 {
		return &SignInResponse{IsValid: true}, nil
	}

	if time.Now().UTC().Unix()-timeout > int64(ACCESS_TOKEN_LIFETIME.Seconds()) {
		deleteToken(ctx, dbDriver, req.AccessToken)
		return &SignInResponse{IsValid: false}, nil
	}

	if authData.ExtendCount > ACCESS_TOKEN_EXTEND_MAX {
		deleteToken(ctx, dbDriver, req.AccessToken)
		return &SignInResponse{IsValid: false}, nil
	}

	authData.ExtendCount = authData.ExtendCount + 1

	if err := dbDriver.UpdateOne(ctx, database.MongoDBSelectOneQueryFilterOne("access_token", req.AccessToken), authData); err != nil {
		return &SignInResponse{IsValid: true}, err
	}

	return &SignInResponse{IsValid: true}, nil
}

func (*authenticationServer) SignOut(ctx context.Context, req *SignOutRequest) (*SignOutResponse, error) {
	return &SignOutResponse{}, deleteToken(ctx, nil, req.AccessToken)
}

func (*authenticationServer) mustEmbedUnimplementedAuthenticationServer() {}

func getSecretKey(id string, ip string, platform string, timestamp int64) string {
	return fmt.Sprintf(`%s-%s-%s-%d`, id, ip, platform, timestamp)
}

func genAccessToken(secretKey string) string {
	hasher := sha256.New()
	hasher.Write([]byte(secretKey))
	return hex.EncodeToString(hasher.Sum([]byte(HASH_DATA)))
}

func getDriver() database.DBDriver {
	driver, err := database.NewMongoDBDriver(config.WartechConfig().Databases["mongodb"].Host, NAME_DATABASE, NAME_TABLE)

	if err != nil {
		panic(err)
	}

	return driver
}

func deleteToken(ctx context.Context, driver database.DBDriver, token string) error {
	if driver == nil {
		driver = getDriver()
		defer driver.Close()
	}

	return driver.DeleteOne(ctx, database.MongoDBSelectOneQueryFilterOne("access_token", token))
}

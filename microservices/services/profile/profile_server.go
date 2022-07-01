package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"wartech-studio.com/monster-reacher/libraries/config"
	"wartech-studio.com/monster-reacher/libraries/database"
)

const NAME_DATABASE = "user"
const NAME_TABLE = "profile"

type profileServer struct{}

type profileDBSchema struct {
	Name        string              `bson:"name"`
	ID          string              `bson:"_id,omitempty"`
	Auth        profileAuthDBSchema `bson:"auth"`
	ServiceAuth map[string]string   `bson:"serviceAuth"`
}

type profileAuthDBSchema struct {
	User     string `bson:"user"`
	Password string `bson:"password"`
	Email    string `bson:"email"`
}

func NewProfileServer() ProfileServer {
	return &profileServer{}
}

func (*profileServer) GetData(ctx context.Context, req *GetDataRequest) (*GetDataResponse, error) {
	driver := getDriver()
	defer driver.Close()
	filter := database.MongoDBSelectOneQueryFilterOne("_id", req.GetId())
	data, err := getProfileData(ctx, driver, filter)
	b, _ := json.Marshal(*data)
	return &GetDataResponse{Data: b}, err
}
func (*profileServer) Authentication(ctx context.Context, req *AuthenticationRequest) (*SuccessResponse, error) {
	driver := getDriver()
	defer driver.Close()
	filter := database.MongoDBSelectOneQueryFilterMany([]string{"auth.user", "auth.password"}, []interface{}{req.GetUser(), req.GetPassword()})
	data, err := getProfileData(ctx, driver, filter)
	return &SuccessResponse{Success: data != nil}, err
}
func (*profileServer) AuthenticationByService(ctx context.Context, req *AuthenticationByServiceRequest) (*SuccessResponse, error) {
	driver := getDriver()
	defer driver.Close()
	filter := database.MongoDBSelectOneQueryFilterOne(fmt.Sprintf("serviceAuth.%s", req.GetName()), req.GetId())
	data, err := getProfileData(ctx, driver, filter)
	return &SuccessResponse{Success: data != nil}, err
}
func (*profileServer) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	driver := getDriver()
	defer driver.Close()
	data := &profileDBSchema{
		Auth: profileAuthDBSchema{
			User:     req.GetUser(),
			Password: req.GetPassword(),
			Email:    req.GetEmail(),
		},
	}
	result, err := driver.PushOne(ctx, data)
	if err != nil {
		return nil, err
	}
	log.Println(result)
	return &RegisterResponse{Id: database.MongoDBDecodeResultToID(result)}, nil
}
func (*profileServer) RegisterByService(ctx context.Context, req *RegisterByServiceRequest) (*RegisterResponse, error) {
	driver := getDriver()
	defer driver.Close()
	serviceAuth := make(map[string]string)
	serviceAuth[req.GetName()] = req.GetId()
	data := &profileDBSchema{
		ServiceAuth: serviceAuth,
	}
	result, err := driver.PushOne(ctx, data)
	if err != nil {
		return nil, err
	}
	log.Println(result)
	return &RegisterResponse{Id: database.MongoDBDecodeResultToID(result)}, nil
}
func (*profileServer) UserIsValid(ctx context.Context, req *UserIsValidRequest) (*SuccessResponse, error) {
	driver := getDriver()
	defer driver.Close()
	filter := database.MongoDBSelectOneQueryFilterOne("auth.user", req.GetUser())
	data, _ := getProfileData(ctx, driver, filter)
	return &SuccessResponse{Success: data != nil}, nil
}
func (*profileServer) NameIsValid(ctx context.Context, req *NameIsValidRequest) (*SuccessResponse, error) {
	driver := getDriver()
	defer driver.Close()
	filter := database.MongoDBSelectOneQueryFilterOne("name", req.GetName())
	data, _ := getProfileData(ctx, driver, filter)
	return &SuccessResponse{Success: data != nil}, nil
}
func (*profileServer) ChangeName(ctx context.Context, req *ChangeNameRequest) (*SuccessResponse, error) {
	driver := getDriver()
	defer driver.Close()
	filter := database.MongoDBSelectOneQueryFilterOne("_id", req.GetId())
	data, err := getProfileData(ctx, driver, filter)
	if err != nil {
		return nil, err
	}
	data.Name = req.GetNewName()
	err = driver.UpdateOne(ctx, filter, data)
	return &SuccessResponse{Success: err == nil}, err
}
func (*profileServer) ChangePassword(ctx context.Context, req *ChangePasswordRequest) (*SuccessResponse, error) {
	driver := getDriver()
	defer driver.Close()
	filter := database.MongoDBSelectOneQueryFilterOne("_id", req.GetId())
	data, err := getProfileData(ctx, driver, filter)
	if err != nil {
		return nil, err
	}
	data.Auth.Password = req.GetNewPassword()
	err = driver.UpdateOne(ctx, filter, data)
	return &SuccessResponse{Success: err == nil}, err
}
func (*profileServer) AddServiceAuth(ctx context.Context, req *AddServiceAuthRequest) (*SuccessResponse, error) {
	driver := getDriver()
	defer driver.Close()
	filter := database.MongoDBSelectOneQueryFilterOne("_id", req.GetId())
	data, err := getProfileData(ctx, driver, filter)
	if err != nil {
		return nil, err
	}
	data.ServiceAuth[req.GetName()] = req.GetId()
	err = driver.UpdateOne(ctx, filter, data)
	return &SuccessResponse{Success: err == nil}, err
}
func (*profileServer) RemoveServiceAuth(ctx context.Context, req *RemoveServiceAuthRequest) (*SuccessResponse, error) {
	driver := getDriver()
	defer driver.Close()
	filter := database.MongoDBSelectOneQueryFilterOne("_id", req.GetId())
	data, err := getProfileData(ctx, driver, filter)
	if err != nil {
		return nil, err
	}
	delete(data.ServiceAuth, req.GetName())
	err = driver.UpdateOne(ctx, filter, data)
	return &SuccessResponse{Success: err == nil}, err
}
func (*profileServer) MergeData(ctx context.Context, req *MergeDataRequest) (*MergeDataResponse, error) {
	driver := getDriver()
	defer driver.Close()
	filter := database.MongoDBSelectOneQueryFilterOne("_id", req.GetIdA())
	dataA, err := getProfileData(ctx, driver, filter)
	if err != nil {
		return nil, err
	}
	filter = database.MongoDBSelectOneQueryFilterOne("_id", req.GetIdB())
	dataB, err := getProfileData(ctx, driver, filter)
	if err != nil {
		return nil, err
	}

	if dataA.Auth.User == "" {
		dataA.Auth = dataB.Auth
	}
	if len(dataA.ServiceAuth) == 0 {
		dataA.ServiceAuth = dataB.ServiceAuth
	}
	err = driver.UpdateOne(ctx, filter, dataA)
	if err != nil {
		return nil, err
	}
	err = driver.DeleteOne(ctx, filter)
	return &MergeDataResponse{Id: dataA.ID}, err
}
func (*profileServer) mustEmbedUnimplementedProfileServer() {}

func getDriver() database.DBDriver {
	driver, err := database.NewMongoDBDriver(config.WartechConfig().Databases["mongodb"].Host, NAME_DATABASE, NAME_TABLE)

	if err != nil {
		panic(err)
	}
	return driver
}
func getProfileData(ctx context.Context, driver database.DBDriver, filter interface{}) (*profileDBSchema, error) {
	result := driver.SelectOne(ctx, filter)
	if err := database.MongoDBSelectOneResultGetError(result); err != nil {
		return nil, err
	}
	data := &profileDBSchema{}
	if err := database.MongoDBDecodeResultToStruct(result, data); err != nil {
		return nil, nil
	}

	return data, nil
}

package account

import (
	"context"

	pb "github.com/Asif-Faizal/Minimum-Viable-Shop/account/pb"
	"github.com/Asif-Faizal/Minimum-Viable-Shop/util"
	"google.golang.org/grpc"
)

type AccountClient struct {
	connection *grpc.ClientConn
	client     pb.AccountServiceClient
}

func NewAccountClient(url string) (*AccountClient, error) {
	connection, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &AccountClient{
		connection: connection,
		client:     pb.NewAccountServiceClient(connection),
	}, nil
}

func (client *AccountClient) Close() {
	client.connection.Close()
}

// CreateOrUpdate Account
func (client *AccountClient) CreateOrUpdateAccount(ctx context.Context, id, name, userType, email, password string) (*pb.CreateOrUpdateAccountResponse, error) {
	response, err := client.client.CreateOrUpdateAccount(ctx, &pb.CreateOrUpdateAccountRequest{
		Id:       id,
		Name:     name,
		Usertype: userType,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

// Get Account by ID
func (client *AccountClient) GetAccountByID(ctx context.Context, id string) (*pb.GetAccountByIDResponse, error) {
	response, err := client.client.GetAccountByID(ctx, &pb.GetAccountByIDRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

// List Accounts
func (client *AccountClient) ListAccounts(ctx context.Context, skip uint32, take uint32) (*pb.ListAccountsResponse, error) {
	response, err := client.client.ListAccounts(ctx, &pb.ListAccountsRequest{
		Skip: skip,
		Take: take,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (client *AccountClient) CheckEmailExists(ctx context.Context, email string) (*pb.CheckEmailExistsResponse, error) {
	response, err := client.client.CheckEmailExists(ctx, &pb.CheckEmailExistsRequest{
		Email: email,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (client *AccountClient) Login(ctx context.Context, email, password, deviceID string, deviceInfo *util.DeviceInfo) (*pb.LoginResponse, error) {
	var pbDeviceInfo *pb.DeviceInfo
	if deviceInfo != nil {
		pbDeviceInfo = &pb.DeviceInfo{
			DeviceType:      deviceInfo.DeviceType,
			DeviceModel:     deviceInfo.DeviceModel,
			DeviceOs:        deviceInfo.DeviceOS,
			DeviceOsVersion: deviceInfo.DeviceOSVersion,
			UserAgent:       deviceInfo.UserAgent,
			IpAddress:       deviceInfo.IPAddress,
		}
	}

	response, err := client.client.Login(ctx, &pb.LoginRequest{
		Email:      email,
		Password:   password,
		DeviceId:   deviceID,
		DeviceInfo: pbDeviceInfo,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (client *AccountClient) Logout(ctx context.Context, accessToken, deviceID string) (*pb.LogoutResponse, error) {
	response, err := client.client.Logout(ctx, &pb.LogoutRequest{
		AccessToken: accessToken,
		DeviceId:    deviceID,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (client *AccountClient) RefreshToken(ctx context.Context, refreshToken, deviceID string) (*pb.RefreshTokenResponse, error) {
	response, err := client.client.RefreshToken(ctx, &pb.RefreshTokenRequest{
		RefreshToken: refreshToken,
		DeviceId:     deviceID,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

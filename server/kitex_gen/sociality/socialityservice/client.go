// Code generated by Kitex v0.6.2. DO NOT EDIT.

package socialityservice

import (
	sociality "GoYin/server/kitex_gen/sociality"
	"context"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	Action(ctx context.Context, req *sociality.DouyinRelationActionRequest, callOptions ...callopt.Option) (r *sociality.DouyinRelationActionResponse, err error)
	GetRelationIdList(ctx context.Context, req *sociality.DouyinGetRelationIdListRequest, callOptions ...callopt.Option) (r *sociality.DouyinGetRelationIdListResponse, err error)
	GetSocialInfo(ctx context.Context, req *sociality.DouyinGetSocialInfoRequest, callOptions ...callopt.Option) (r *sociality.DouyinGetSocialInfoResponse, err error)
	BatchGetSocialInfo(ctx context.Context, req *sociality.DouyinBatchGetSocialInfoRequest, callOptions ...callopt.Option) (r *sociality.DouyinBatchGetSocialInfoResponse, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kSocialityServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kSocialityServiceClient struct {
	*kClient
}

func (p *kSocialityServiceClient) Action(ctx context.Context, req *sociality.DouyinRelationActionRequest, callOptions ...callopt.Option) (r *sociality.DouyinRelationActionResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Action(ctx, req)
}

func (p *kSocialityServiceClient) GetRelationIdList(ctx context.Context, req *sociality.DouyinGetRelationIdListRequest, callOptions ...callopt.Option) (r *sociality.DouyinGetRelationIdListResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetRelationIdList(ctx, req)
}

func (p *kSocialityServiceClient) GetSocialInfo(ctx context.Context, req *sociality.DouyinGetSocialInfoRequest, callOptions ...callopt.Option) (r *sociality.DouyinGetSocialInfoResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetSocialInfo(ctx, req)
}

func (p *kSocialityServiceClient) BatchGetSocialInfo(ctx context.Context, req *sociality.DouyinBatchGetSocialInfoRequest, callOptions ...callopt.Option) (r *sociality.DouyinBatchGetSocialInfoResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.BatchGetSocialInfo(ctx, req)
}

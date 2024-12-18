package shopify

import (
	"context"
	"github.com/bold-commerce/go-shopify/v4"
	"github.com/Shridhar2104/logilo/shopify/pb"
	"google.golang.org/grpc"
)

// ShopifyClient wraps the goshopify client.
type ShopifyClient struct {
	app    *goshopify.App
	client *goshopify.Client
}

type Client struct {
	conn *grpc.ClientConn
	service pb.ShopifyServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	c:= pb.NewShopifyServiceClient(conn)
	return &Client{conn: conn, service: c}, nil
}


// NewShopifyClient initializes a new Shopify client with app credentials.
func NewShopifyClient(apiKey, apiSecret, redirectURL string) *ShopifyClient {
	app := &goshopify.App{
		ApiKey:      apiKey,
		ApiSecret:   apiSecret,
		RedirectUrl: redirectURL,
		Scope:       "read_products,read_orders",
	}

	return &ShopifyClient{app: app}
}

func (c *Client) Close() {
	c.conn.Close()
}


// GetAuthorizationURL generates the OAuth URL for Shopify app installation.
func (c *ShopifyClient) GetAuthorizationURL(shopName, state string) (string, error) {
	authUrl, err := c.app.AuthorizeUrl(shopName, state)
	if err != nil {
		return "", err
	}
	return authUrl, nil
}

// GetAccessToken fetches a permanent access token after the OAuth callback.
func (c *ShopifyClient) GetAccessToken(ctx context.Context, shopName, code string) (string, error) {
	return c.app.GetAccessToken(ctx, shopName, code)
}

// InitializeAPIClient creates an authenticated API client with the access token.
func (c *ShopifyClient) InitializeAPIClient(shopName, token string) (error) {
	var err error
	c.client, err = goshopify.NewClient(*c.app, shopName, token)
	return err
}


func (c *Client) GetOrdersForShopAndAccount(ctx context.Context, shopName string, accountId string) ([]*Order, error) {
	res, err := c.service.GetOrdersForShopAndAccount(ctx, &pb.GetOrdersForShopAndAccountRequest{
		ShopName: shopName,
		AccountId: accountId,
	})
	if err != nil {
		return nil, err
	}
	orders := make([]*Order, len(res.Orders))
	return orders, nil
}

func (c *Client) SyncOrders(ctx context.Context, shopName string, sinceId string, limit int, token string) error {
	_, err:=  c.service.SyncOrders(ctx, &pb.SyncOrdersRequest{
		ShopName: shopName,
		SinceId: sinceId,
		Limit: int32(limit),
		Token: token,
	})

	return err
}

func (c *Client) UpdateOrder(ctx context.Context, shopName string, order *Order) (*Order, error) {
	res, err := c.service.UpdateOrder(ctx, &pb.UpdateOrderRequest{	
		ShopName: shopName,
		Order: &pb.Order{
			Id: order.ID,

		},
	})
	if err != nil {
		return nil, err
	}
	return &Order{
		ID: res.Order.Id,
		ShopName: shopName,
		AccountId: res.Order.AccountId,
		TotalPrice: float64(res.Order.TotalPrice),
		OrderId: res.Order.OrderId,
	}, nil
}

func (c *Client) StoreToken(ctx context.Context, shopName string, accountId string, token string) error {
	_, err := c.service.StoreToken(ctx, &pb.StoreTokenRequest{
		ShopName: shopName,
		AccountId: accountId,
		Token: token,
	})
	return err
}


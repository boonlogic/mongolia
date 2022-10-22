package mongodm

import (
	"gitlab.boonlogic.com/development/expert/mongolia/pkg/mongodm/options"
	"go.mongodb.org/mongo-driver/mongo"
)

type Client struct {
	db *mongo.Database
}

// Connect creates a new Client and then initializes it using the Connect method. This is equivalent to calling
// NewClient followed by Client.Connect.
//
// The NewClient function does not do any I/O and returns an error if the given options are invalid.
//
// The Client.Ping method can be used to verify that the deployment is successfully connected and the
// Client was correctly configured.
func Connect() {}

// NewClient creates a new client to connect to a deployment specified by the options.
func NewClient(opts *options.ClientOptions) (*Client, error) {}

func (c *Client) Connect() error {}

// Ping sends a ping command to verify that the client can connect to the deployment.
//
// Do not use Ping in production. It reduces application resiliance because applications starting up
// will error if the server is temporarily unavailable or is failing over (e.g. during autoscaling).
func (c *Client) Ping() error {}

func (c *Client) configure(opts *options.ClientOptions) error {
	// Set default options
}

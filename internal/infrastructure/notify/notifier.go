package notify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/dpe27/es-krake/internal/infrastructure/httpclient"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/nethttp"
	"github.com/dpe27/es-krake/pkg/utils"
)

type DiscordNotifier interface {
	SendMessage(msg *Message) error
}

type discordNotifier struct {
	cli        httpclient.HttpClient
	logger     *log.Logger
	webhookUrl string
}

func NewDiscordNotifier(webhookUrl string, cli httpclient.HttpClient) DiscordNotifier {
	return &discordNotifier{
		logger:     log.With("service", "discord_notifier"),
		webhookUrl: webhookUrl,
		cli:        cli,
	}
}

// SendMessage implements DiscordNotifier.
func (d *discordNotifier) SendMessage(msg *Message) error {
	if msg == nil {
		return errors.New("invalid discord webhook message")
	}

	data, err := json.Marshal(*msg)
	if err != nil {
		d.logger.Error(context.Background(), utils.ErrorMarshalFailed, "error", err)
		return err
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		d.webhookUrl,
		strings.NewReader(string(data)),
	)
	if err != nil {
		d.logger.Error(context.Background(), utils.ErrorCreateReq, "error", err)
		return err
	}

	req.Header.Add(nethttp.HeaderContentType, nethttp.MIMEApplicationJSON)
	opts := httpclient.ReqOptBuidler().
		Log().LogReqBodyOnlyError().
		LogResBody().
		LoggedResBody([]string{}).
		LoggedReqBody([]string{}).
		Build()

	res, err := d.cli.Do(req, opts)
	if err != nil {
		return nil
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			d.logger.Error(context.Background(), utils.ErrorCloseResponseBody, "error", err)
		}
	}()

	if res.StatusCode/100 != 2 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			d.logger.Error(context.Background(), "failed to read the Discord webhook response body", "error", err)
		}
		return fmt.Errorf("the Discord notification webhook failed with status %v: %v", res.StatusCode, string(body))
	}
	return nil
}

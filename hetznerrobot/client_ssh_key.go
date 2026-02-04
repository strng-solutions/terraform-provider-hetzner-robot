package hetznerrobot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type SshKeyWrapper struct {
	Key SshKey `json:"key"`
}

type SshKey struct {
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
	Type        string `json:"type"`
	Size        int    `json:"size"`
	Data        string `json:"data"`
	CreatedAt   string `json:"created_at"`
}

func (c *HetznerRobotClient) getSshKey(ctx context.Context, keyFingerprint string) (*SshKey, error) {
	bytes, err := c.makeAPICall(ctx, "GET", fmt.Sprintf("%s/key/%s", c.url, keyFingerprint), nil, []int{http.StatusOK, http.StatusCreated, http.StatusAccepted})
	if err != nil {
		return nil, err
	}

	sshKeyWrapper := SshKeyWrapper{}
	if err = json.Unmarshal(bytes, &sshKeyWrapper); err != nil {
		return nil, err
	}

	return &sshKeyWrapper.Key, nil
}

func (c *HetznerRobotClient) createSshKey(ctx context.Context, name string, data string) (*SshKey, error) {
	body := url.Values{}
	body.Set("name", name)
	body.Set("data", data)

	bytes, err := c.makeAPICall(ctx, "POST", fmt.Sprintf("%s/key", c.url), body, []int{http.StatusOK, http.StatusCreated, http.StatusAccepted})
	if err != nil {
		return nil, err
	}

	sshKeyWrapper := SshKeyWrapper{}
	if err = json.Unmarshal(bytes, &sshKeyWrapper); err != nil {
		return nil, err
	}

	return &sshKeyWrapper.Key, nil
}

func (c *HetznerRobotClient) updateSshKey(ctx context.Context, keyFingerprint string, newName string) (*SshKey, error) {
	body := url.Values{}
	body.Set("name", newName)

	bytes, err := c.makeAPICall(ctx, "PUT", fmt.Sprintf("%s/key/%s", c.url, keyFingerprint), body, []int{http.StatusOK, http.StatusCreated, http.StatusAccepted})
	if err != nil {
		return nil, err
	}

	sshKeyWrapper := SshKeyWrapper{}
	if err = json.Unmarshal(bytes, &sshKeyWrapper); err != nil {
		return nil, err
	}

	return &sshKeyWrapper.Key, nil
}

func (c *HetznerRobotClient) deleteSshKey(ctx context.Context, keyFingerprint string) error {
	_, err := c.makeAPICall(ctx, "DELETE", fmt.Sprintf("%s/key/%s", c.url, keyFingerprint), nil, []int{http.StatusOK, http.StatusCreated, http.StatusAccepted})
	if err != nil {
		return err
	}

	return nil
}

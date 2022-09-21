package etcd

import (
	"errors"
	"strconv"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	// ErrNil result nil error
	ErrNil = errors.New("result nil")
)

// PutResult is a helper that converts etcd put operation response to an error
func PutResult(_ *clientv3.PutResponse, err error) error {
	return err
}

// GetResult is a helper that converts the response of an etcd get operation to a mapa[string]string
func GetResult(resp *clientv3.GetResponse, err error) (map[string]string, error) {
	if err != nil {
		return nil, err
	}

	m := make(map[string]string, resp.Count)
	for _, v := range resp.Kvs {
		m[string(v.Key)] = string(v.Value)
	}

	return m, nil
}

// GetResultString is a helper that converts the first value in the response of an etcd get operation to a string
//
// Return ErrNil when response count euqal 0.
func GetResultString(resp *clientv3.GetResponse, err error) (string, error) {
	if err != nil {
		return "", err
	}

	if resp.Count < 1 {
		return "", ErrNil
	}

	return string(resp.Kvs[0].Value), nil
}

// GetResultStrings is a helper that converts the first value in the response of an etcd get operation to a []string
//
// Return ErrNil when response count euqal 0.
func GetResultStrings(resp *clientv3.GetResponse, err error) ([]string, error) {
	if err != nil {
		return nil, err
	}

	if resp.Count < 1 {
		return nil, ErrNil
	}

	s := make([]string, 0, resp.Count)
	for _, v := range resp.Kvs {
		s = append(s, string(v.Value))
	}

	return s, nil
}

// GetResultInt is a helper that converts the first value in the response of an etcd get operation to a int
//
// Return ErrNil when response count euqal 0.
func GetResultInt(resp *clientv3.GetResponse, err error) (int, error) {
	if err != nil {
		return 0, err
	}

	if resp.Count < 1 {
		return 0, ErrNil
	}

	return strconv.Atoi(string(resp.Kvs[0].Value))
}

// GetResultInts is a helper that converts the first value in the response of an etcd get operation to a []int
//
// Return ErrNil when response count euqal 0.
func GetResultInts(resp *clientv3.GetResponse, err error) ([]int, error) {
	if err != nil {
		return nil, err
	}

	if resp.Count < 1 {
		return nil, ErrNil
	}

	s := make([]int, 0, resp.Count)
	for _, v := range resp.Kvs {
		tmp, err := strconv.Atoi(string(v.Value))
		if err != nil {
			return nil, err
		}

		s = append(s, tmp)
	}

	return s, nil
}

// GetResultCount is a helper that return the count in the response of an etcd get operation
func GetResultCount(resp *clientv3.GetResponse, err error) (int64, error) {
	if err != nil {
		return 0, err
	}

	return resp.Count, nil
}

// DeleteResult is a helper that return the deleted raw number of etcd delete operation
func DeleteResult(resp *clientv3.DeleteResponse, err error) (int64, error) {
	if err != nil {
		return 0, err
	}

	return resp.Deleted, nil
}

// TxnResult is a helper that converts etcd txn operation response to an error
func TxnResult(_ *clientv3.TxnResponse, err error) error {
	if err != nil {
		return err
	}

	return nil
}

// LeaseGrantResult is a helper that return lease id in the response of an etcd lease grant operation
func LeaseGrantResult(resp *clientv3.LeaseGrantResponse, err error) (clientv3.LeaseID, error) {
	if err != nil {
		return 0, err
	}

	return resp.ID, nil
}

// LeaseRevokeResult is a helper that converts etcd lease revoke operation response to an error
func LeaseRevokeResult(_ *clientv3.LeaseRevokeResponse, err error) error {
	if err != nil {
		return err
	}

	return nil
}

// LeaseTimeToLiveResult is a helper that return lease ttl in the response of an etcd lease grant operation
func LeaseTimeToLiveResult(resp *clientv3.LeaseTimeToLiveResponse, err error) (int64, error) {
	if err != nil {
		return 0, err
	}

	return resp.TTL, nil
}

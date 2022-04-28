package server

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/apex/log"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/kaidyth/ender/contexts"
	pb "github.com/kaidyth/ender/protos"
	"google.golang.org/grpc"
)

var db *badger.DB

type instance struct {
	pb.UnimplementedEnderServiceServer
}

func (s *instance) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	var response pb.GetResponse

	var value, nonce []byte
	err := db.View(func(txn *badger.Txn) error {
		itemValue, err := txn.Get([]byte(in.Key + "_value"))
		if err == nil {
			err = itemValue.Value(func(val []byte) error {
				value = append([]byte{}, val...)
				return nil
			})
			if err != nil {
				return err
			}
		} else {
			return err
		}

		itemNonce, err := txn.Get([]byte(in.Key + "_nonce"))
		if err == nil {
			err = itemNonce.Value(func(val []byte) error {
				nonce = append([]byte{}, val...)
				return nil
			})
			if err != nil {
				return err
			}
		} else {
			return err
		}

		return err
	})

	response.Value = value
	response.Nonce = nonce

	return &response, err
}

func (s *instance) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	var response pb.DeleteResponse
	err := db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(in.Key))
	})

	response.Ok = (err == nil)

	return &response, err
}

func (s *instance) Exists(ctx context.Context, in *pb.ExistsRequest) (*pb.ExistsResponse, error) {
	var response pb.ExistsResponse
	request := pb.GetRequest{
		Key: in.Key,
	}
	_, err := s.Get(ctx, &request)
	response.Exists = (err == nil)

	return &response, err
}

func (s *instance) Set(ctx context.Context, in *pb.SetRequest) (*pb.SetResponse, error) {
	var response pb.SetResponse
	response.Ok = false

	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(in.Label+"_value"), in.Value)
		if err != nil {
			return err
		}
		err = txn.Set([]byte(in.Label+"_nonce"), in.Nonce)
		if err != nil {
			return err
		}

		response.Ok = true
		return nil
	})

	return &response, err
}

func NewServer(ctx context.Context, socketAddress string) error {
	if err := os.MkdirAll(socketAddress, 0700); err != nil {
		return fmt.Errorf("unable to create directory")
	}
	if err := os.RemoveAll(socketAddress); err != nil {
		return fmt.Errorf("unable to delete socket")
	}
	addr, _ := net.ResolveUnixAddr("unix", socketAddress)
	listener, err := net.ListenUnix("unix", addr)
	defer db.Close()
	os.Chmod(socketAddress, 0700)

	if err != nil {
		return err
	}

	opt := badger.DefaultOptions("").
		WithInMemory(true)
	db, err = badger.Open(opt)
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, contexts.ContextDb, db)

	defer listener.Close()
	log.Debugf("Ender chest now available at %s", socketAddress)

	server := grpc.NewServer()
	pb.RegisterEnderServiceServer(server, &instance{})

	return server.Serve(listener)
}

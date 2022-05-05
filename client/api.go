package client

// #cgo pkg-config: libsodium
// #include <stdlib.h>
// #include <sodium.h>
import "C"
import (
	"context"
	"fmt"
	"net"
	"time"
	"unsafe"

	"github.com/99designs/keyring"
	"github.com/apex/log"
	"github.com/jamesruan/sodium"
	pb "github.com/kaidyth/ender/protos"
	"google.golang.org/grpc"
)

// GenerateRandomBytes: Helper function to generate random bytes for encryption
func GenerateRandomBytes(length int) []byte {
	hash := make([]byte, length)

	C.randombytes_buf(unsafe.Pointer(&hash[0]), C.size_t(len(hash)))
	return hash
}

// getClient: Grabs the appropriate EnderServiceClient
func getClient(socketAddress string) (*pb.EnderServiceClient, error) {
	conn, err := grpc.Dial(
		socketAddress,
		grpc.WithInsecure(),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}))

	if err != nil {
		log.Error("Unable to connect to Ender")
		return nil, err
	}

	client := pb.NewEnderServiceClient(conn)

	return &client, nil
}

// getKey: Retrieves, and optionally sets the encryption key and returns it for internal use
func getKey(chest string) ([]byte, error) {
	ring, err := GetKeyring(chest)
	if err != nil {
		return nil, err
	}

	item, err := ring.Get("key")

	// If the data is empty, generate a new key and store it
	if len(item.Data) == 0 || err != nil {
		kp := sodium.MakeBoxKP()
		err = ring.Set(keyring.Item{Key: "key", Data: kp.SecretKey.Bytes})
		if err != nil {
			return nil, err
		} else {
			return kp.SecretKey.Bytes, nil
		}
	}

	return item.Data, nil
}

func Set(socketAddress string, chest string, key string, value string) (bool, error) {
	if client, err := getClient(socketAddress); err == nil {
		if k, err := getKey(chest); err == nil {

			sk := sodium.SecretBoxKey{Bytes: k}
			nonce := sodium.SecretBoxNonce{Bytes: GenerateRandomBytes(24)}
			byteData := sodium.Bytes([]byte(value))
			data := byteData.SecretBox(nonce, sk)
			request := pb.SetRequest{
				Label: key,
				Value: data,
				Nonce: nonce.Bytes,
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			result, err := (*client).Set(ctx, &request)

			if err == nil && result.Ok {
				return result.Ok, nil
			}

			return false, err
		}
	}

	return false, fmt.Errorf("unable to set key")
}

// @todo: this is not tamper safe
func Del(socketAddress string, chest string, key string) (bool, error) {
	if client, err := getClient(socketAddress); err == nil {
		var request pb.DeleteRequest
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		request.Key = key

		result, err := (*client).Delete(ctx, &request)
		if err == nil && result.Ok {
			return result.Ok, nil
		}

		return false, err
	}
	return false, fmt.Errorf("unable to delete key")
}

func Get(socketAddress string, chest string, key string) (string, error) {
	if client, err := getClient(socketAddress); err == nil {
		if k, err := getKey(chest); err == nil {
			request := pb.GetRequest{
				Key: key,
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			result, err := (*client).Get(ctx, &request)

			if err == nil {
				sk := sodium.SecretBoxKey{Bytes: k}
				v := sodium.Bytes(result.Value)
				nonce := sodium.SecretBoxNonce{Bytes: result.Nonce}
				data, err := v.SecretBoxOpen(nonce, sk)
				if err == nil {
					return string(data), nil
				}

				return "", err
			}
		}
	}

	return "", fmt.Errorf("unable to retrieve element")
}

// @todo: this is not tamper safe
func Exists(socketAddress string, chest string, key string) (bool, error) {
	if client, err := getClient(socketAddress); err == nil {
		request := pb.ExistsRequest{
			Key: key,
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		result, err := (*client).Exists(ctx, &request)
		if err == nil {
			return result.Exists, nil
		} else {
			return false, err
		}
	}
	return false, fmt.Errorf("key does not exists")
}

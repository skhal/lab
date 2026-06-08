// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package market

import (
	"context"
	"fmt"

	"github.com/skhal/lab/book/irex/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Quote calls MarketService.Quote to get a list of quotes on a symbol.
func Quote(req *pb.QuoteRequest) (*pb.QuoteResponse, error) {
	uri := fmt.Sprintf("unix://%s", defaultSocket)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient(uri, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	c := pb.NewMarketServiceClient(conn)
	res, err := c.Quote(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

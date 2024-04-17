package nft

import (
	"fmt"
	"strings"

	"github.com/CreatureDev/xrpl-go/client"
	"github.com/CreatureDev/xrpl-go/model/client/path"
	"github.com/CreatureDev/xrpl-go/model/transactions/types"
)

func SellOffersForNFT(cl *client.XRPLClient, id types.NFTokenID) ([]path.NFTokenOffer, error) {
	offersRequest := path.NFTokenSellOffersRequest{
		NFTokenID: id,
	}
	offersResponse, _, err := cl.Path.NFTokenSellOffers(&offersRequest)
	if err != nil {
		if strings.Contains(err.Error(), "objectNotFound") {
			return []path.NFTokenOffer{}, nil
		}
		return nil, fmt.Errorf("fetch sell offers: %w", err)
	}
	offers := offersResponse.Offers
	for offersResponse.Marker != nil {
		offersRequest.Marker = offersResponse.Marker
		offersResponse, _, err := cl.Path.NFTokenSellOffers(&offersRequest)
		if err != nil {
			return nil, fmt.Errorf("fetch pages sell offers: %w", err)
		}
		offers = append(offers, offersResponse.Offers...)
	}
	return offers, nil
}

func BuyOffersForNFT(cl *client.XRPLClient, id types.NFTokenID) ([]path.NFTokenOffer, error) {
	offersRequest := path.NFTokenBuyOffersRequest{
		NFTokenID: id,
	}
	offersResponse, _, err := cl.Path.NFTokenBuyOffers(&offersRequest)
	if err != nil {
		if strings.Contains(err.Error(), "objectNotFound") {
			return []path.NFTokenOffer{}, nil
		}
		return nil, fmt.Errorf("fetch sell offers: %w", err)
	}
	offers := offersResponse.Offers
	for offersResponse.Marker != nil {
		offersRequest.Marker = offersResponse.Marker
		offersResponse, _, err := cl.Path.NFTokenBuyOffers(&offersRequest)
		if err != nil {
			return nil, fmt.Errorf("fetch pages sell offers: %w", err)
		}
		offers = append(offers, offersResponse.Offers...)
	}
	return offers, nil
}

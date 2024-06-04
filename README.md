Market is a service for buying and selling digital goods using the XRPL

WIP - this project is still in active development, and most core features have not been implemented yet

![diagram](https://github.com/CreatureDev/market/blob/master/doc/diagram.png?raw=true)

XRPL Market consists of four main components
NFT Library
Publisher
Distributor
Market frontend

The NFT Library contains all the functions required to mint, sell, and purchase XRPL NFTs. This also contains utility functions for creating and validating license information encoded into NFTs.
This library, as with the rest of the project, is completely open source and free to use in any capacity.

The Publisher is used to manage products that can then be sold through a Market frontend instace, and distributed through a Distribution network. This service is used to mint NFTs for products, and manage licenses.

The Distributor is a distributed network for consumers to download purchased product through. This service validates ownership of a product through NFTs and coordinates the distribution of products to consumers. By selling access to publishers distribution network providers can earn income.

The Market frontend is what connects consumers to producers and distributors. It can be modified to accomodate the selling of any digital good or service. By brokering NFT sales between producers and consumers Market frontend providers can earn income without cutting into the profits of producers.



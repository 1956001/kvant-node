package accounts

import (
	"bytes"
	"github.com/kvant-node/core/types"
	"github.com/kvant-node/crypto"
	"github.com/kvant-node/rlp"
	"math/big"
	"sort"
)

type Model struct {
	Nonce        uint64
	MultisigData Multisig

	address  types.Address
	coins    []types.CoinSymbol
	balances map[types.CoinSymbol]*big.Int

	hasDirtyCoins bool
	dirtyBalances map[types.CoinSymbol]struct{}
	isDirty       bool // nonce or multisig data

	isNew bool

	markDirty func(types.Address)
}

type Multisig struct {
	Weights   []uint
	Threshold uint
	Addresses []types.Address
}

func (m *Multisig) Address() types.Address {
	b, err := rlp.EncodeToBytes(m)
	if err != nil {
		panic(err)
	}

	var addr types.Address
	copy(addr[:], crypto.Keccak256(b)[12:])

	return addr
}

func (m *Multisig) GetWeight(address types.Address) uint {
	for i, addr := range m.Addresses {
		if addr == address {
			return m.Weights[i]
		}
	}

	return 0
}

func (model *Model) setNonce(nonce uint64) {
	model.Nonce = nonce
	model.isDirty = true
	model.markDirty(model.address)
}

func (model *Model) getBalance(coin types.CoinSymbol) *big.Int {
	return model.balances[coin]
}

func (model *Model) hasDirtyBalances() bool {
	return len(model.dirtyBalances) > 0
}

func (model *Model) isBalanceDirty(coin types.CoinSymbol) bool {
	_, exists := model.dirtyBalances[coin]
	return exists
}

func (model *Model) getOrderedCoins() []types.CoinSymbol {
	keys := make([]types.CoinSymbol, 0, len(model.balances))
	for k := range model.balances {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return bytes.Compare(keys[i].Bytes(), keys[j].Bytes()) == 1
	})

	return keys
}

func (model *Model) setBalance(coin types.CoinSymbol, amount *big.Int) {
	if amount.Cmp(big.NewInt(0)) == 0 {
		if !model.hasCoin(coin) {
			return
		}

		var newCoins []types.CoinSymbol
		for _, c := range model.coins {
			if coin == c {
				continue
			}

			newCoins = append(newCoins, c)
		}

		model.hasDirtyCoins = true
		model.coins = newCoins
		model.balances[coin] = amount
		model.dirtyBalances[coin] = struct{}{}
		model.markDirty(model.address)

		return
	}

	if !model.hasCoin(coin) {
		model.hasDirtyCoins = true
		model.coins = append(model.coins, coin)
	}
	model.dirtyBalances[coin] = struct{}{}
	model.markDirty(model.address)
	model.balances[coin] = amount
}

func (model *Model) hasCoin(coin types.CoinSymbol) bool {
	for _, c := range model.coins {
		if c == coin {
			return true
		}
	}

	return false
}

func (model *Model) IsMultisig() bool {
	return len(model.MultisigData.Weights) > 0
}

func (model *Model) Multisig() Multisig {
	return model.MultisigData
}

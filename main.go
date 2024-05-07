package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://rpc-amoy.polygon.technology")
	if err != nil {
		fmt.Println("Erro ao conectar ao nó da rede Polygon:", err)
		return
	}
	defer client.Close()

	contractAddress := common.HexToAddress("0xcE655827f8f0F285cD23145FC0099dC574dD760c")
	privateKey, err := crypto.HexToECDSA("27be364ca6b304d810d1362a65b18efdd9a04c4bcb04e0207a8c310452bd10df")
	if err != nil {
		fmt.Println("Erro ao carregar a chave privada:", err)
		return
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("Erro ao converter a chave privada para chave pública ECDSA")
		return
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	instance, err := NewElection(contractAddress, client)
	if err != nil {
		fmt.Println("Erro ao criar instância do contrato:", err)
		return
	}

	// Registrar candidato
	candidatos := []string{"Candidato 1", "Candidato 2", "Candidato 3"}
	for _, nome := range candidatos {
		tx, err := instance.RegisterCandidate(&bind.TransactOpts{
			From:     fromAddress,
			Signer:   func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
				return types.SignTx(tx, signer, privateKey)
			},
			GasLimit: uint64(300000), 
			Value:    big.NewInt(0),  
		}, nome)
		if err != nil {
			log.Fatal("Erro ao registrar candidato:", err)
		}
		fmt.Printf("Candidato '%s' registrado. Aguardando confirmação da transação...\n", nome)
		err = waitForTransaction(client, tx)
		if err != nil {
			log.Fatal("Erro ao aguardar confirmação da transação:", err)
		}
		fmt.Println("Transação confirmada!")
	}

	// // Registrar voto
	// candidatoID := uint8(0)
	// tx, err := instance.Vote(&bind.TransactOpts{
	// 	From:     fromAddress,
	// 	Signer:   func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
	// 		return types.SignTx(tx, signer, privateKey)
	// 	},
	// 	GasLimit: uint64(300000),
	// 	Value:    big.NewInt(0),
	// }, candidatoID)
	// if err != nil {
	// 	log.Fatal("Erro ao registrar voto:", err)
	// }
	// fmt.Printf("Voto para o candidato %d registrado.
	// err = waitForTransaction(client, tx)
	// if err != nil {
	// 	log.Fatal("Erro ao aguardar confirmação da transação:", err)
	// }
	// fmt.Println("Transação confirmada!")
	// }

	// Contar votos
	// tx, err := instance.CountVotes(&bind.TransactOpts{
	//     From: fromAddress,
	//     Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
	//         return types.SignTx(tx, signer, privateKey)
	//     },
	// })

	// Auditoria
	// registros, err := instance.Audit(&bind.CallOpts{})

	// Verifique se houve erros nas operações
	// if err != nil {
	//     fmt.Println("Erro ao realizar operação no contrato:", err)
	//     return
	// }

	// Aguarde a transação ser confirmada (opcional)
	// err = waitForTransaction(client, tx)
	// if err != nil {
	//     fmt.Println("Erro ao aguardar confirmação da transação:", err)
	//     return
	// }

	// Operações bem-sucedidas
	// fmt.Println("Operação concluída com sucesso!")
	// }

	// waitForTransaction aguarda a confirmação de uma transação no nó da rede
	func waitForTransaction(client *ethclient.Client, tx *types.Transaction) error {
		ctx := context.Background()
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
	
		for {
			_, pending, err := client.TransactionByHash(ctx, tx.Hash())
			if err != nil {
				return err
			}
			if !pending {
				break
			}
			<-ticker.C
		}
	
		return nil
	}
}

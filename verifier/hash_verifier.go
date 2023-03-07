package verifier

import (
	"time"

	"github.com/bnb-chain/greenfield-challenger/common"
	"github.com/bnb-chain/greenfield-challenger/executor"

	"github.com/bnb-chain/greenfield-challenger/config"
	"github.com/bnb-chain/greenfield-challenger/db/dao"
	"github.com/bnb-chain/greenfield-challenger/db/model"
)

const batchSize = 10

type Verifier struct {
	daoManager            *dao.DaoManager
	config                *config.Config
	executor              *executor.Executor
	deduplicationInterval uint64
	heartbeatInterval     uint64
}

func NewGreenfieldHashVerifier(cfg *config.Config, dao *dao.DaoManager, executor *executor.Executor,
	deduplicationInterval, heartbeatInterval uint64) *Verifier {
	return &Verifier{
		config:                cfg,
		daoManager:            dao,
		executor:              executor,
		deduplicationInterval: deduplicationInterval,
		heartbeatInterval:     heartbeatInterval,
	}
}

func (p *Verifier) VerifyHashLoop() {
	for {
		err := p.verifyHash()
		if err != nil {
			time.Sleep(common.RetryInterval)
		}
	}
}

func (p *Verifier) verifyHash() error {
	// Read unprocessed event from db with lowest challengeId
	events, err := p.daoManager.EventDao.GetEarliestEventByStatus(model.Unprocessed, 10)
	if err != nil {
		return err
	}
	if len(events) == 0 {
		time.Sleep(common.RetryInterval)
		return nil
	}

	for _, event := range events {
		err = p.verifyForSingleEvent(event)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Verifier) verifyForSingleEvent(event *model.Event) error {
	// skip event if
	// 1) no challenger field and
	// 2) the event is not for heartbeat
	// 3) the event with same storage provider and object id has been processed recently and
	if event.ChallengerAddress == "" && event.ChallengeId%p.heartbeatInterval != 0 {
		found, err := p.daoManager.EventDao.IsEventExistsBetween(event.ObjectId, event.SpOperatorAddress,
			event.ChallengeId-p.deduplicationInterval, event.ChallengeId-1)
		if err != nil {
			return err
		}
		if found {
			return p.daoManager.UpdateEventStatusByChallengeId(event.ChallengeId, model.Duplicated)
		}
	}

	//TODO: replace with real logic, mock for test purpose
	/*
			// Call StorageProvider API to get piece hashes
			challengeInfo := sp.ChallengeInfo{
				ObjectId:        string(event.ObjectId),
				PieceIndex:      int(event.SegmentIndex),
				RedundancyIndex: int(event.RedundancyIndex),
			}
			// TODO: What to use for authinfo?
			authInfo := sp.NewAuthInfo(false, "")
			client, err := p.executor.GetGnfdClient()
			spClient := p.executor.GetSPClient()
			if err != nil {
				return err
			}
			challengeRes, err := spClient.ChallengeSP(context.Background(), challengeInfo, authInfo)
			if err != nil {
				return err
			}

			// Call blockchain for storage obj
			// TODO: Will be changed to use ObjectID instead so will have to wait
			headObjQueryReq := &storagetypes.QueryHeadObjectRequest{
				//BucketName:,
				//ObjectName:,
			}
			storageObj, err := client.StorageQueryClient.HeadObject(context.Background(), headObjQueryReq)
			if err != nil {
				return err
			}
			// Hash pieceData -> Replace pieceData hash in checksums -> Validate against original checksum stored on-chain
			// RootHash = dataHash + piecesHash
			newChecksums := storageObj.ObjectInfo.GetChecksums() // 0-6
			bytePieceData, err := ioutil.ReadAll(challengeRes.PieceData)
			if err != nil {
				return err
			}
			hashPieceData := hash.CalcSHA256(bytePieceData)
			newChecksums[event.SegmentIndex] = hashPieceData
			total := bytes.Join(newChecksums, []byte(""))
			rootHash := []byte(hash.CalcSHA256Hex(total))

			if bytes.Equal(rootHash, storageObj.ObjectInfo.Checksums[event.RedundancyIndex+1]) {
				return p.daoManager.EventDao.UpdateEventStatusVerifyResultByChallengeId(event.ChallengeId, model.Verified, model.HashMismatched)
			}

		return p.daoManager.EventDao.UpdateEventStatusVerifyResultByChallengeId(event.ChallengeId, model.Verified, model.HashMatched)
	*/
	if event.ChallengeId%3 == 0 {
		return p.daoManager.EventDao.UpdateEventStatusVerifyResultByChallengeId(event.ChallengeId, model.Verified, model.HashMismatched)
	}
	return p.daoManager.EventDao.UpdateEventStatusVerifyResultByChallengeId(event.ChallengeId, model.Verified, model.HashMatched)
}

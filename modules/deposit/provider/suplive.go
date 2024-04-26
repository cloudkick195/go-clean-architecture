package depositProvider

import (
	"bytes"
	"encoding/json"
	"errors"
	"go_clean_architecture/commons"
	"go_clean_architecture/commons/models"
	"go_clean_architecture/utils"
	"net/http"
	"strconv"
	"time"
)

type ConfigSupliveProvider struct {
	Uid           int64
	Key           string
	Api           string
	Commission    float64
	ComisionRefer float64
}

type SupliveProvider struct {
	Config       ConfigSupliveProvider
	Model        *models.Config
	depositModel *models.Deposit
}

func NewSupliveProvider(m *models.Config, depositModel *models.Deposit) IProvider {
	var config ConfigSupliveProvider
	_ = json.Unmarshal([]byte(m.Config), &config)

	return &SupliveProvider{
		Model:        m,
		Config:       config,
		depositModel: depositModel,
	}
}
func (p *SupliveProvider) Prepare(input *models.Member) {
	price := int64(4)
	p.depositModel.CoinAmount = p.depositModel.Amount / price
	p.depositModel.PriceCoin = price
	if input.IsAgency {
		p.depositModel.Commission = int64(float64(p.depositModel.Amount) * p.Config.Commission)
	}
	p.depositModel.TotalCoin = p.depositModel.CoinAmount
}
func (p *SupliveProvider) GetComissionRefer() int64 {
	return int64(float64(p.depositModel.Amount) * p.Config.ComisionRefer)
}

func (p *SupliveProvider) Deposit() (string, interface{}, error) {
	num, _ := strconv.ParseInt(p.depositModel.ProviderID, 10, 64)

	reqDataProvider := MoneyMerchantExternalPayModel{
		Uid:       p.Config.Uid,
		TargetId:  num,
		Num:       p.depositModel.TotalCoin,
		EventTime: time.Now().Unix(),
	}
	byteReqData, _ := json.Marshal(reqDataProvider)
	byteReqDataStr := string(byteReqData)
	cipherByte, err := utils.RSAEncrypt(byteReqDataStr, p.Config.Key)

	if err != nil {
		return byteReqDataStr, nil, commons.ErrInvalidRequest(err)
	}

	// Tạo chuỗi JSON
	reqSupermeet := &ReqSupermeet{
		Message: cipherByte,
		Uid:     p.Config.Uid,
	}

	jsonStr, _ := json.Marshal(reqSupermeet)
	apiURL := p.Config.Api
	resp, err := http.Post(apiURL, "application/json", bytes.NewReader(jsonStr))
	if err != nil {
		return byteReqDataStr, nil, commons.ErrInvalidRequest(err)
	}
	defer resp.Body.Close()

	var convertResData ResSupermeet
	err = json.NewDecoder(resp.Body).Decode(&convertResData)
	if err != nil {
		return byteReqDataStr, err, commons.ErrInvalidRequest(err)
	}

	if convertResData.DmError != 0 {
		if convertResData.DmError == 4101 {
			return byteReqDataStr, convertResData, models.ErrWrongDepositID
		}
		return byteReqDataStr, convertResData, commons.ErrInvalidRequest(errors.New("cannot deposit coins"))
	}

	return byteReqDataStr, convertResData, nil
}

package go3s

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

var (
	Err401 = fmt.Errorf("solscan: 401 unauthorized")
	Err403 = fmt.Errorf("solscan: 403 forbidden")
	Err404 = fmt.Errorf("solscan: 404 not found")
	Err429 = fmt.Errorf("solscan: 429 too many requests")
	Err500 = fmt.Errorf("solscan: 500 internal server error")
)

func createParams[Opt any](optParams *Opt, requiredParams ...string) url.Values {
	if optParams == nil {
		optParams = new(Opt)
	}
	params := toParams(*optParams)
	if len(requiredParams)%2 != 0 {
		panic("solscan: requiredParams is key,value,key,value... list")
	}
	for i := 0; i < len(requiredParams); i += 2 {
		k := requiredParams[i]
		v := requiredParams[i+1]
		params.Add(k, v)
	}
	return params
}

type Getter[D any] interface {
	Do(ctx context.Context) (D, error)
	URL() string
}

func ExportBodyUnmarshal(body []byte) ([]byte, error) {
	return body, nil
}

func DefaultRespBodyUnmarshal[D any](body []byte) (D, error) {
	var respData RespData[D]
	err := json.Unmarshal(body, &respData)
	if err != nil {
		return *new(D), err
	}
	return respData.Data, nil
}

func DefaultRespStatusHandler(resp *http.Response) error {
	statusCode := resp.StatusCode
	switch statusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("solscan: 400 bad request, but can not read body: %s", err.Error())
		}
		var respError RespError
		err = json.Unmarshal(body, &respError)
		if err != nil {
			return fmt.Errorf("solscan: 400 bad request, but can not unmarshal body: %s", err.Error())
		}
		return &respError.Errors
	case http.StatusUnauthorized:
		return Err401
	case http.StatusForbidden:
		return Err403
	case http.StatusNotFound:
		return Err404
	case http.StatusTooManyRequests:
		return Err429
	case http.StatusInternalServerError:
		return Err500
	default:
		return fmt.Errorf("solscan: %d unknown error", statusCode)
	}
}

type GetterOption struct {
	RetryInterval time.Duration
	MaxRetries    int
}

var defaultGetterOption = &GetterOption{
	RetryInterval: time.Second,
	MaxRetries:    1,
}

type SimpleGetter[D any] struct {
	BaseURL           string
	Path              string
	Params            url.Values
	Headers           map[string][]string
	Limiter           *rate.Limiter
	RespStatusHandler func(resp *http.Response) error
	RespBodyUnmarshal func(body []byte) (D, error)
	Option            *GetterOption
}

func (g *SimpleGetter[D]) URL() string {
	return fmt.Sprintf("%s/%s?%s", g.BaseURL, strings.Trim(g.Path, "/"), g.Params.Encode())
}

func (g *SimpleGetter[D]) Do(ctx context.Context) (D, error) {
	option := g.Option
	if option == nil {
		option = defaultGetterOption
	}
	retryInterval := option.RetryInterval
	maxRetries := option.MaxRetries
	if retryInterval == 0 {
		retryInterval = time.Second
	}
	if maxRetries == 0 {
		maxRetries = 1
	}
	if maxRetries == 1 {
		return g.do(ctx)
	}
	for i := 0; i < maxRetries; i++ {
		d, err := g.do(ctx)
		if err != nil {
			time.Sleep(retryInterval)
			logger.Error("solscan: failed to get response", "retrying", i+1, "error", err)
			continue
		}
		return d, nil
	}
	return *new(D), fmt.Errorf("solscan: failed to get response after %d retries", maxRetries)
}

func (g *SimpleGetter[D]) do(ctx context.Context) (D, error) {
	if g.Limiter != nil {
		err := g.Limiter.Wait(ctx)
		if err != nil {
			return *new(D), err
		}
	}
	ul := fmt.Sprintf("%s/%s", g.BaseURL, strings.Trim(g.Path, "/"))
	if len(g.Params) > 0 {
		ul += "?" + g.Params.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, "GET", ul, nil)
	if err != nil {
		return *new(D), err
	}
	for k, v := range g.Headers {
		req.Header[k] = v
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return *new(D), err
	}
	defer resp.Body.Close()

	if g.RespStatusHandler != nil {
		err = g.RespStatusHandler(resp)
	} else {
		err = DefaultRespStatusHandler(resp)
	}

	if err != nil {
		return *new(D), err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return *new(D), fmt.Errorf("solscan: can not read body: %s", err.Error())
	}

	if g.RespBodyUnmarshal != nil {
		return g.RespBodyUnmarshal(body)
	}

	return DefaultRespBodyUnmarshal[D](body)
}

func CreateSliceDataFinishChecker[D any](pageSize int64) func(d []D) bool {
	return func(d []D) bool {
		return len(d) < int(pageSize)
	}
}

func CreateWithTotalItemsFinishChecker[D any](pageSize int64) func(d RespDataWithTotal[D]) bool {
	return func(d RespDataWithTotal[D]) bool {
		return len(d.Items) < int(pageSize)
	}
}

func CreateWithTotalTransactionsDataFinishChecker(pageSize int64) func(d RespDataWithTotal[Transaction]) bool {
	return func(d RespDataWithTotal[Transaction]) bool {
		return len(d.Transactions) < int(pageSize)
	}
}

func CreateWithTotalDataFinishChecker[D any](pageSize int64) func(d RespDataWithTotal[D]) bool {
	return func(d RespDataWithTotal[D]) bool {
		return len(d.Data) < int(pageSize)
	}
}

func CreateSliceResultsHandler[D any](totalSize int64) func([][]D) ([]D, error) {
	return func(results [][]D) ([]D, error) {
		l := 0
		for _, result := range results {
			l += len(result)
		}
		r := make([]D, 0, l)
		for _, result := range results {
			r = append(r, result...)
		}
		return r[:int(math.Min(float64(totalSize), float64(len(r))))], nil
	}
}

func CreateWithTotalItemsResultsHandler[D any](totalSize int64) func([]RespDataWithTotal[D]) (RespDataWithTotal[D], error) {
	return func(results []RespDataWithTotal[D]) (RespDataWithTotal[D], error) {
		l := 0
		for _, result := range results {
			l += len(result.Items)
		}
		total := int64(0)
		r := make([]D, 0, l)
		for _, result := range results {
			r = append(r, result.Items...)
			if result.Total > total {
				total = result.Total
			}
		}
		return RespDataWithTotal[D]{
			Items: r[:int(math.Min(float64(totalSize), float64(len(r))))],
			Total: total,
		}, nil
	}
}

func CreateWithTotalTransactionsResultsHandler(totalSize int64) func([]RespDataWithTotal[Transaction]) (RespDataWithTotal[Transaction], error) {
	return func(results []RespDataWithTotal[Transaction]) (RespDataWithTotal[Transaction], error) {
		l := 0
		for _, result := range results {
			l += len(result.Transactions)
		}
		total := int64(0)
		r := make([]Transaction, 0, l)
		for _, result := range results {
			r = append(r, result.Transactions...)
			if result.Total > total {
				total = result.Total
			}
		}
		return RespDataWithTotal[Transaction]{
			Transactions: r[:int(math.Min(float64(totalSize), float64(len(r))))],
			Total:        total,
		}, nil
	}
}

func CreateWithTotalDataResultsHandler[D any](totalSize int64) func([]RespDataWithTotal[D]) (RespDataWithTotal[D], error) {
	return func(results []RespDataWithTotal[D]) (RespDataWithTotal[D], error) {
		l := 0
		for _, result := range results {
			l += len(result.Data)
		}
		total := int64(0)
		r := make([]D, 0, l)
		for _, result := range results {
			r = append(r, result.Data...)
			if result.Total > total {
				total = result.Total
			}
		}
		return RespDataWithTotal[D]{
			Data:  r[:int(math.Min(float64(totalSize), float64(len(r))))],
			Total: total,
		}, nil
	}
}

type CcrtDataFinishChecker[D any] func(d D) bool

type CcrtResultsHandler[D any] func(results []D) (D, error)

type CcrtGetter[D any] struct {
	Getters           []Getter[D]
	MaxConcurrency    int64
	DataFinishChecker CcrtDataFinishChecker[D]
	ResultsHandler    CcrtResultsHandler[D]
}

func (g *CcrtGetter[D]) URL() string {
	return ""
}

func (g *CcrtGetter[D]) Do(ctx context.Context) (D, error) {
	if g.MaxConcurrency == 0 {
		g.MaxConcurrency = 1
	}
	l := len(g.Getters)
	results := make([]D, l)
	var isRespEmpty bool
	for i := 0; i < l; i += int(g.MaxConcurrency) {
		end := i + int(g.MaxConcurrency)
		if end > l {
			end = l
		}
		logger.Info("solscan: concurrency", "start", i, "end", end, "total", l, "url", g.Getters[i].URL())
		group := g.Getters[i:end]
		eg, ctx := errgroup.WithContext(ctx)
		for j, simpleGetter := range group {
			j := j
			simpleGetter := simpleGetter
			eg.Go(func() error {
				d, err := simpleGetter.Do(ctx)
				if err != nil {
					return err
				}
				if g.DataFinishChecker != nil && g.DataFinishChecker(d) {
					isRespEmpty = true
				}
				results[i+j] = d
				return nil
			})
		}
		err := eg.Wait()
		if err != nil {
			return *new(D), err
		}
		if isRespEmpty {
			break
		}
	}
	return g.ResultsHandler(results)
}

type PagingParams[D any] struct {
	StartPage         int64
	TotalSize         int64
	MaxConcurrency    int64
	DataFinishChecker CcrtDataFinishChecker[D]
	ResultsHandler    CcrtResultsHandler[D]
}

type PagingGetter[D any] struct {
	BaseURL      string
	Path         string
	Params       url.Values
	Headers      map[string][]string
	Limiter      *rate.Limiter
	GetterOption *GetterOption
	PagingParams *PagingParams[D]
}

func (g *PagingGetter[D]) URL() string {
	return ""
}

func (g *PagingGetter[D]) Do(ctx context.Context) (D, error) {
	if g.PagingParams == nil {
		sg := SimpleGetter[D]{
			BaseURL: g.BaseURL,
			Path:    g.Path,
			Params:  g.Params,
			Headers: g.Headers,
			Limiter: g.Limiter,
		}
		return sg.Do(ctx)
	}
	pageSize, err := strconv.ParseInt(g.Params.Get("page_size"), 10, 64)
	if err != nil {
		return *new(D), fmt.Errorf("solscan: can not get page_size: %s", err.Error())
	}
	pages := int64(math.Ceil(float64(g.PagingParams.TotalSize) / float64(pageSize)))
	getters := make([]Getter[D], pages)
	for i := int64(0); i < pages; i++ {
		p := url.Values{}
		for k, v := range g.Params {
			for _, vv := range v {
				p.Add(k, vv)
			}
		}
		p.Set("page", strconv.FormatInt(i+g.PagingParams.StartPage, 10))
		sg := &SimpleGetter[D]{
			BaseURL: g.BaseURL,
			Path:    g.Path,
			Params:  p,
			Headers: g.Headers,
			Limiter: g.Limiter,
			Option:  g.GetterOption,
		}
		getters[i] = sg
	}
	ccrt := &CcrtGetter[D]{
		Getters:           getters,
		MaxConcurrency:    g.PagingParams.MaxConcurrency,
		DataFinishChecker: g.PagingParams.DataFinishChecker,
		ResultsHandler:    g.PagingParams.ResultsHandler,
	}
	return ccrt.Do(ctx)
}

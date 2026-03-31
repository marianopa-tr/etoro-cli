package api

// --- Market Data ---

type InstrumentSearchResponse struct {
	Page       int          `json:"page"`
	PageSize   int          `json:"pageSize"`
	TotalItems int          `json:"totalItems"`
	Items      []Instrument `json:"items"`
}

type Instrument struct {
	InstrumentID           int     `json:"instrumentId"`
	DisplayName            string  `json:"displayname"`
	Symbol                 string  `json:"internalSymbolFull"`
	InstrumentType         string  `json:"instrumentType"`
	InstrumentTypeID       int     `json:"instrumentTypeID"`
	ExchangeID             int     `json:"exchangeID"`
	IsOpen                 bool    `json:"isOpen"`
	IsCurrentlyTradable    bool    `json:"isCurrentlyTradable"`
	IsExchangeOpen         bool    `json:"isExchangeOpen"`
	IsBuyEnabled           bool    `json:"isBuyEnabled"`
	IsDelisted             bool    `json:"isDelisted"`
	IsHiddenFromClient     bool    `json:"isHiddenFromClient"`
	CurrentRate            float64 `json:"currentRate"`
	DailyPriceChange       float64 `json:"dailyPriceChange"`
	WeeklyPriceChange      float64 `json:"weeklyPriceChange"`
	MonthlyPriceChange     float64 `json:"monthlyPriceChange"`
	ThreeMonthPriceChange  float64 `json:"threeMonthPriceChange"`
	OneYearPriceChange     float64 `json:"oneYearPriceChange"`
	Popularity7Day         int     `json:"popularityUniques7Day"`
	Logo50x50              string  `json:"logo50x50"`
	InternalAssetClassName string  `json:"internalAssetClassName"`
	InternalExchangeName   string  `json:"internalExchangeName"`
	InternalClosingPrice   float64 `json:"internalClosingPrice"`
}

type InstrumentDisplayData struct {
	InstrumentID       int    `json:"instrumentID"`
	DisplayName        string `json:"instrumentDisplayName"`
	Symbol             string `json:"symbolFull"`
	InstrumentTypeID   int    `json:"instrumentTypeID"`
	ExchangeID         int    `json:"exchangeID"`
	PriceSource        string `json:"priceSource"`
	StocksIndustryID   int    `json:"stocksIndustryID"`
	HasExpirationDate  bool   `json:"hasExpirationDate"`
	IsInternalInstrument bool `json:"isInternalInstrument"`
}

type InstrumentsResponse struct {
	InstrumentDisplayDatas []InstrumentDisplayData `json:"instrumentDisplayDatas"`
}

type LiveRatesResponse struct {
	Rates []Rate `json:"rates"`
}

type Rate struct {
	InstrumentID  int     `json:"instrumentID"`
	Ask           float64 `json:"ask"`
	Bid           float64 `json:"bid"`
	LastExecution float64 `json:"lastExecution"`
	ConversionAsk float64 `json:"conversionRateAsk"`
	ConversionBid float64 `json:"conversionRateBid"`
	Date          string  `json:"date"`
}

type ExchangesResponse struct {
	ExchangeInfo []Exchange `json:"exchangeInfo"`
}

type Exchange struct {
	ExchangeID int    `json:"exchangeID"`
	Name       string `json:"name"`
}

type InstrumentTypesResponse struct {
	InstrumentTypes []InstrumentTypeInfo `json:"instrumentTypes"`
}

type InstrumentTypeInfo struct {
	InstrumentTypeID int    `json:"instrumentTypeID"`
	Name             string `json:"name"`
}

// --- Portfolio ---

type PortfolioResponse struct {
	ClientPortfolio ClientPortfolio `json:"clientPortfolio"`
}

type ClientPortfolio struct {
	Positions     []Position     `json:"positions"`
	Credit        float64        `json:"credit"`
	Mirrors       []Mirror       `json:"mirrors"`
	Orders        []Order        `json:"orders"`
	OrdersForOpen []OrderForOpen `json:"ordersForOpen"`
	BonusCredit   float64        `json:"bonusCredit"`
	UnrealizedPnL float64        `json:"unrealizedPnL"`
}

type Position struct {
	PositionID    int     `json:"positionId"`
	CID           int     `json:"cid"`
	OpenDateTime  string  `json:"openDateTime"`
	OpenRate      float64 `json:"openRate"`
	InstrumentID  int     `json:"instrumentId"`
	IsBuy         bool    `json:"isBuy"`
	TakeProfitRate float64 `json:"takeProfitRate"`
	StopLossRate  float64 `json:"stopLossRate"`
	Amount        float64 `json:"amount"`
	Leverage      int     `json:"leverage"`
	Units         float64 `json:"units"`
	TotalFees     float64 `json:"totalFees"`
	InitialAmount float64 `json:"initialAmountInDollars"`
	MirrorID      int     `json:"mirrorId"`
}

type Mirror struct {
	MirrorID         int        `json:"mirrorId"`
	CID              int        `json:"cid"`
	ParentCID        int        `json:"parentCid"`
	StopLossPercent  float64    `json:"stopLossPercentage"`
	IsPaused         bool       `json:"isPaused"`
	AvailableAmount  float64    `json:"availableAmount"`
	StopLossAmount   float64    `json:"stopLossAmount"`
	InitialInvestment float64   `json:"initialInvestment"`
	Positions        []Position `json:"positions"`
	ParentUsername   string     `json:"parentUsername"`
	ClosedNetProfit  float64    `json:"closedPositionsNetProfit"`
	StartedCopyDate  string     `json:"startedCopyDate"`
}

type Order struct {
	OrderID       int     `json:"orderId"`
	InstrumentID  int     `json:"instrumentId"`
	IsBuy         bool    `json:"isBuy"`
	Amount        float64 `json:"amount"`
	Units         float64 `json:"units"`
	Leverage      int     `json:"leverage"`
	Rate          float64 `json:"rate"`
	StopLossRate  float64 `json:"stopLossRate"`
	TakeProfitRate float64 `json:"takeProfitRate"`
	OpenDateTime  string  `json:"openDateTime"`
	IsTslEnabled  bool    `json:"isTslEnabled"`
}

type OrderForOpen struct {
	OrderID       int     `json:"orderId"`
	InstrumentID  int     `json:"instrumentId"`
	IsBuy         bool    `json:"isBuy"`
	Amount        float64 `json:"amount"`
	AmountInUnits float64 `json:"amountInUnits"`
	Leverage      int     `json:"leverage"`
	StopLossRate  float64 `json:"stopLossRate"`
	TakeProfitRate float64 `json:"takeProfitRate"`
	IsTslEnabled  bool    `json:"isTslEnabled"`
	OpenDateTime  string  `json:"openDateTime"`
}

// --- Trading ---

type OpenByAmountRequest struct {
	InstrumentID   int      `json:"InstrumentID"`
	IsBuy          bool     `json:"IsBuy"`
	Leverage       int      `json:"Leverage"`
	Amount         float64  `json:"Amount"`
	StopLossRate   *float64 `json:"StopLossRate,omitempty"`
	TakeProfitRate *float64 `json:"TakeProfitRate,omitempty"`
	IsTslEnabled   *bool    `json:"IsTslEnabled,omitempty"`
}

type OpenByUnitsRequest struct {
	InstrumentID   int      `json:"InstrumentID"`
	IsBuy          bool     `json:"IsBuy"`
	Leverage       int      `json:"Leverage"`
	Units          float64  `json:"AmountInUnits"`
	StopLossRate   *float64 `json:"StopLossRate,omitempty"`
	TakeProfitRate *float64 `json:"TakeProfitRate,omitempty"`
	IsTslEnabled   *bool    `json:"IsTslEnabled,omitempty"`
}

type LimitOrderRequest struct {
	InstrumentID   int      `json:"InstrumentID"`
	IsBuy          bool     `json:"IsBuy"`
	Leverage       int      `json:"Leverage"`
	Amount         float64  `json:"Amount,omitempty"`
	Units          float64  `json:"AmountInUnits,omitempty"`
	Rate           float64  `json:"Rate"`
	StopLossRate   *float64 `json:"StopLossRate,omitempty"`
	TakeProfitRate *float64 `json:"TakeProfitRate,omitempty"`
	IsTslEnabled   *bool    `json:"IsTslEnabled,omitempty"`
	IsNoStopLoss   *bool    `json:"IsNoStopLoss,omitempty"`
	IsNoTakeProfit *bool    `json:"IsNoTakeProfit,omitempty"`
}

type ClosePositionRequest struct {
	UnitsToDeduct *float64 `json:"UnitsToDeduct,omitempty"`
}

type OrderResponse struct {
	OrderID    int    `json:"orderId"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

// --- Trade History ---

type TradeHistoryEntry struct {
	NetProfit        float64 `json:"netProfit"`
	CloseRate        float64 `json:"closeRate"`
	CloseTimestamp   string  `json:"closeTimestamp"`
	PositionID       int64   `json:"positionId"`
	InstrumentID     int     `json:"instrumentId"`
	IsBuy            bool    `json:"isBuy"`
	Leverage         int     `json:"leverage"`
	OpenRate         float64 `json:"openRate"`
	OpenTimestamp    string  `json:"openTimestamp"`
	StopLossRate     float64 `json:"stopLossRate"`
	TakeProfitRate   float64 `json:"takeProfitRate"`
	Investment       float64 `json:"investment"`
	InitialInvestment float64 `json:"initialInvestment"`
	Fees             float64 `json:"fees"`
	Units            float64 `json:"units"`
}

// --- Watchlists ---

type WatchlistsResponse struct {
	Watchlists []Watchlist `json:"watchlists"`
}

type Watchlist struct {
	WatchlistID          string          `json:"WatchlistId"`
	Name                 string          `json:"Name"`
	WatchlistType        string          `json:"WatchlistType"`
	TotalItems           int             `json:"TotalItems"`
	IsDefault            bool            `json:"IsDefault"`
	IsUserSelectedDefault bool           `json:"IsUserSelectedDefault"`
	WatchlistRank        int             `json:"WatchlistRank"`
	Items                []WatchlistItem `json:"Items"`
}

type WatchlistItem struct {
	ItemID   int    `json:"ItemId"`
	ItemType string `json:"ItemType"`
	ItemRank int    `json:"ItemRank"`
}

type CreateWatchlistRequest struct {
	Name  string          `json:"Name"`
	Items []WatchlistItem `json:"Items,omitempty"`
}

type AddWatchlistItemsRequest struct {
	Items []WatchlistItem `json:"Items"`
}

type RemoveWatchlistItemsRequest struct {
	Items []WatchlistItem `json:"Items"`
}

type CuratedListsResponse struct {
	CuratedLists []CuratedList `json:"CuratedLists"`
}

type CuratedList struct {
	UUID        string            `json:"Uuid"`
	Name        string            `json:"Name"`
	Description string            `json:"Description"`
	ImageURL    string            `json:"ListImageUrl"`
	Items       []CuratedListItem `json:"Items"`
}

type CuratedListItem struct {
	InstrumentID int `json:"InstrumentId"`
}

// --- Users / Social ---

type UserSearchResponse struct {
	TotalCount int          `json:"totalCount"`
	Items      []UserSummary `json:"items"`
}

type UserSummary struct {
	UserName     string  `json:"userName"`
	DisplayName  string  `json:"displayName"`
	CID          int     `json:"cid"`
	Gain         float64 `json:"gain"`
	RiskScore    int     `json:"riskScore"`
	Copiers      int     `json:"copiers"`
	AUM          float64 `json:"aum"`
	IsFund       bool    `json:"isFund"`
	IsPopularInvestor bool `json:"isPopularInvestor"`
}

type UserGainResponse struct {
	Monthly []GainEntry `json:"monthly"`
	Yearly  []GainEntry `json:"yearly"`
}

type GainEntry struct {
	Timestamp string  `json:"timestamp"`
	Gain      float64 `json:"gain"`
}

type TradeInfoResponse struct {
	UserName              string  `json:"userName"`
	FullName              string  `json:"fullName"`
	WeeksSinceRegistration int   `json:"weeksSinceRegistration"`
	IsPopularInvestor     bool    `json:"isPopularInvestor"`
	IsFund                bool    `json:"isFund"`
	Gain                  float64 `json:"gain"`
	DailyGain             float64 `json:"dailyGain"`
	ThisWeekGain          float64 `json:"thisWeekGain"`
	RiskScore             int     `json:"riskScore"`
	MaxDailyRiskScore     int     `json:"maxDailyRiskScore"`
	MaxMonthlyRiskScore   int     `json:"maxMonthlyRiskScore"`
	Copiers               int     `json:"copiers"`
	CopiedTrades          int     `json:"copiedTrades"`
	CopyTradesPct         float64 `json:"copyTradesPct"`
	Trades                int     `json:"trades"`
	WinRatio              float64 `json:"winRatio"`
	DailyDD               float64 `json:"dailyDd"`
	WeeklyDD              float64 `json:"weeklyDd"`
	PeakToValley          float64 `json:"peakToValley"`
	ProfitableWeeksPct    float64 `json:"profitableWeeksPct"`
	ProfitableMonthsPct   float64 `json:"profitableMonthsPct"`
	AvgPosSize            float64 `json:"avgPosSize"`
	AUMTier               int     `json:"aumTier"`
	AUMTierDesc           string  `json:"aumTierDesc"`
}

type CopiersResponse struct {
	Copiers []CopierInfo `json:"copiers"`
}

type CopierInfo struct {
	Gender               string `json:"Gender"`
	Club                 string `json:"Club"`
	Country              string `json:"Country"`
	CopyStartedCategory  string `json:"CopyStartedAtCategory"`
	AmountCategory       string `json:"AmountCategory"`
}

// --- Feeds ---

type FeedResponse struct {
	Discussions []FeedPost     `json:"discussions"`
	Pagination  FeedPagination `json:"pagination"`
	Paging      FeedPagination `json:"paging"`
}

type FeedPost struct {
	ID   string   `json:"id"`
	Post PostData `json:"post"`
}

type PostData struct {
	ID          string    `json:"id"`
	Owner       PostOwner `json:"owner"`
	Message     PostMessage `json:"message"`
	Created     string    `json:"created"`
	Type        string    `json:"type"`
}

type PostOwner struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type PostMessage struct {
	Text string `json:"text"`
}

type FeedPagination struct {
	Next   string `json:"next"`
	Offset int    `json:"offSet"`
	Take   int    `json:"take"`
}

type CreatePostRequest struct {
	Owner    int              `json:"owner"`
	Message  string           `json:"message"`
	Tags     *PostTagsRequest `json:"tags,omitempty"`
}

type PostTagsRequest struct {
	Instruments []PostTagInstrument `json:"instruments,omitempty"`
}

type PostTagInstrument struct {
	ID int `json:"id"`
}

// --- Agent Portfolios ---

type CreateAgentPortfolioRequest struct {
	InvestmentAmountInUsd     float64  `json:"investmentAmountInUsd"`
	AgentPortfolioName        string   `json:"agentPortfolioName"`
	AgentPortfolioDescription string   `json:"agentPortfolioDescription,omitempty"`
	UserTokenName             string   `json:"userTokenName"`
	ScopeIDs                  []int    `json:"scopeIds"`
	IPsWhitelist              []string `json:"ipsWhitelist,omitempty"`
	ExpiresAt                 string   `json:"expiresAt,omitempty"`
}

type CreateAgentPortfolioResponse struct {
	AgentPortfolioID             string                         `json:"agentPortfolioId"`
	AgentPortfolioName           string                         `json:"agentPortfolioName"`
	AgentPortfolioGCID           int                            `json:"agentPortfolioGcid"`
	AgentPortfolioVirtualBalance float64                        `json:"agentPortfolioVirtualBalance"`
	MirrorID                     int                            `json:"mirrorId"`
	UserTokens                   []CreateAgentPortfolioToken    `json:"userTokens"`
	UserTokenCreated             *bool                          `json:"userTokenCreated,omitempty"`
}

type CreateAgentPortfolioToken struct {
	UserTokenID   string   `json:"userTokenId"`
	UserToken     string   `json:"userToken"`
	UserTokenName string   `json:"userTokenName"`
	ClientID      string   `json:"clientId"`
	IPsWhitelist  []string `json:"ipsWhitelist"`
	ScopeIDs      []int    `json:"scopeIds"`
	ExpiresAt     string   `json:"expiresAt,omitempty"`
}

type GetAgentPortfoliosResponse struct {
	AgentPortfolios []AgentPortfolioItem `json:"agentPortfolios"`
}

type AgentPortfolioItem struct {
	AgentPortfolioID             string                    `json:"agentPortfolioId"`
	AgentPortfolioName           string                    `json:"agentPortfolioName"`
	AgentPortfolioGCID           int                       `json:"agentPortfolioGcid"`
	AgentPortfolioVirtualBalance float64                   `json:"agentPortfolioVirtualBalance"`
	MirrorID                     int                       `json:"mirrorId"`
	CreatedAt                    string                    `json:"createdAt"`
	UserTokens                   []AgentPortfolioTokenItem `json:"userTokens"`
}

type AgentPortfolioTokenItem struct {
	UserTokenID             string   `json:"userTokenId"`
	UserTokenName           string   `json:"userTokenName"`
	ClientID                string   `json:"clientId"`
	ExternalApplicationName string   `json:"externalApplicationName"`
	IPsWhitelist            []string `json:"ipsWhitelist"`
	ExpiresAt               string   `json:"expiresAt,omitempty"`
	ScopeIDs                []int    `json:"scopeIds"`
	CreatedAt               string   `json:"createdAt"`
}

type CreateUserTokenRequest struct {
	UserTokenName string   `json:"userTokenName"`
	ScopeIDs      []int    `json:"scopeIds"`
	IPsWhitelist  []string `json:"ipsWhitelist,omitempty"`
	ExpiresAt     string   `json:"expiresAt,omitempty"`
}

type CreateUserTokenResponse struct {
	UserTokenID string `json:"userTokenId"`
	UserToken   string `json:"userToken"`
}

type UpdateUserTokenRequest struct {
	ScopeIDs     []int    `json:"scopeIds,omitempty"`
	IPsWhitelist []string `json:"ipsWhitelist,omitempty"`
	ExpiresAt    string   `json:"expiresAt,omitempty"`
}

// --- User Profile ---

type UserProfileResponse struct {
	Users []UserProfile `json:"users"`
}

type UserProfile struct {
	GCID     int    `json:"gcid"`
	RealCID  int    `json:"realCID"`
	DemoCID  int    `json:"demoCID"`
	Username string `json:"username"`
	IsPi     bool   `json:"isPi"`
}

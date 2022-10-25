package facade

// // AnyTime is mock for time
// type AnyTime struct{}
//
// // Match for AnyTime mocking
// func (a AnyTime) Match(v driver.Value) bool {
//	_, ok := v.(time.Time)
//	return ok
// }
//
// type HandlerFacadeTestSuite struct {
//	suite.Suite
//	mock       *sqlmock.Sqlmock
//	db         *sql.DB
//	gormDB     *gorm.DB
//	v          *validator.Validate
//	client     *database.Client
//	logger     *zap.Logger
//	logic      *HandlerFacade
//	r          *mux.Router
//	redisCache *redis.Client
//	redisMock  *redismock.ClientMock
// }
//
// func (suite *HandlerFacadeTestSuite) SetupSuite() {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		suite.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	suite.mock = &mock
//	suite.db = db
//	dialector := postgres.New(postgres.Config{
//		DSN:                  "sqlmock_db_0",
//		DriverName:           "postgres",
//		Conn:                 suite.db,
//		PreferSimpleProtocol: true,
//	})
//	gormDB, err := gorm.Open(dialector, &gorm.Config{})
//	if err != nil {
//		suite.T().Fatalf("an error '%s' was not expected when init gorm", err)
//	}
//	if gormDB == nil {
//		suite.T().Fatalf("gormDB is null")
//	}
//	suite.gormDB = gormDB
//
//	v := validator.New()
//	_ = v.RegisterValidation("checkURL", func(fl validator.FieldLevel) bool {
//		arr, ok := fl.Field().Interface().([]string)
//		if !ok {
//			return false
//		}
//		for _, a := range arr {
//			_, err := url.ParseRequestURI(a)
//			if err != nil {
//				return false
//			}
//		}
//		return true
//	})
//	suite.v = v
//
//	suite.client = database.NewClient(suite.gormDB)
//
//	config := zap.NewDevelopmentConfig()
//	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
//	logger, _ := config.Build()
//	suite.logger = logger
//
//	redisDB, redisMock := redismock.NewClientMock()
//	suite.redisCache = redis.NewClientForTest(redisDB)
//	suite.redisMock = &redisMock
// }
//
// func (suite *HandlerFacadeTestSuite) TearDownSuite() {
//	_ = suite.db.Close()
// }
//
// func (suite *HandlerFacadeTestSuite) SetupTest() {
//	suite.logic = NewHandlerFacade(suite.redisCache, suite.client, suite.v, suite.logger)
//
//	getAdHandler, _ := suite.logic.GetHandler("get_ad")
//	listAdsHandler, _ := suite.logic.GetHandler("list_ads")
//	createAdHandler, _ := suite.logic.GetHandler("create_ad")
//
//	suite.r = mux.NewRouter()
//
//	suite.r.HandleFunc("/ads/{id:[0-9]+}", getAdHandler).Methods("GET")
//	suite.r.HandleFunc("/ads", listAdsHandler).Methods("GET")
//	suite.r.HandleFunc("/ads", createAdHandler).Methods("POST")
// }
//
// func (suite *HandlerFacadeTestSuite) TestCreateAd() {
//	id := 0
//	expectedResult := entities.CreateAdAnswer{Status: "success", ID: &id}
//	item := entities.AdJSONItem{Title: "title", Price: decimal.NewFromInt32(0)}
//	mockedRow := sqlmock.NewRows([]string{"id"}).AddRow(strconv.Itoa(id))
//
//	dbItem := database.AdItem{
//		Title: item.Title, Description: item.Description,
//		Price: item.Price, ImageURLs: []database.ImageURL{}, MainImageURL: nil}
//	dbCachedItem := dbItem.CreateMap([]string{"description", "image_urls"})
//	bs, _ := easyjson.Marshal(dbCachedItem)
//	(*suite.redisMock).ExpectSet(fmt.Sprintf("item:%d", id), string(bs), suite.redisCache.GetDuration()).SetVal("OK")
//
//	(*suite.mock).ExpectBegin()
//	(*suite.mock).ExpectQuery(regexp.QuoteMeta(
//		`INSERT INTO "ad_items" ("created_at","updated_at","deleted_at","title","description","price")
//			VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id","id"`)).
//		WithArgs(AnyTime{}, AnyTime{}, nil, item.Title, item.Description, item.Price).
//		WillReturnRows(mockedRow)
//	(*suite.mock).ExpectCommit()
//	(*suite.mock).ExpectClose()
//
//	defer func() {
//		_ = suite.logger.Sync()
//	}()
//
//	h, _ := suite.logic.GetHandler("create_ad")
//
//	bs, err := json.Marshal(item)
//	if err != nil {
//		suite.T().Fatalf("err during Marshal item; item = %v ; err = %v", item, err)
//	}
//	req := httptest.NewRequest("POST", "/ads", bytes.NewBuffer(bs))
//	res := httptest.NewRecorder()
//
//	h(res, req)
//	resp := res.Result()
//	if resp.StatusCode != http.StatusOK {
//		suite.T().Fatalf("resp.StatusCode != 200; resp.StatusCode = %d", resp.StatusCode)
//	}
//
//	body, _ := io.ReadAll(resp.Body)
//	var actualResult entities.CreateAdAnswer
//	_ = json.Unmarshal(body, &actualResult)
//
//	if !reflect.DeepEqual(actualResult, expectedResult) {
//		suite.T().Fatalf("actualResult != expectedResult; actualResult = %v ; expectedResult = %v",
//			actualResult, expectedResult)
//	}
// }
//
// func (suite *HandlerFacadeTestSuite) TestCreateAdInvalidStruct() {
//	item := entities.AdJSONItem{Price: decimal.NewFromInt32(0)}
//
//	defer func() {
//		_ = suite.logger.Sync()
//	}()
//
//	h, _ := suite.logic.GetHandler("create_ad")
//
//	bs, err := json.Marshal(item)
//	if err != nil {
//		suite.T().Fatalf("err during Marshal item; item = %v ; err = %v", item, err)
//	}
//	req := httptest.NewRequest("POST", "/ads", bytes.NewBuffer(bs))
//	res := httptest.NewRecorder()
//
//	h(res, req)
//	resp := res.Result()
//	if resp.StatusCode != http.StatusUnprocessableEntity {
//		suite.T().Fatalf("resp.StatusCode != 422; resp.StatusCode = %d", resp.StatusCode)
//	}
// }
//
// func (suite *HandlerFacadeTestSuite) TestCreateAdDBError() {
//	item := entities.AdJSONItem{Title: "title", Price: decimal.NewFromInt32(0)}
//	(*suite.mock).ExpectCommit()
//
//	defer func() {
//		_ = suite.logger.Sync()
//	}()
//
//	h, _ := suite.logic.GetHandler("create_ad")
//
//	bs, err := json.Marshal(item)
//	if err != nil {
//		suite.T().Fatalf("err during Marshal item; item = %v ; err = %v", item, err)
//	}
//	req := httptest.NewRequest("POST", "/ads", bytes.NewBuffer(bs))
//	res := httptest.NewRecorder()
//
//	h(res, req)
//	resp := res.Result()
//	if resp.StatusCode != http.StatusInternalServerError {
//		suite.T().Fatalf("resp.StatusCode != 500; resp.StatusCode = %d", resp.StatusCode)
//	}
// }
//
// func (suite *HandlerFacadeTestSuite) TestGetAdNotFound() {
//	defer func() {
//		_ = suite.logger.Sync()
//	}()
//
//	// h, _ := suite.logic.GetHandler("get_ad")
//
//	req := httptest.NewRequest("GET", "/ads/problem/", nil)
//	res := httptest.NewRecorder()
//
//	suite.r.ServeHTTP(res, req)
//	resp := res.Result()
//	if resp.StatusCode != http.StatusNotFound {
//		suite.T().Fatalf("resp.StatusCode != 404; resp.StatusCode = %d", resp.StatusCode)
//	}
// }
//
// func (suite *HandlerFacadeTestSuite) TestGetAdNegativeID() {
//	defer func() {
//		_ = suite.logger.Sync()
//	}()
//
//	req := httptest.NewRequest("GET", "/ads/-1", nil)
//	res := httptest.NewRecorder()
//
//	suite.r.ServeHTTP(res, req)
//	resp := res.Result()
//	if resp.StatusCode != http.StatusNotFound {
//		suite.T().Fatalf("resp.StatusCode != 404; resp.StatusCode = %d", resp.StatusCode)
//	}
// }
//
// func (suite *HandlerFacadeTestSuite) TestGetAdInvalidFields() {
//	defer func() {
//		_ = suite.logger.Sync()
//	}()
//
//	req := httptest.NewRequest("GET", "/ads/1?fields=wrong", nil)
//	res := httptest.NewRecorder()
//
//	suite.r.ServeHTTP(res, req)
//	resp := res.Result()
//	if resp.StatusCode != http.StatusUnprocessableEntity {
//		suite.T().Fatalf("resp.StatusCode != 422; resp.StatusCode = %d", resp.StatusCode)
//	}
// }
//
// func (suite *HandlerFacadeTestSuite) TestListAdsStringOffset() {
//	defer func() {
//		_ = suite.logger.Sync()
//	}()
//
//	req := httptest.NewRequest("GET", "/ads?offset=wrong", nil)
//	res := httptest.NewRecorder()
//
//	suite.r.ServeHTTP(res, req)
//	resp := res.Result()
//	if resp.StatusCode != http.StatusUnprocessableEntity {
//		suite.T().Fatalf("resp.StatusCode != 422; resp.StatusCode = %d", resp.StatusCode)
//	}
// }
//
// func (suite *HandlerFacadeTestSuite) TestListAdsNegativeOffset() {
//	defer func() {
//		_ = suite.logger.Sync()
//	}()
//
//	req := httptest.NewRequest("GET", "/ads?offset=-1", nil)
//	res := httptest.NewRecorder()
//
//	suite.r.ServeHTTP(res, req)
//	resp := res.Result()
//	if resp.StatusCode != http.StatusUnprocessableEntity {
//		suite.T().Fatalf("resp.StatusCode != 422; resp.StatusCode = %d", resp.StatusCode)
//	}
// }
//
// func (suite *HandlerFacadeTestSuite) TestListAdsInvalidBy() {
//	defer func() {
//		_ = suite.logger.Sync()
//	}()
//
//	req := httptest.NewRequest("GET", "/ads?by=-42", nil)
//	res := httptest.NewRecorder()
//
//	suite.r.ServeHTTP(res, req)
//	resp := res.Result()
//	if resp.StatusCode != http.StatusUnprocessableEntity {
//		suite.T().Fatalf("resp.StatusCode != 422; resp.StatusCode = %d", resp.StatusCode)
//	}
// }
//
// func (suite *HandlerFacadeTestSuite) TestListAdsInvalidAsc() {
//	defer func() {
//		_ = suite.logger.Sync()
//	}()
//
//	req := httptest.NewRequest("GET", "/ads?asc=-42", nil)
//	res := httptest.NewRecorder()
//
//	suite.r.ServeHTTP(res, req)
//	resp := res.Result()
//	if resp.StatusCode != http.StatusUnprocessableEntity {
//		suite.T().Fatalf("resp.StatusCode != 422; resp.StatusCode = %d", resp.StatusCode)
//	}
// }
//
// func TestHandlerFacadeTestSuite(t *testing.T) {
//	suite.Run(t, new(HandlerFacadeTestSuite))
// }

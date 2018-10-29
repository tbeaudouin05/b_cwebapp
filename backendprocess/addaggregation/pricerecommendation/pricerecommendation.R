
Sys.setenv(SPARK_HOME = "D:/program/spark-2.2.0-bin-hadoop2.7")

.libPaths(c(file.path(Sys.getenv("SPARK_HOME"), "R", "lib"), .libPaths()))

library(SparkR)

sparkR.session(
  appName = "PricingRecommendation"
  #, spark.master = "spark://192.168.10.178:7077"
  , sparkConfig = list(
    spark.driver.memory = "8g"
    , spark.executor.memory= "8g"
    , spark.mongodb.input.uri = "mongodb://127.0.0.1/competition_analysis.test?readPreference=primaryPreferred"
    , spark.mongodb.output.uri = "mongodb://127.0.0.1/competition_analysis.test"
    , spark.mongodb.output.replaceDocument="false"
  )
  , sparkPackages = "org.mongodb.spark:mongo-spark-connector_2.11:2.2.5"
)

# get BmlCatalogConfig to filter for good_match
BmlCatalogConfig <- read.df("",
  source = "com.mongodb.spark.sql.DefaultSource"
  , database = "competition_analysis"
  , collection = "bml_catalog_config"
  , pipeline = "[{'$match': {'$and':[{'good_match':{'$exists':'true'}}, {'avg_price':{'$gt':0}}]}},
  {'$project': {'_id':'$_id'
  ,'fk_dgk_catalog_config':'$fk_dgk_catalog_config'
  ,'avg_price':'$avg_price'
  ,'avg_special_price':'$avg_special_price'}}]"
)
createOrReplaceTempView(BmlCatalogConfig, "BmlCatalogConfig")
BmlCatalogConfigR <- collect(BmlCatalogConfig)

# create IDBmlCatalogConfigRStr to only filter good_match in history
IDBmlCatalogConfigR <- BmlCatalogConfigR$`_id`
IDBmlCatalogConfigRStr = ""
for (i in IDBmlCatalogConfigR) {
  IDBmlCatalogConfigRStr <- paste0(IDBmlCatalogConfigRStr,i,",")
}
IDBmlCatalogConfigRStr = substr(IDBmlCatalogConfigRStr,1,nchar(IDBmlCatalogConfigRStr)-1)
IDBmlCatalogConfigRStr = gsub(" ","",IDBmlCatalogConfigRStr)

# get BmlCatalogConfigHist to analyze each SKU history
BmlCatalogConfigHist <- read.df("",
  source = "com.mongodb.spark.sql.DefaultSource"
  , database = "competition_analysis"
  , collection = "bml_catalog_config_hist"
)
createOrReplaceTempView(BmlCatalogConfigHist, "BmlCatalogConfigHist")

# filter for good_match, last 90 days & sum_of_paid_price > 0 and group by week & SKU
# and create price
BmlCatalogConfigHistF <- sql(paste0("SELECT 
                                    fk_bml_catalog_config _id
                                    , WEEKOFYEAR(config_snapshot_at) week
                                    , AVG(CASE WHEN avg_special_price > 0
                                          THEN avg_special_price ELSE avg_price END) price
                                    , SUM(sum_of_paid_price) sum_of_paid_price
                                    
                                    FROM BmlCatalogConfigHist
                                    WHERE config_snapshot_at > date_sub(current_timestamp(),90)
                                    AND sum_of_paid_price > 0                                     -- make sure there was at least one sale
                                    AND avg_price > 0                                             -- make sure at least avg_price exists
                                    AND fk_bml_catalog_config IN(",IDBmlCatalogConfigRStr,")
                                    GROUP BY 
                                    WEEKOFYEAR(config_snapshot_at)
                                    , fk_bml_catalog_config"))
createOrReplaceTempView(BmlCatalogConfigHistF, "BmlCatalogConfigHistF")

# rank each week of each SKU by sum_of_paid_price
# filter only for best Bamilo week ie rank = 1
# join with BmlCatalogConfig to get fk_dgk_catalog_config
# create fk_dgk_catalog_config_week to join with DgkWithIDDgkWeek (to compare prices of Bamilo vs. Digikala during Bamilo best week)
BmlWithIDDgkWeek <- sql("SELECT
                      tmp._id
                      , CONCAT(bcc.fk_dgk_catalog_config, tmp.week) fk_dgk_catalog_config_week
                      , tmp.price
                      , tmp.sum_of_paid_price

                        FROM (SELECT
                        _id
                        , week
                        , price
                        , sum_of_paid_price
                        , DENSE_RANK() OVER (PARTITION BY CONCAT(_id,week) ORDER BY sum_of_paid_price DESC) rank
                        FROM BmlCatalogConfigHistF) tmp

                      LEFT JOIN BmlCatalogConfig bcc
                      ON tmp._id = bcc._id AND tmp.rank = 1")
createOrReplaceTempView(BmlWithIDDgkWeek, "BmlWithIDDgkWeek")

# create IDDgkCatalogConfigRStr to only filter good_match in dgk history
IDDgkCatalogConfigR <- BmlCatalogConfigR$fk_dgk_catalog_config
IDDgkCatalogConfigRStr = ""
for (i in IDDgkCatalogConfigR) {
  IDDgkCatalogConfigRStr <- paste0(IDDgkCatalogConfigRStr,"'",i,"'",",")
}
IDDgkCatalogConfigRStr = substr(IDDgkCatalogConfigRStr,1,nchar(IDDgkCatalogConfigRStr)-1)
IDDgkCatalogConfigRStr = gsub(" ","",IDDgkCatalogConfigRStr)

# get DgkCatalogConfigHist to analyze each dgk SKU history
DgkCatalogConfigHist <- read.df("",
                                source = "com.mongodb.spark.sql.DefaultSource"
                                , database = "competition_analysis"
                                , collection = "dgk_catalog_config_hist")
createOrReplaceTempView(DgkCatalogConfigHist, "DgkCatalogConfigHist")

# filter for good_match & last 90 days and group by week & SKU
# and create fk_dgk_catalog_config_week & dgk_price
DgkWithIDDgkWeek <- sql(paste0("SELECT 
                              CONCAT(fk_dgk_catalog_config, WEEKOFYEAR(config_snapshot_at)) fk_dgk_catalog_config_week
                              , AVG(CASE WHEN avg_special_price > 0
                                    THEN avg_special_price ELSE avg_price END) dgk_price

                              FROM DgkCatalogConfigHist
                              WHERE config_snapshot_at > date_sub(current_timestamp(),90)
                              AND avg_special_price > 0                                    -- make sure at least avg_price exists
                              AND fk_dgk_catalog_config IN(",IDDgkCatalogConfigRStr,")
                              GROUP BY 
                              WEEKOFYEAR(config_snapshot_at)
                              , fk_dgk_catalog_config"))
createOrReplaceTempView(DgkWithIDDgkWeek, "DgkWithIDDgkWeek")

# join BmlWithIDDgkWeek and DgkWithIDDgkWeek on fk_dgk_catalog_config_week
# to compare prices of Bamilo vs. Digikala during the best week of Bamilo
# INNER JOIN because if we don't have bamilo price then there is no valuable information, same if we don't have digikala price for that week
BmlDgkDiff <- sql("SELECT
            bw._id
            , 1 - bw.price/dw.dgk_price diff_price_best_week
            FROM BmlWithIDDgkWeek bw
            INNER JOIN DgkWithIDDgkWeek dw
            ON bw.fk_dgk_catalog_config_week = dw.fk_dgk_catalog_config_week")
createOrReplaceTempView(BmlDgkDiff, "BmlDgkDiff")

# get DgkCatalogConfig to get the latest price of Digikala
DgkCatalogConfig <- read.df("",
                            source = "com.mongodb.spark.sql.DefaultSource"
                            , database = "competition_analysis"
                            , collection = "dgk_catalog_config"
                            , pipeline = "{'$project': {'_id':'$_id'
                              ,'avg_price':'$avg_price'
                              ,'avg_special_price':'$avg_special_price'}}")
createOrReplaceTempView(DgkCatalogConfig, "DgkCatalogConfig")

# join BmlCatalogConfig, DgkCatalogConfig and BmlDgkDiff on IDBmlCatalogConfig
# calculate bml current_price and dgk_current_price thanks to BmlCatalogConfig and DgkCatalogConfig
# create is_not_good_price_match to check if current_price and dgk_current_price are too different (> 60% difference, in which case we assume the SKUs are not good matches for this process)
# create recommended_price
BmlJoinDgkJoinBestWeekDiff <- sql("SELECT
                                bcc._id

                                , CASE WHEN bcc.avg_special_price > 0
                                       THEN bcc.avg_special_price ELSE bcc.avg_price END current_price

                                , CASE WHEN dcc.avg_special_price > 0
                                       THEN dcc.avg_special_price ELSE dcc.avg_price END dgk_current_price

                                -- if abs(1 - current_price / dgk_current_price) > 0.6 then is_not_good_price_match = TRUE
                                , CASE WHEN abs(1 - (CASE WHEN bcc.avg_special_price > 0
                                       THEN bcc.avg_special_price ELSE bcc.avg_price END) / (CASE WHEN dcc.avg_special_price > 0
                                       THEN dcc.avg_special_price ELSE dcc.avg_price END)) > 0.6
                                       THEN TRUE ELSE FALSE END is_not_good_price_match

                                -- recommended_price = dgk_current_price * (1 - diff_price_best_week)
                                , (CASE WHEN dcc.avg_special_price > 0
                                       THEN dcc.avg_special_price ELSE dcc.avg_price END) * (1 - bdd.diff_price_best_week) recommended_price

                                FROM BmlCatalogConfig bcc
                                JOIN DgkCatalogConfig dcc
                                ON bcc.fk_dgk_catalog_config = dcc._id
                                JOIN BmlDgkDiff bdd
                                ON bcc._id = bdd._id")
createOrReplaceTempView(BmlJoinDgkJoinBestWeekDiff, "BmlJoinDgkJoinBestWeekDiff")

# we recommend a price only if the recommended price is lower than dgk price AND different than bml price (otherwise, it's already the good price)
# Also, if the products do not have a good price match (ie if dgk price vs. bml price are too different), then We assume that there was a matching mistake 
BmlReco <- sql("SELECT
             _id
            , bround(recommended_price,0) recommended_price
             FROM BmlJoinDgkJoinBestWeekDiff
             WHERE recommended_price < dgk_current_price
             AND bround(recommended_price,0) <> current_price
             AND is_not_good_price_match = FALSE")
createOrReplaceTempView(BmlReco, "BmlReco")

BmlNotGoodPriceMatch <- sql("SELECT
                         _id
                        , is_not_good_price_match
                         FROM BmlJoinDgkJoinBestWeekDiff
                         WHERE is_not_good_price_match = TRUE")


write.df(BmlReco, "",
  source = "com.mongodb.spark.sql.DefaultSource",
  mode = "append", database = "competition_analysis", collection = "bml_catalog_config"
)

write.df(BmlNotGoodPriceMatch, "",
         source = "com.mongodb.spark.sql.DefaultSource",
         mode = "append", database = "competition_analysis", collection = "bml_catalog_config"
)

sparkR.session.stop()
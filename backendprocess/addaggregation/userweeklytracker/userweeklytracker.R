
  Sys.setenv(SPARK_HOME='D:/program/spark-2.2.0-bin-hadoop2.7')

.libPaths(c(file.path(Sys.getenv("SPARK_HOME"),"R","lib"),.libPaths()))

library(SparkR)

# NOT to replace documents, you added this option: spark.mongodb.output.replaceDocument="false"
sparkR.session(appName = "userweeklytracker", sparkConfig = list(spark.driver.memory = "1g"
                                                    ,spark.mongodb.input.uri="mongodb://127.0.0.1/competition_analysis.dgk_catalog_config?readPreference=primaryPreferred"
                                                    ,spark.mongodb.output.uri="mongodb://127.0.0.1/test.myCollection"
                                                    ,spark.mongodb.output.replaceDocument="false")
               , sparkPackages = "org.mongodb.spark:mongo-spark-connector_2.11:2.2.5"
)

user <- read.df("", source = "com.mongodb.spark.sql.DefaultSource", database ="competition_analysis", collection = "user")
createOrReplaceTempView(user, "user")

email <- sql("SELECT DISTINCT _id, email FROM user")

bmlCatalogConfig <- read.df("", source = "com.mongodb.spark.sql.DefaultSource"
                            , database ="competition_analysis"
                            , collection = "bml_catalog_config"
                            , pipeline = "{'$match': {'matched_by_email': {'$exists':'true'}}}")
createOrReplaceTempView(bmlCatalogConfig, "bmlCatalogConfig")

TotalSKUMatch <- sql("SELECT matched_by_email AS email2, count(_id) AS total_matched_sku FROM bmlCatalogConfig GROUP BY matched_by_email")
Last7DaySKUMatch <- sql("SELECT matched_by_email AS email3, count(_id) AS last_7_day_matched_sku 
                        FROM bmlCatalogConfig 
                        WHERE good_match_at > date_sub(current_timestamp(),7)
                        GROUP BY matched_by_email")
Last14DaySKUMatch <- sql("SELECT matched_by_email AS email4, count(_id) AS last_14_day_matched_sku 
                        FROM bmlCatalogConfig 
                        WHERE good_match_at > date_sub(current_timestamp(),14) AND good_match_at <= date_sub(current_timestamp(),7)
                        GROUP BY matched_by_email")

BmlSupplierConfigCountTopPageLast <- read.df("", source = "com.mongodb.spark.sql.DefaultSource"
                                  , database ="competition_analysis"
                                  , collection = "bml_agg_statistic_hist"
                                  , pipeline = "{'$match': {'matched_by_email': {'$exists':'true'}, 'type': 'BmlSupplierConfigCountTopPageLast'}}")
createOrReplaceTempView(BmlSupplierConfigCountTopPageLast, "BmlSupplierConfigCountTopPageLast")

TotalSupplierMatch <- sql("SELECT matched_by_email AS email5, count(_id) AS total_matched_supplier FROM BmlSupplierConfigCountTopPageLast GROUP BY matched_by_email")
Last7DaySupplierMatch <- sql("SELECT matched_by_email AS email6, count(_id) AS last_7_day_matched_supplier 
                        FROM BmlSupplierConfigCountTopPageLast 
                        WHERE good_match_at > date_sub(current_timestamp(),7)
                        GROUP BY matched_by_email")
Last14DaySupplierMatch <- sql("SELECT matched_by_email AS email7, count(_id) AS last_14_day_matched_supplier 
                        FROM BmlSupplierConfigCountTopPageLast 
                        WHERE good_match_at > date_sub(current_timestamp(),14) AND good_match_at <= date_sub(current_timestamp(),7)
                        GROUP BY matched_by_email")

NotCompetitiveGoodMatched <- read.df("", source = "com.mongodb.spark.sql.DefaultSource"
                                  , database ="competition_analysis"
                                  , collection = "bml_dgk_agg_statistic_hist"
                                  , pipeline = "{'$match': {'matched_by_email': {'$exists':'true'}, 'type': 'NotCompetitiveGoodMatched'}}")
createOrReplaceTempView(NotCompetitiveGoodMatched, "NotCompetitiveGoodMatched")

TotalNonCompetitiveSkuMatch <- sql("SELECT matched_by_email AS email8, count(_id) AS total_not_competitive FROM NotCompetitiveGoodMatched GROUP BY matched_by_email")

createOrReplaceTempView(email, "email")
createOrReplaceTempView(TotalSKUMatch, "TotalSKUMatch")
createOrReplaceTempView(Last7DaySKUMatch, "Last7DaySKUMatch")
createOrReplaceTempView(Last14DaySKUMatch, "Last14DaySKUMatch")
createOrReplaceTempView(TotalSupplierMatch, "TotalSupplierMatch")
createOrReplaceTempView(Last7DaySupplierMatch, "Last7DaySupplierMatch")
createOrReplaceTempView(Last14DaySupplierMatch, "Last14DaySupplierMatch")
createOrReplaceTempView(TotalNonCompetitiveSkuMatch, "TotalNonCompetitiveSkuMatch")

UserTableForMongoUpsert <- sql("SELECT
                        email._id
                        ,email.email
                        ,COALESCE(TotalSKUMatch.total_matched_sku,0) total_matched_sku
                        ,COALESCE(Last7DaySKUMatch.last_7_day_matched_sku,0) last_7_day_matched_sku
                        ,COALESCE(Last14DaySKUMatch.last_14_day_matched_sku,0) last_14_day_matched_sku
                        ,COALESCE(TotalSupplierMatch.total_matched_supplier,0) total_matched_supplier
                        ,COALESCE(Last7DaySupplierMatch.last_7_day_matched_supplier,0) last_7_day_matched_supplier
                        ,COALESCE(Last14DaySupplierMatch.last_14_day_matched_supplier,0) last_14_day_matched_supplier
                        ,COALESCE(TotalNonCompetitiveSkuMatch.total_not_competitive,0) total_not_competitive
                        FROM email
                        LEFT JOIN TotalSKUMatch
                        ON TotalSKUMatch.email2 = email.email
                        LEFT JOIN Last7DaySKUMatch
                        ON Last7DaySKUMatch.email3 = email.email
                        LEFT JOIN Last14DaySKUMatch
                        ON Last14DaySKUMatch.email4 = email.email
                        LEFT JOIN TotalSupplierMatch
                        ON TotalSupplierMatch.email5 = email.email
                        LEFT JOIN Last7DaySupplierMatch
                        ON Last7DaySupplierMatch.email6 = email.email
                        LEFT JOIN Last14DaySupplierMatch
                        ON Last14DaySupplierMatch.email7 = email.email
                        LEFT JOIN TotalNonCompetitiveSkuMatch
                        ON TotalNonCompetitiveSkuMatch.email8 = email.email")

# head(UserTableForMongoUpsert)

write.df(UserTableForMongoUpsert, "", source = "com.mongodb.spark.sql.DefaultSource",
         mode = "append", database = "competition_analysis", collection = "user")



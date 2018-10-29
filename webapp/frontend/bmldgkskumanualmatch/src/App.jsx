import React, { Component } from 'react';
import Table from './table/table.jsx';
import FilterCategory1 from './filter/filtercategory/filtercategory1.jsx';
import FilterCategory2 from './filter/filtercategory/filtercategory2.jsx';
import SearchSkuName from './filter/searchskuname.jsx';
import FliterPriceRange from './filter/pricerange.jsx';
import PageNumberWrapper from './pagination/pagenumberwrapper.jsx';
import ShowBamiloSKU from './sku/sku.jsx'
import UpdateFrequency from './sku/updatefrequency.jsx'
import ShowMatchedSKU from './sku/matched.jsx'
import ProductUrl from './addurl/addproducturl.jsx'
import envConfig from './envConfig.json'


export default class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            table: [],// data.table,
            filterCategory1OptionList: [], // filterCategory2OptionList,
            filterCategory2OptionList: [],
            ws: null,
            pageNumber: 1,
            pageNumberList: [], // data.pageNumberList,
            FrequencyOptionList: [],
            category1: "All",
            category2: 'All',
            skuName: '',
            BmlSkuName: '',
            BmlPrice: '',
            BmlImgLink: '',
            BmlBrand: '',
            BmlSKULink: '',
            BmlID: "173027",
            dgkFK: 0,
            checkCycle: 0,
            UpdateStatus: " ",
            MatchingStatus: " ",
            DgkImgLink: "",
            minPrice: '-5',
            maxPrice: '9999999999',
            UrlValidationStatus: [],
            User: {
                Email: null,
                Name: null,
            }


        } // , BmlSkuName, BmlPrice, BmlImgLink, BmlBrand, BmlSKULink
    }


    componentDidMount() {

        this.setState({ BmlID: document.getElementById("BmlIDCatalogConfig").getAttribute("value") })
        this.setState({ DgkImgLink: document.getElementById("DgkImgLink").getAttribute("value") })
        this.setState({
            User: {
                Email: document.getElementById("Email").getAttribute("value"),
                Name: document.getElementById("Name").getAttribute("value"),
            }
        })


        console.log(`BmlID  : ` + this.state.BmlID)
        console.log(`DgkImgLink  : ` + this.state.DgkImgLink)

        let ws = new WebSocket('ws://' + envConfig.WebsocketIP + `:` + envConfig.WebsocketPort);
        
        this.setState({ ws: ws })

        ws.onmessage = function (event) {
            console.log(event)
            let parsedEvent = JSON.parse(event.data)

            switch (parsedEvent.name) {
                case 'manualMatchingTablePage change':
                    this.setState({ table: parsedEvent.data.Table })
                    this.setState({ pageNumberList: parsedEvent.data.PageNumberList })
                    this.setState({ BmlSkuName: parsedEvent.data.BmlSkuName })
                    this.setState({ BmlPrice: parsedEvent.data.BmlPrice })
                    this.setState({ BmlImgLink: parsedEvent.data.BmlImgLink })
                    this.setState({ BmlBrand: parsedEvent.data.BmlBrand })
                    this.setState({ BmlSKULink: parsedEvent.data.BmlSKULink })
                    break
                case 'dgkFilterCategory1 get':
                    this.setState({ filterCategory1OptionList: parsedEvent.data.OptionList })
                    break
                case 'dgkFilterCategory2 get':
                    this.setState({ filterCategory2OptionList: parsedEvent.data.OptionList })
                    break
                case 'FrequencyOptionList get':
                    this.setState({ FrequencyOptionList: parsedEvent.data.FrequencyOptionList })
                    break
                case 'update status':
                    this.setState({ UpdateStatus: parsedEvent.data })
                    break
                case 'matching status':
                    this.setState({ DgkImgLink: parsedEvent.data })
                    break
                case 'unmatch status':
                    this.setState({ DgkImgLink: parsedEvent.data })
                    break
                case 'productUrlValidation status':
                    this.setState({ UrlValidationStatus: parsedEvent.data })
                    console.log(this.state.UrlValidationStatus)
                    break
                default: console.log(`Event was not recognized`)

            }

        }.bind(this)

        ws.onopen = function () {
            this.setState({ DgkImgLink: document.getElementById("DgkImgLink").getAttribute("value") })
            this.requestManualMatchingPage(this.state.BmlID, 1, this.state.category1, this.state.category2, '', this.state.minPrice, this.state.maxPrice);
            this.requestDgkCategoryFilter1OptionList();
            this.requestDgkCategoryFilter2OptionList(this.state.category1);
            this.requestFrequencyOptionList();
        }.bind(this)
    }

    requestDgkCategoryFilter1OptionList() {
        const { ws } = this.state;
        let dgkCategoryFilter1OptionListEvent = {
            name: 'dgkCategoryFilter1OptionList request',
            data: {
                isRequested: true
            }
        }
        ws.send(JSON.stringify(dgkCategoryFilter1OptionListEvent))
    }
    requestDgkCategoryFilter2OptionList(category1) {
        const { ws } = this.state;
        let dgkCategoryFilter2OptionListEvent = {
            name: 'dgkCategoryFilter2OptionList request',
            data: {
                isRequested: true,
                category1: category1
            }
        }
        ws.send(JSON.stringify(dgkCategoryFilter2OptionListEvent))

    }
    requestFrequencyOptionList() {
        const { ws } = this.state;
        let FrequencyOptionListEvent = {
            name: 'FrequencyOptionList request',
            data: {
                isRequested: true
            }
        }
        ws.send(JSON.stringify(FrequencyOptionListEvent))
    }

    requestManualMatchingPage(bmlID, pageNumber, category1, category2, skuName, minPrice ,maxPrice) {
        const { ws } = this.state;
        let manualMatchingtablePageRequestEvent = {
            name: 'ManualMatchingTablePage request',
            data: {
                bmlID: bmlID,
                pageNumber: pageNumber,
                category1: category1,
                category2: category2,
                skuName: skuName,
                minPrice: minPrice,
                maxPrice: maxPrice,
            }
        }
        ws.send(JSON.stringify(manualMatchingtablePageRequestEvent))

    }
    requestApplyManualMatching(bmlID, dgkFK, DgkImgLink) {
        const { ws } = this.state;
        let applyManualMatchingEvent = {
            name: 'ApplyManualMatching request',
            data: {
                IDBmlCatalogConfig: bmlID,
                FKDgkCatalogConfig: dgkFK,
                DgkImgLink: DgkImgLink,
                Email: this.state.User.Email,
                Name: this.state.User.Name,
            }
        }

        ws.send(JSON.stringify(applyManualMatchingEvent))
    }
    requestApplyUnmatch(bmlID) {
        const { ws } = this.state;
        let applyUnmatchEvent = {
            name: 'ApplyUnmatch request',
            data: {
                IDBmlCatalogConfig: bmlID
            }
        }

        ws.send(JSON.stringify(applyUnmatchEvent))
        console.log("request check")
    }
    requestUpdateFrequency(bmlID, checkCycle) {
        const { ws } = this.state;
        let updateFrequencyEvent = {
            name: 'UpdateFrequency request',
            data: {
                bmlID: bmlID,
                checkCycle: checkCycle,
            }
        }
        console.log(`UpdateFreq  : ` + checkCycle)
        ws.send(JSON.stringify(updateFrequencyEvent))
    }

   
    requestApplyProductUrl(ProductUrl,BmlID) {
        const { ws } = this.state;
        console.log(ProductUrl)

        let ApplyProductUrlEvent = {
            name: 'ApplyProductUrl request',
            data: {
                ProductUrl: ProductUrl,
                BmlID: BmlID,

            }
        }

        ws.send(JSON.stringify(ApplyProductUrlEvent))

    }


    // should set new state for page number and fetch state of catgeory and skuName then call requestTablePage
    changePageNumber(pageNumber) {
        const { category1, category2, skuName , minPrice, maxPrice} = this.state;
        this.setState({ pageNumber: pageNumber })
        this.requestManualMatchingPage(this.state.BmlID, pageNumber, category1, category2, skuName, minPrice ,maxPrice )

    }
    ApplyManualMatching(dgkFK, DgkImgLink) {
        const { BmlID } = this.state;
        this.setState({
            dgkFK: dgkFK,
            DgkImgLink: DgkImgLink
        })
        this.requestApplyManualMatching(BmlID, dgkFK, DgkImgLink)
    }
    ApplyUnmatch() {
        const { BmlID } = this.state;
        this.requestApplyUnmatch(BmlID)
        console.log("function check")
    }
    UpdateFrequency(checkCycle) {

        console.log(`BmlID  : ` + this.state.BmlID)
        console.log(`DgkImgLink  : ` + this.state.DgkImgLink)

        const { BmlID } = this.state;
        this.setState({ checkCycle: checkCycle })
        console.log(`UpdateFreq in func  : ` + checkCycle)
        this.requestUpdateFrequency(BmlID, checkCycle)
    }

    // choose category for specified SKU - if SKU = all or blank then show all SKU for the category
    // goes back to pageNumber = 1
    chooseCategory1(category1) {
        const { skuName,minPrice,maxPrice } = this.state;
        this.setState({ category1: category1, category2: 'All' })
        this.requestManualMatchingPage(this.state.BmlID, 1, category1, "All", skuName, minPrice ,maxPrice )
        this.requestDgkCategoryFilter2OptionList(category1);
    }
    chooseCategory2(category2) {
        const { category1, skuName, minPrice ,maxPrice } = this.state;
        this.setState({ category2: category2 })
        this.requestManualMatchingPage(this.state.BmlID, 1, category1, category2, skuName, minPrice ,maxPrice )

    }

    // search for skuName within the category - if category = all or blank then seach among all SKUs
    // goes back to pageNumber = 1
    searchSkuName(skuName) {
        const { category1,category2,minPrice,maxPrice  } = this.state;
        this.setState({ skuName: skuName })
        this.requestManualMatchingPage(this.state.BmlID, 1, category1, category2, skuName,minPrice,maxPrice)

    }
    ApplyProductUrl(productUrl) {
        const { BmlID } = this.state;
        this.requestApplyProductUrl(productUrl,BmlID)
    }
    ApplyPriceFilter(minPrice,maxPrice){
        const { category1,category2,skuName  } = this.state;
console.log(minPrice,maxPrice)
        this.setState({ 
            minPrice: minPrice,
            maxPrice: maxPrice
        })
        this.requestManualMatchingPage(this.state.BmlID, 1, category1, category2, skuName, minPrice ,maxPrice )

    }


    // THIS ONE SHOULD BE SENT BY THE GO SERVER
    /*changeTablePage(tableInput) {
        const { ws } = this.state;
        let tablePageChangeEvent = {
            name: 'tablePage change',
            data: tableInput
        }
        ws.send(JSON.stringify(tablePageChangeEvent))
    }*/


    render() {
        return (
            <div class="container">
                <div className="row">

                    <div className="col-xs-3">
                    </div>
                    <div className="col-xs-9">
                        <ul>
                            <li className='filterOption'>
                                <FilterCategory1 optionList1={this.state.filterCategory1OptionList} chooseCategory1={this.chooseCategory1.bind(this)} /*pageTwo={this.pageTwo.bind(this)}*/ />
                            </li>
                            <li className='filterOption'>
                                <FilterCategory2 optionList2={this.state.filterCategory2OptionList} chooseCategory2={this.chooseCategory2.bind(this)} /*pageTwo={this.pageTwo.bind(this)}*/ />
                            </li>
                            <li className='filterOption'>
                                <SearchSkuName searchSkuName={this.searchSkuName.bind(this)} />
                            </li>
                        </ul>
                    </div>

                </div>
                <div className="row">


                    <div className="col-xs-4">
                        <div calssName='text-center'>
                            <ShowBamiloSKU BmlSkuName={this.state.BmlSkuName} BmlPrice={this.state.BmlPrice} BmlImgLink={this.state.BmlImgLink} BmlBrand={this.state.BmlBrand} BmlSKULink={this.state.BmlSKULink} />
                        </div>
                        <div className='row-gap-sm'></div>
                        <div  className='filterOption'>
                            <UpdateFrequency UpdateStatus={this.state.UpdateStatus} optionList={this.state.FrequencyOptionList} UpdateFrequency={this.UpdateFrequency.bind(this)} />
                        </div>
                        <div className='row-gap-sm2'></div>
                        <div  className='filterOption'>
                            <ProductUrl ValidationStatus={this.state.UrlValidationStatus} AddProducUrlStatus={this.state.AddProducUrlStatus} ApplyProductUrl={this.ApplyProductUrl.bind(this)} />
                        </div>
                    </div>


                    <div className="col-xs-6">

                        <div>
                            <Table table={this.state.table} ApplyManualMatching={this.ApplyManualMatching.bind(this)} UpdateFrequency={this.UpdateFrequency.bind(this)} />
                        </div>
                        <PageNumberWrapper pageNumberList={this.state.pageNumberList} changePageNumber={this.changePageNumber.bind(this)} />
                    </div>
                    <div className='col-xs-2'>

                   
                    <div className='row-gap'>
                    <FliterPriceRange  ApplyPriceFilter={this.ApplyPriceFilter.bind(this)} />

                    </div>
                        <ShowMatchedSKU DgkImgLink={this.state.DgkImgLink} ApplyUnmatch={this.ApplyUnmatch.bind(this)} />

                    </div>

                </div>
            </div>

        )

    }

}
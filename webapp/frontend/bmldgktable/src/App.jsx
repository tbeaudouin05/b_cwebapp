import React, { Component } from 'react';
import Table from './table/table.jsx';
import FilterBiCategory from './filter/filterbicategory/filterbicategory.jsx';
import FilterCategory1 from './filter/filtercategory1/filtercategory1.jsx';
import FilterCategory2 from './filter/filtercategory2/filtercategory2.jsx';
import FilterCategory3 from './filter/filtercategory3/filtercategory3.jsx';
import SearchSkuName from './filter/searchskuname.jsx';
import PageNumberWrapper from './pagination/pagenumberwrapper.jsx';
import SortFilter from './filter/sortfilter/sortfilter.jsx';
import envConfig from './envConfig.json'

/*var data = {
    table: [
        {
            rowValue: {
                bmlIDCatalogConfig: 12345,
                bmlSKUName: "شومیز نخی زنانه",
                bmlImgLink: "https://media.bamilo.com/p/navales-2936-852971-1-product.jpg",
                bmlSKULink: "https://www.bamilo.com/navales-شومیز-نخی-زنانه-179258.html",
                bmlSKUPrice: '300,000',
                dgkScore: '275',
                dgkSKUName: "شومیز نخی زنانه",
                dgkImgLink: "https://media.bamilo.com/p/navales-2936-852971-1-product.jpg",
                dgkSKULink: "https://www.bamilo.com/navales-شومیز-نخی-زنانه-179258.html",
                dgkSKUPrice: '300,000'
            },
            rowKey: 1
        }
    ],
    pageNumberList: [1, 2, 3]
}

var filterCategory2OptionList = [{ optionValue: "category1", optionText: "category1" }, { optionValue: "category2", optionText: "category2" }]*/

export default class App extends Component {

    constructor(props) {
        super(props);
        this.state = {
            table: [],// data.table,
            filterBiCategoryOptionList: [],
            filterCategory1OptionList: [],
            filterCategory2OptionList: [], // filterCategory2OptionList,
            filterCategory3OptionList: [],
            sortingFilter: [],
            ws: null,
            goodMatchCount: 0,
            pageNumber: 1,
            pageNumberList: [], // data.pageNumberList,
            biCategory: '',
            category1: '',
            category2: '',
            category3: '',
            skuName: '',
            sorting: '-count_of_soi',
            //numberOfRow: 10,
            User: {
                Email: null,
                Name: null,
            }
        }
    }


    componentDidMount() {

        this.setState({
            User: {
                Email: document.getElementById("Email").getAttribute("value"),
                Name: document.getElementById("Name").getAttribute("value"),
            }
        })

        let ws = new WebSocket('ws://' + envConfig.WebsocketIP + `:` + envConfig.WebsocketPort);

        this.setState({ ws: ws })

        ws.onmessage = function (event) {
            console.log(event)
            let parsedEvent = JSON.parse(event.data)

            switch (parsedEvent.name) {
                case 'tablePage change':
                    this.setState({ table: parsedEvent.data.Table })
                    this.setState({ pageNumberList: parsedEvent.data.PageNumberList })
                    break
                case 'filterBiCategory get':
                    this.setState({ filterBiCategoryOptionList: parsedEvent.data.OptionList })
                    break
                case 'filterCategory1 get':
                    this.setState({ filterCategory1OptionList: parsedEvent.data.OptionList })
                    break
                case 'filterCategory2 get':
                    this.setState({ filterCategory2OptionList: parsedEvent.data.OptionList })
                    break
                case 'filterCategory3 get':
                    this.setState({ filterCategory3OptionList: parsedEvent.data.OptionList })
                    break
                case 'sortingFilter get':
                    this.setState({ sortingFilter: parsedEvent.data })
                    break
                /*case 'goodMatchCount get':
                    this.setState({ goodMatchCount: parsedEvent.data })
                    break*/
                default: console.log(`Event was not recognized`)

            }

        }.bind(this)

        ws.onopen = function () {
            this.requestTablePage(1, '', '', '', '', '', this.state.sorting);
            //this.requestGoodMatchCount('', '', this.state.category3)
            this.requestBiCategoryFilterOptionList();
            this.requestCategoryFilter1OptionList();
            this.requestCategoryFilter2OptionList();
            this.requestCategoryFilter3OptionList();
            this.requestSortingFilterOptionList();
        }.bind(this)
    }

    /*requestGoodMatchCount(category1, category2, category3) {
        const { ws } = this.state;
        let goodMatchCountRequestEvent = {
            name: 'goodMatchCount request',
            data: {
                category1: category1,
                category2: category2,
                category3: category3
            }
        }
        ws.send(JSON.stringify(goodMatchCountRequestEvent))
    }*/
    requestBiCategoryFilterOptionList() {
        const { ws } = this.state;
        let categoryFilter1OptionListEvent = {
            name: 'biCategoryFilterOptionList request',
            data: {
                isRequested: true,
            }
        }
        ws.send(JSON.stringify(categoryFilter1OptionListEvent))
    }

    requestCategoryFilter1OptionList() {
        const { ws } = this.state;
        let categoryFilter1OptionListEvent = {
            name: 'categoryFilter1OptionList request',
            data: {
                isRequested: true,
                /*category2: category2,
                category3: category3*/
            }
        }
        ws.send(JSON.stringify(categoryFilter1OptionListEvent))
    }

    requestCategoryFilter2OptionList() {
        const { ws } = this.state;
        let categoryFilter2OptionListEvent = {
            name: 'categoryFilter2OptionList request',
            data: {
                isRequested: true,
                /*category1: category1,
                 category3: category3*/
            }
        }
        ws.send(JSON.stringify(categoryFilter2OptionListEvent))
    }

    requestCategoryFilter3OptionList() {
        const { ws } = this.state;
        let categoryFilter3OptionListEvent = {
            name: 'categoryFilter3OptionList request',
            data: {
                isRequested: true,
                /*category1: category1,
                 category2: category2*/
            }
        }
        ws.send(JSON.stringify(categoryFilter3OptionListEvent))
    }

    requestSortingFilterOptionList() {
        const { ws } = this.state;
        let sortingFilterOptionListEvent = {
            name: 'requestSortingFilterOptionList request',
            data: {
                isRequested: true
            }
        }
        ws.send(JSON.stringify(sortingFilterOptionListEvent))
    }

    requestTablePage(pageNumber, biCategory, category1, category2, category3, skuName, sorting) {
        const { ws } = this.state;
        console.log(`skuName: ` + skuName)
        let tablePageRequestEvent = {
            name: 'tablePage request',
            data: {
                pageNumber: pageNumber,
                biCategory: biCategory,
                category1: category1,
                category2: category2,
                category3: category3,
                skuName: skuName,
                sorting: sorting
            }
        }
        ws.send(JSON.stringify(tablePageRequestEvent))
    }

    /* requestCsvOutput(numberOfRow, category1, category2, category3, skuName) {
         const { ws } = this.state;
         let csvOutputRequestEvent = {
             name: 'csvOutput request',
             data: {
                 numberOfRow: numberOfRow,
                 category1: category1,
                 category2: category2,
                 category3: category3,
                 skuName: skuName
             }
         }
         ws.send(JSON.stringify(csvOutputRequestEvent))
     }*/

    requestSetGoodMatch(isGoodMatched, id) {
        const { ws } = this.state;
        let setSetGoodMatchRequestEvent = {
            name: 'SetGoodMatch request',
            data: {
                IsGoodMatched: isGoodMatched,
                ID: id,
                Email: this.state.User.Email,
                Name: this.state.User.Name,
            }
        }
        ws.send(JSON.stringify(setSetGoodMatchRequestEvent))
    }


    // should set new state for page number and fetch state of catgeory and skuName then call requestTablePage
    changePageNumber(pageNumber) {
        const { biCategory, category1, category2, category3, skuName, sorting } = this.state;
        this.setState({ pageNumber: pageNumber })
        this.requestTablePage(pageNumber, biCategory, category1, category2, category3, skuName, sorting)

    }

    chooseBiCategory(biCategory) {
        const { skuName, sorting, category1, category2, category3 } = this.state;
        this.setState({ biCategory: biCategory })
        this.requestTablePage(1, biCategory, category1, category2, category3, skuName, sorting)
    }

    // goes back to pageNumber = 1
    chooseCategory1(category1) {
        const { skuName, sorting, biCategory, category2, category3 } = this.state;
        this.setState({ category1: category1 })
        this.requestTablePage(1, biCategory, category1, category2, category3, skuName, sorting)
    }

    chooseCategory2(category2) {
        const { skuName, sorting, biCategory, category1, category3 } = this.state;
        this.setState({ category2: category2 })
        this.requestTablePage(1, biCategory, category1, category2, category3, skuName, sorting)
    }

    chooseCategory3(category3) {
        const { skuName, sorting, biCategory, category1, category2 } = this.state;
        this.setState({ category3: category3 })
        this.requestTablePage(1, biCategory, category1, category2, category3, skuName, sorting)
    }

    chooseSorting(sorting) {
        const { biCategory, category1, category2, category3, skuName, pageNumber } = this.state;
        this.setState({ sorting: sorting })
        this.requestTablePage(pageNumber, biCategory, category1, category2, category3, skuName, sorting)

    }

    // search for skuName within the category2 - if category2 = all or blank then seach among all SKUs
    // goes back to pageNumber = 1
    searchSkuName(skuName) {
        const { biCategory, category1, category2, category3, sorting } = this.state;
        this.setState({ skuName: skuName })
        this.requestTablePage(1, biCategory, category1, category2, category3, skuName, sorting)

    }

    SetGoodMatch(isGoodMatched, rowKey) {
        this.state.table[parseInt(rowKey) - 1].RowValue.GoodMatch = isGoodMatched
        this.requestSetGoodMatch(isGoodMatched, this.state.table[parseInt(rowKey) - 1].RowValue.BmlIDCatalogConfig)
    }

    seeSKUHistory(sentToServer) {
        const { ws } = this.state;
        let skuHistoryRequestEvent = {
            name: 'skuHistory request',
            data: sentToServer
        }
        ws.send(JSON.stringify(skuHistoryRequestEvent))
    }

    /*<div className='col-xs-6'>
        <p className='good_match_count'>Good Match Count: <strong>{this.state.goodMatchCount}</strong></p>
    </div>*/
    render() {
        return (
            <div>
                <div className='row'>
                    <div className='col-xs-1'>
                        <p className='category_label'>Categories:</p>
                    </div>
                    <div className='col-xs-2'>
                        <li className='filterOption'>
                            <FilterBiCategory optionList={this.state.filterBiCategoryOptionList} chooseBiCategory={this.chooseBiCategory.bind(this)} />
                        </li>
                    </div>
                    <div className='col-xs-2'>
                        <li className='filterOption'>
                            <FilterCategory1 optionList={this.state.filterCategory1OptionList} chooseCategory1={this.chooseCategory1.bind(this)} />
                        </li>
                    </div>
                    <div className='col-xs-2'>
                        <li className='filterOption'>
                            <FilterCategory2 optionList={this.state.filterCategory2OptionList} chooseCategory2={this.chooseCategory2.bind(this)} />
                        </li>
                    </div>
                    <div className='col-xs-2'>
                        <li className='filterOption'>
                            <FilterCategory3 optionList={this.state.filterCategory3OptionList} chooseCategory3={this.chooseCategory3.bind(this)} />
                        </li>
                    </div>
                </div>
                <div className='row'>
                    <div className='col-xs-3'>
                        <li className='filterOption'>
                            <SearchSkuName searchSkuName={this.searchSkuName.bind(this)} />
                        </li>
                    </div>
                    <div className='col-xs-3'>
                        <li className='filterOption'>
                            <SortFilter optionList={this.state.sortingFilter} chooseSorting={this.chooseSorting.bind(this)} /*pageTwo={this.pageTwo.bind(this)}*/ />
                        </li>
                    </div>
                </div>

                <div>
                    <Table table={this.state.table} seeSKUHistory={this.seeSKUHistory.bind(this)} SetGoodMatch={this.SetGoodMatch.bind(this)} />
                </div>
                <PageNumberWrapper pageNumberList={this.state.pageNumberList} changePageNumber={this.changePageNumber.bind(this)} />
            </div>
        )

    }

}

/*  CSV OUTPUT BUTTON
                  <div className='col-xs-1'>
                        <div className='row-gap'></div>
                        <form className="row_button" action="/bmldgktablecsvoutput" method="post">
                            <input type="hidden" name="NumberOfRow" defaultValue={this.state.numberOfRow} ref={this.numberOfRow} />
                            <input type="hidden" name="PageNumber" defaultValue={this.state.pageNumber} ref={this.pageNumber} />
                            <input type="hidden" name="Category1" defaultValue={this.state.category1} ref={this.category1} />
                            <input type="hidden" name="Category2" defaultValue={this.state.category2} ref={this.category2} />
                            <input type="hidden" name="Category3" defaultValue={this.state.category3} ref={this.category3} />
                            <input type="hidden" name="SearchedBy" defaultValue={this.state.skuName} ref={this.skuName} />
                            <input type="hidden" name="SortedBy" defaultValue={this.state.sorting} ref={this.sorting} />
                            <button type="submit" value="Submit" className="btn btn-secondary " >csv output</button>
                        </form>

                    </div>*/
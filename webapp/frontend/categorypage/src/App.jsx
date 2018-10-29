import React, { Component } from 'react';
import CategoryAssortmentTable from './cattable/categoryassortmenttable.jsx'
import FilterCategory1 from './filter/filtercategory/filtercategory1.jsx'
import envConfig from './envConfig.json'

export default class App extends Component {
    constructor(props) {
        super(props);
        this.state = {

            UrlValidationStatus: "",
            DgkCategoryAssortmentList: [],
            BmlCategoryAssortmentList: [],

            BmlFilterCategory1OptionList: [],
            DgkFilterCategory1OptionList: [],
            BmlCategory: "",
            DgkCategory: "",
            dgkConfigCount: 0,
            bmlLogoImage: "https://upload.wikimedia.org/wikipedia/commons/0/0a/Bamilo_Logo2.png",
            dgkLogoImage: "https://upload.wikimedia.org/wikipedia/fa/2/29/Digikala_logo.png",
        }
    }

    componentDidMount() {

        let ws = new WebSocket('ws://' + envConfig.WebsocketIP + `:` + envConfig.WebsocketPort);
        
        this.setState({ ws: ws })

        ws.onmessage = function (event) {
            console.log(event)
            let parsedEvent = JSON.parse(event.data)

            switch (parsedEvent.name) {
                case 'sellerUrlValidation status':
                    this.setState({ UrlValidationStatus: parsedEvent.data })
                    console.log(this.state.UrlValidationStatus)
                    break
                case 'sellerList get':
                    this.setState({ SellerInfo: parsedEvent.data })
                    break
                case 'dgkCategoryAssortment get':
                    this.setState({ DgkCategoryAssortmentList: parsedEvent.data })
                    break
                case 'bmlCategoryAssortment get':
                    this.setState({ BmlCategoryAssortmentList: parsedEvent.data })
                    break
                case 'dgkFilterCategory1 get':
                    this.setState({ DgkFilterCategory1OptionList: parsedEvent.data.OptionList })
                    break
                case 'bmlFilterCategory1 get':
                    this.setState({ BmlFilterCategory1OptionList: parsedEvent.data.OptionList })
                    break

                default: console.log(`Event was not recognized`)

            }

        }.bind(this)

        ws.onopen = function () {
            this.requestDgkCategoryAssortmentList("");
            this.requestDgkCategoryFilter1OptionList();
            this.requestBmlCategoryAssortmentList("");
            this.requestBmlCategoryFilter1OptionList();

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
    requestDgkCategoryAssortmentList(category) {
        const { ws } = this.state;
        let requestDgkCategoryAssortmentListEvent = {
            name: 'dgkCategoryAssortment request',
            data: {
                isRequested: true,
                DgkCategory: category
            }
        }
        console.log("filter cat:", category)
        // console.log("filter cat2:", this.state.DgkCategory)
        ws.send(JSON.stringify(requestDgkCategoryAssortmentListEvent))
    }
    requestBmlCategoryFilter1OptionList() {
        const { ws } = this.state;
        let bmlCategoryFilter1OptionListEvent = {
            name: 'bmlCategoryFilter1OptionList request',
            data: {
                isRequested: true
            }
        }
        ws.send(JSON.stringify(bmlCategoryFilter1OptionListEvent))
    }
    requestBmlCategoryAssortmentList(category) {
        const { ws } = this.state;
        let requestBmlCategoryAssortmentListEvent = {
            name: 'bmlCategoryAssortment request',
            data: {
                isRequested: true,
                BmlCategory: category
            }
        }
        console.log("filter cat:", category)
        ws.send(JSON.stringify(requestBmlCategoryAssortmentListEvent))
    }



    ChooseCategory1Bml(category) {
        this.setState({ BmlCategory: category })
        this.requestBmlCategoryAssortmentList(category)

    }
    ChooseCategory1Dgk(category) {
        this.setState({ DgkCategory: category })
        this.requestDgkCategoryAssortmentList(category)

    }



    render() {
        return (
            <div class="container">

                <div className="row">
                    <div className="col-xs-5">
                        <div className="slim_img__wrap">
                            <p >  </p>
                            <img src={this.state.bmlLogoImage} height="30" width="150" />
                        </div>
                    </div>
                    <div className="col-xs-1"></div>


                    <div className="col-xs-6">
                        <div className="slim_img__wrap">
                            <img src={this.state.dgkLogoImage} height="50" width="150" />
                        </div>
                    </div>
                </div>

                <div className="row">


                    <div className="col-xs-5">
                    <div className="col-xs-6">
                    <div className="filterOption">
                        < FilterCategory1 filterCategory1OptionList={this.state.BmlFilterCategory1OptionList} ChooseCategory1={this.ChooseCategory1Bml.bind(this)} />
                    </div>
                    </div>
                    </div>
                    <div className="col-xs-1">
                        <div>
                        </div>

                    </div>
                    <div className="col-xs-5">
                            <div className="col-xs-6">
                            <div className="filterOption">
                                < FilterCategory1 filterCategory1OptionList={this.state.DgkFilterCategory1OptionList} ChooseCategory1={this.ChooseCategory1Dgk.bind(this)} />
                                </div>
                            </div>
                            </div>





                </div>
                <p > </p>
                <div className="row">

                    <div className="col-xs-5">
                        <div>
                            <CategoryAssortmentTable CategoryAssortmentList={this.state.BmlCategoryAssortmentList} />
                        </div>

                    </div>
                    <div className="col-xs-1">
                        <div>
                        </div>

                    </div>
                    <div className="col-xs-5">
                        <div>
                            <CategoryAssortmentTable CategoryAssortmentList={this.state.DgkCategoryAssortmentList} />
                        </div>

                    </div>
                </div>
            </div>

        )

    }

}
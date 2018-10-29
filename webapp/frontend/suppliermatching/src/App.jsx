import React, { Component } from 'react';
import SellerUrl from './addurl/addsellerurl.jsx'
import BmlSupplierTable from './table/bmlsuppliertable.jsx'
import DgkSupplierTable from './table/dgksuppliertable.jsx'
import SearchDgkSellerName from './filter/searchdgksellername.jsx'
import SearchBmlSellerName from './filter/searchbmlsellername.jsx'
import BmlDgkSupplierTable from './table/bmldgkmatchedsuppliertable.jsx'
import envConfig from './envConfig.json'
export default class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            DgkSellerInfo: [],
            BmlSellerInfo: [],

            DgkSearchedBy: "",
            BmlSearchedBy: "",
            UrlValidationStatus: "",
            CategoryAssortmentList: [],
            filterCategory1OptionList: [],
            category: "",
            dgkConfigCount: 0,
            bmlLogoImage: "https://upload.wikimedia.org/wikipedia/commons/0/0a/Bamilo_Logo2.png",
            dgkLogoImage: "https://upload.wikimedia.org/wikipedia/fa/2/29/Digikala_logo.png",
            SelectedDgkSellerID: 0,
            SelectedDgkSellerName: "",
            SelectedBmlSellerID: 0,
            SelectedBmlSellerName: "",
            MatchStatus: "",
            DgkSearchedByMatch: "",
            DgkSearchedByMatch: "", ///I put ability to search just by one in backend 
            bothSellerInfo: [],
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
                case 'sellerUrlValidation status':
                    this.setState({ UrlValidationStatus: parsedEvent.data })
                    console.log(this.state.UrlValidationStatus)
                    break
                case 'dgkSellerList get':
                    this.setState({ DgkSellerInfo: parsedEvent.data })
                    break
                case 'bmlSellerList get':
                    this.setState({ BmlSellerInfo: parsedEvent.data })
                    break
                case 'matchedSellerList get':
                    this.setState({ bothSellerInfo: parsedEvent.data })
                    break
                case 'MatchSupplier status':
                    this.setState({ MatchStatus: parsedEvent.data })
                    break

                default: console.log(`Event was not recognized`)

            }

        }.bind(this)

        ws.onopen = function () {
            this.requestMatchedSellerList("", "");
            this.requestDgkSellerList(this.state.DgkSearchedBy);
            this.requestBmlSellerList(this.state.DgkSearchedBy);

        }.bind(this)
    }


    requestDgkSellerList(DgkSearchedBy) {
        const { ws } = this.state;
        let requestDgkSellerListEvent = {
            name: 'DgkSellerList request',
            data: {
                isRequested: true,
                DgkSearchedBy: DgkSearchedBy
            }
        }
        ws.send(JSON.stringify(requestDgkSellerListEvent))
    }

    requestBmlSellerList(BmlSearchedBy) {
        const { ws } = this.state;
        let requestBmlSellerListEvent = {
            name: 'BmlSellerList request',
            data: {
                isRequested: true,
                BmlSearchedBy: BmlSearchedBy
            }
        }
        ws.send(JSON.stringify(requestBmlSellerListEvent))
    }

    requestMatchedSellerList(BmlSearchedByMatch, DgkSearchedByMatch) {
        const { ws } = this.state;
        let requestMatchedSellerListEvent = {
            name: 'MatchedSellerListEvent request',
            data: {
                isRequested: true,
                BmlSearchedByMatch: BmlSearchedByMatch,
                DgkSearchedByMatch: DgkSearchedByMatch
            }
        }
        ws.send(JSON.stringify(requestMatchedSellerListEvent))
    }

    requestApplySellerUrl(SellerUrl) {
        const { ws } = this.state;
        console.log(SellerUrl)

        let ApplySellerUrlEvent = {
            name: 'ApplySellerUrl request',
            data: {
                SellerUrl: SellerUrl,
            }
        }

        ws.send(JSON.stringify(ApplySellerUrlEvent))

    }

    requestMatchSeller(selectedDgkSellerID, selectedDgkSellerName, selectedBmlSellerID, selectedBmlSellerName) {
        const { ws } = this.state;
        console.log("match:", selectedDgkSellerID, selectedBmlSellerID)

        let MatchSellerEvent = {
            name: 'MatchSeller request',
            data: {
                SelectedDgkSellerID: selectedDgkSellerID,
                SelectedDgkSellerName: selectedDgkSellerName,
                SelectedBmlSellerID: selectedBmlSellerID,
                SelectedBmlSellerName: selectedBmlSellerName,
                Email: this.state.User.Email,
                Name: this.state.User.Name,
            }
        }

        ws.send(JSON.stringify(MatchSellerEvent))

    }

    SearchDgkSellerName(DgksearchedBy) {
        this.requestDgkSellerList(DgksearchedBy)
    }
    SearchBmlSellerName(BmlsearchedBy) {
        this.requestBmlSellerList(BmlsearchedBy)
    }

    ApplySellerUrl(sellerUrl) {
        this.requestApplySellerUrl(sellerUrl)
    }

    SelectDgkSeller(id, selectedDgkSellerName) {
        this.setState({
            SelectedDgkSellerID: id,
            SelectedDgkSellerName: selectedDgkSellerName,
        })
        console.log("SelectedDgkSellerName", selectedDgkSellerName, this.state.SelectedDgkSellerName)
        console.log("id", id, this.state.SelectedDgkSellerID)
    }

    SelectBmlSeller(id, selectedBmlSellerName) {
        this.setState({
            SelectedBmlSellerID: id,
            SelectedBmlSellerName: selectedBmlSellerName,
        })
        console.log("SelectedBmlSellerName", selectedBmlSellerName, this.state.SelectedBmlSellerName)
        console.log("id", id, this.state.SelectedBmlSellerID)
    }
    onClick(e) {
        e.preventDefault();

        this.requestMatchSeller(this.state.SelectedDgkSellerID, this.state.SelectedDgkSellerName, this.state.SelectedBmlSellerID, this.state.SelectedBmlSellerName);
        this.requestMatchedSellerList(this.state.BmlSearchedByMatch, this.state.DgkSearchedByMatch);

    }


    render() {
        return (
            <div class="container">

                <div className="row">

                    <div className="col-xs-3">

                        <div className="slim_img__wrap">
                            <p >  </p>
                            <img src={this.state.bmlLogoImage} alt="bamilo" height="30" width="150" />
                        </div>


                    </div>
                    <div className="col-xs-1"></div>


                    <div className="col-xs-3">
                        <div className="slim_img__wrap">
                            <img src={this.state.dgkLogoImage} alt="digikala" height="50" width="150" />

                        </div>

                    </div>
                    <div className="col-xs-2"></div>

                    <div className="col-xs-3"><p >  </p> {this.state.MatchStatus} </div>

                </div>

                <div className="row">



                    <div className="col-xs-3">

                        <SearchBmlSellerName SearchBmlSellerName={this.SearchBmlSellerName.bind(this)} />
                    </div>
                    <div className="col-xs-1"><div>
                    </div>

                    </div>

                    <div className="col-xs-3">

                        <div className="row">
                            <div className="col-xs-5">
                                < SearchDgkSellerName SearchDgkSellerName={this.SearchDgkSellerName.bind(this)} />

                            </div>

                        </div>


                    </div>
                    <div className="col-xs-3"></div>
                    <div className="col-xs-2">
                        <div className="row">
                            <form className="row_button" onClick={this.onClick.bind(this)} target="_blank">
                                <button type="submit" className="btn btn-primary btn-lg medium-size">Match</button>
                            </form>
                        </div>
                    </div>

                </div>
                <p > </p>
                <div className="row">

                    <div className="col-xs-3">
                        <div>
                            <BmlSupplierTable SellerInfo={this.state.BmlSellerInfo} SelectBmlSeller={this.SelectBmlSeller.bind(this)} />

                        </div>

                    </div>
                    <div className="col-xs-1">
                        <div>
                        </div>

                    </div>
                    <div className="col-xs-3">
                        <DgkSupplierTable SellerInfo={this.state.DgkSellerInfo} SelectDgkSeller={this.SelectDgkSeller.bind(this)} />
                        <SellerUrl ValidationStatus={this.state.UrlValidationStatus} ApplySellerUrl={this.ApplySellerUrl.bind(this)} />

                    </div>
                    <div className="col-xs-2"></div>
                    <div className="col-xs-3">
                        <strong  >Matched Suppliers: </strong>
                        <BmlDgkSupplierTable SellerInfo={this.state.bothSellerInfo} />
                    </div>

                </div>
            </div>

        )

    }

}
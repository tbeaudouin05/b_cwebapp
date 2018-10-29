import React, { Component } from 'react';
import BmlDgkChart from './chart/chart.jsx';
import envConfig from './envConfig.json';

/*let chartData = {
    Label: ["2018-07-28 13:40:08.653 +0000 UTC","","","","2018-07-30 03:33:48.768 +0000 UTC"], //['7/28 2018', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '7/28 2018'],
    Value: {
        BmlPrice: [4,5,6,7,8],
        BmlSales: [0,0,0,0,0],
        DgkPrice: [1,2,3,4,5]
    },
    SKUName: {
        BmlSKUName: 'bmlSKUName',
        DgkSKUName: 'dgkSKUName'
    }
}*/

export default class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            Data: {
                // do not change labels & series to label & series! This is part of external package!
                labels: [],
                series: []
            },
            BmlSKUName: '',
            DgkSKUName: '',
            //LegendName: ['Bamilo Sales', 'Bamilo Price', 'Digikala Price'],
            //LegendDiv: null,
            BmlDgkSKUHistoricalDataRequest: {
                BmlIDCatalogConfig: null,
                BmlSKUName: null,
                BmlImgLink: null,
                BmlSupplierName: null,
                BmlBrand: null,
                BmlConfigSnapshot: null,
                BmlSumOfStockQuantity: null,
                BmlSKULink: null,

                DgkIDCatalogConfig: null,
                DgkSKUName: null,
                DgkImgLink: null,
                DgkSupplierName: null,
                DgkBrand: null,
                DgkConfigSnapshot: null,
                DgkStock: null,
                DgkSKULink: null,
            }
        }
    }


    componentDidMount() {

        // this.setState({ LegendDiv: document.getElementById("legend_position") })
        this.setState({ 
            BmlDgkSKUHistoricalDataRequest: {
                BmlIDCatalogConfig: document.getElementById("BmlIDCatalogConfig").getAttribute("value"),
                BmlSKUName: document.getElementById("BmlSKUName").getAttribute("value"),
                BmlImgLink: document.getElementById("BmlImgLink").getAttribute("value") ,
                BmlSupplierName: document.getElementById("BmlSupplierName").getAttribute("value"),
                BmlBrand: document.getElementById("BmlBrand").getAttribute("value"),
                BmlSumOfStockQuantity: document.getElementById("BmlSumOfStockQuantity").getAttribute("value"),
                BmlConfigSnapshot: document.getElementById("BmlConfigSnapshot").getAttribute("value"),
                BmlSKULink: document.getElementById("BmlSKULink").getAttribute("value"),

                DgkIDCatalogConfig: document.getElementById("DgkIDCatalogConfig").getAttribute("value"),
                DgkSKUName: document.getElementById("DgkSKUName").getAttribute("value"),
                DgkImgLink: document.getElementById("DgkImgLink").getAttribute("value"),
                DgkSupplierName: document.getElementById("DgkSupplierName").getAttribute("value"),
                DgkBrand: document.getElementById("DgkBrand").getAttribute("value"),
                DgkConfigSnapshot: document.getElementById("DgkConfigSnapshot").getAttribute("value"),
                DgkStock: document.getElementById("DgkStock").getAttribute("value"),
                DgkSKULink: document.getElementById("DgkSKULink").getAttribute("value"),
            }
        })

        let ws = new WebSocket('ws://' + envConfig.WebsocketIP + `:` + envConfig.WebsocketPort);

        this.setState({ ws: ws })
        ws.onmessage = function (event) {
            console.log(event)
            let parsedEvent = JSON.parse(event.data)

            console.log(`parsedEvent.name: ` + parsedEvent.name)
            console.log(`parsedEvent.data: ` + parsedEvent.data)
            console.log(`parsedEvent.data.Value: ` + parsedEvent.data.Value)
            console.log(`parsedEvent.data.Value.BmlSales: ` + parsedEvent.data.Value.BmlSales)
            console.log(`parsedEvent.data.Value.BmlPrice: ` + parsedEvent.data.Value.BmlPrice)
            console.log(`parsedEvent.data.Value.DgkPrice: ` + parsedEvent.data.Value.DgkPrice)
            console.log(`parsedEvent.data.SKUName.BmlSKUName: ` + parsedEvent.data.SKUName.BmlSKUName)
            console.log(`parsedEvent.data.SKUName.BmlSKUName: ` + parsedEvent.data.SKUName.BmlSKUName)
            console.log(`parsedEvent.data.SKUName.DgkSKUName: ` + parsedEvent.data.SKUName.DgkSKUName)

            switch (parsedEvent.name) {
                case 'BmlDgkSKUHistoricalData send':
                    this.setState({
                        Data: {
                            // do not change labels & series to label & series! This is part of external package!
                            labels: parsedEvent.data.Label,
                            series: [
                                { name: 'Bamilo Price', data: parsedEvent.data.Value.BmlPrice },
                                { name: 'Digikala Price', data: parsedEvent.data.Value.DgkPrice },
                                { name: 'Bamilo Sales', data: parsedEvent.data.Value.BmlSales }
                            ]
                        }
                    });
                    this.setState({ BmlSKUName: parsedEvent.data.SKUName.BmlSKUName });
                    this.setState({ DgkSKUName: parsedEvent.data.SKUName.DgkSKUName });
                    break
                default: console.log(`Event was not recognized`)
            }

        }.bind(this)

        ws.onopen = function () {
            this.setState({ BmlSKUName: document.getElementById("BmlSKUName").getAttribute("value") })
            this.setState({ DgkSKUName: document.getElementById("DgkSKUName").getAttribute("value") })

            this.requestBmlDgkSKUHistoricalData();
        }.bind(this)

    }

    requestBmlDgkSKUHistoricalData() {
    
        const { ws } = this.state;
        let bmlDgkSKUHistoricalDataEvent = {
            name: 'BmlDgkSKUHistoricalData request',
            data: this.state.BmlDgkSKUHistoricalDataRequest
        }  
        ws.send(JSON.stringify(bmlDgkSKUHistoricalDataEvent))

    }

    render() {
        return (
            <div>
                <h4 className='chart_title'>{this.state.BmlSKUName} <strong>VS.</strong> {this.state.DgkSKUName}</h4>
                <div className="row">
                    <div className="col-xs-1">
                        <div id="legend_position">
                            <li className='ct-series-0'>Bamilo Price</li>
                            <li className='ct-series-1'>Digikala Price</li>
                            <li className='ct-series-2'>Bamilo Sales</li>
                        </div>
                    </div>
                    <div className="col-xs-11">
                        <BmlDgkChart Data={this.state.Data} LegendName={this.state.LegendName} LegendDiv={this.state.LegendDiv} />
                    </div>
                </div>
            </div>

        )

    }

}
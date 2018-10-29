import React, { Component } from 'react';
import envConfig from './envConfig.json';



export default class App extends Component {
    constructor(props) {
        super(props);
        this.state = {

            bmlLogoImage: "https://upload.wikimedia.org/wikipedia/commons/0/0a/Bamilo_Logo2.png",
            torobLogoImage: "http://setak.sharif.ir/main/wp-content/uploads/2015/02/%D8%AA%D8%B1%D8%A8.png",
            User: {
                Email: null,
                Name: null,
            },
           
            name: "",
            email: "",
                totalSKU: 0,
                totalSupplier: 0,
                thisWeekSKU: 0,
                thisWeekSupplier: 0,
                lastWeekSKU: 0,
                lastWeekSupplier: 0,
                totalNotCompetitive: 0,
                Message: ""
               

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
                case 'userInfo get':
                    this.setState({ name: parsedEvent.data.name,
                        email: parsedEvent.data.email,
                        totalSKU: parsedEvent.data.totalSKU,
                        totalSupplier: parsedEvent.data.totalSupplier,
                        thisWeekSKU: parsedEvent.data.thisWeekSKU,
                        thisWeekSupplier: parsedEvent.data.thisWeekSupplier,
                        lastWeekSKU: parsedEvent.data.lastWeekSKU,
                        lastWeekSupplier: parsedEvent.data.lastWeekSupplier,
                        Message: parsedEvent.data.Message,
                        totalNotCompetitive: parsedEvent.data.TotalNotCompetitive
                     })
                    break

                default: console.log(`Event was not recognized`)
            }

        }.bind(this)

        ws.onopen = function () {
            this.requestPersonalizedInfo();
        }.bind(this)
    }
    requestPersonalizedInfo() {
        const { ws } = this.state;
        console.log("name: ",this.state.User.Name)
        let requestPersonalizedInfoEvent = {
            name: 'PersonalizedInfo request',
            data: {
                Email: this.state.User.Email,
                Name: this.state.User.Name,
            }
        }
        ws.send(JSON.stringify(requestPersonalizedInfoEvent))
    }

/*
   <h5 > Total matched SKUs: {this.state.totalSKU} </h5>
                    <h5 > This week(last 7 days) matched SKUs: {this.state.thisWeekSKU}</h5>
                    <h5 > Last week matched SKUs: {this.state.lastWeekSKU} </h5>
                    <h5 > Total matched suppliers: {this.state.totalSupplier} </h5>
                    <h5 > This week matched suppliers {this.state.thisWeekSupplier} </h5>
                    <h5 > Last week matched suppliers: {this.state.lastWeekSupplier} </h5>
                    */
    render() {
        return (
            <div class="container">
                <div className="row">
                <div className="col-xs-8">
                    <div className="row">
                    <p >  </p>
                    <div className="space">
                        </div>
                    <h4 style={{ marginLeft: 10 }} > Hi {this.state.name}! </h4>
                    <div className="space">
                        </div>
                        <form className="row_button" action="http://192.168.100.160:8089/d/JYW2Rzxmk/noncompetitiveskus_dashboard?orgId=1" target="_blank">
                                <button type="submit" className="btn btn-primary mb-2">See NonCompetitive Skus</button>
                            </form>
                    <h5 style={{ marginLeft: 5 , width: 310 }} >  {this.state.Message} </h5>
                    <ul class="list-group" style={{width: 310 }}  >
                    
        <li class="list-group-item d-flex justify-content-between align-items-center">Total matched SKUs:
        <span class="badge badge-primary badge-pill">{this.state.totalSKU} </span></li>
        <li class="list-group-item d-flex justify-content-between align-items-center">Total not-competitive matched SKUs:
        <span class="badge badge-primary badge-pill">{this.state.totalNotCompetitive}</span></li>
        <li class="list-group-item d-flex justify-content-between align-items-center">Last 7 days matched SKUs:
        <span class="badge badge-primary badge-pill">{this.state.thisWeekSKU}</span></li>
        <li class="list-group-item d-flex justify-content-between align-items-center">Total matched suppliers:
        <span class="badge badge-primary badge-pill">{this.state.totalSupplier}</span></li>
      </ul>
                 
                   
                    </div>
                   
                    </div>

                    <div className="col-xs-4 ">
                        <div className="less-space">
                        </div>
                        <div className="row">
                            <div className="img__wrap">
                                <p >  </p>
                                <img src={this.state.torobLogoImage} height="95" width="400" />

                            </div>
                        </div>
                        <div className="space">
                        </div>
                        <div className="row ">

                            <form className="row_button " action="/bmldgktable" method="get" target="_blank">

                                <button type="submit" value="Submit" className="btn btn-success btn-lg same-size" >Go to SKU comparison table</button>
                            </form>


                        </div>

                        <div className="row ">
                            <form className="row_button " action="/suppliermatching" method="get" target="_blank">

                                <button type="submit" value="Submit" className="btn btn-success btn-lg same-size" >Go to supplier manual matching</button>
                            </form>
                        </div>
                        <div className="row">

                            <form className="row_button" action="http://192.168.100.160:8089/d/z6p9oMomz/categories-competitiveness?orgId=1" method="get" target="_blank">

                                <button type="submit" value="Submit" className="btn btn-danger btn-lg same-size" >Go to category competitiveness dashboard</button>
                            </form>

                        </div>
                        <div className="row">
                            <form className="row_button" action="http://192.168.100.160:8089/d/CWFa_HTiz/supplier-assortment-count?orgId=1" method="get" target="_blank">

                                <button type="submit" value="Submit" className="btn btn-danger btn-lg same-size" >Go to supplier assortment count dashboard</button>
                            </form>

                        </div>
                        <div className="row">

                            <form className="row_button" action="http://192.168.100.160:8089/d/3cH6D8omk/category-assortment-count?orgId=1" method="get" target="_blank">

                                <button type="submit" value="Submit" className="btn btn-danger btn-lg same-size" >Go to category assortment count dashboard</button>
                            </form>

                        </div>

                        <div className="row">

                            <form className="row_button" action="/adduser" method="get" target="_blank">

                                <button type="submit" value="Submit" className="btn btn-warning btn-lg same-size" >Add user access</button>
                            </form>

                        </div>


                    </div>

                </div>
            </div>


        )

    }

}
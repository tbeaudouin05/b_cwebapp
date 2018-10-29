import React, { Component } from 'react';

export default class BmlDgkTableRow extends Component {

    constructor(props) {
        super(props);
        this.state = {
            GoodMatch: this.props.RowValue.GoodMatch,
        }

        this.GoodMatch = React.createRef();

        this.BmlIDCatalogConfig = React.createRef();
        this.BmlSKUName = React.createRef();
        this.BmlImgLink = React.createRef();
        this.BmlSupplierName = React.createRef();
        this.BmlBrand = React.createRef();
        this.BmlConfigSnapshot = React.createRef();
        this.BmlSumOfStockQuantity = React.createRef();
        this.BmlSKULink = React.createRef();

        this.DgkIDCatalogConfig = React.createRef();
        this.DgkSKUName = React.createRef();
        this.DgkImgLink = React.createRef();
        this.DgkSupplierName = React.createRef();
        this.DgkBrand = React.createRef();
        this.DgkConfigSnapshot = React.createRef();
        this.DgkStock = React.createRef();
        this.DgkSKULink = React.createRef();
        this.handleInputChange = this.handleInputChange.bind(this);
    }
/*
    componentDidMount() {
        this.setState({ GoodMatch: this.props.RowValue.GoodMatch })

        let ws1 = new WebSocket('ws://localhost:4000');
        ws1.onopen = function () {

        }
    }*/


    handleInputChange(event) {

        const target = event.target;

        const rowKey = target.name;

        this.props.SetGoodMatch(!this.props.RowValue.GoodMatch, rowKey)

        this.setState({ GoodMatch: !this.props.RowValue.GoodMatch });


    }
   

    //
    render() {
        return (
            <tr key={this.props.RowKey}>
                <td key='1'>
                    <div className="row">
                        <div className="col-xs-6">
                            <div className="img__wrap">
                                <img src={this.props.RowValue.BmlImgLink} alt={this.props.RowValue.BmlSKUName} className="img-thumbnail" height="150" width="150" />
                                <p className="img__description">
                                    <strong>{this.props.RowValue.BmlSKUName}</strong><br /><br />
                                    Stock: <strong>{this.props.RowValue.BmlSumOfStockQuantity}</strong><br />
                                    Supplier: <strong>{this.props.RowValue.BmlSupplierName}</strong><br />
                                    Brand: <strong>{this.props.RowValue.BmlBrand}</strong><br />
                                    Last Update: <strong>{this.props.RowValue.BmlConfigSnapshot}</strong><br />
                                </p>
                            </div>
                        </div>
                        <div className="col-xs-3">
                            <form className="row_button" action={this.props.RowValue.BmlSKULink} target="_blank">
                                <button type="submit" className="btn btn-primary mb-2">See on Bamilo</button>
                            </form>
                            <form className="row_button" action="/seebmldgkskuhist" method="post" target="_blank">
                                <input type="hidden" name="BmlIDCatalogConfig" defaultValue={this.props.RowValue.BmlIDCatalogConfig} ref={this.BmlIDCatalogConfig} />
                                <input type="hidden" name="BmlSKUName" defaultValue={this.props.RowValue.BmlSKUName} ref={this.BmlSKUName} />
                                <input type="hidden" name="BmlImgLink" defaultValue={this.props.RowValue.BmlImgLink} ref={this.BmlImgLink} />
                                <input type="hidden" name="BmlSupplierName" defaultValue={this.props.RowValue.BmlSupplierName} ref={this.BmlSupplierName} />
                                <input type="hidden" name="BmlBrand" defaultValue={this.props.RowValue.BmlBrand} ref={this.BmlBrand} />
                                <input type="hidden" name="BmlConfigSnapshot" defaultValue={this.props.RowValue.BmlConfigSnapshot} ref={this.BmlConfigSnapshot} />
                                <input type="hidden" name="BmlSumOfStockQuantity" defaultValue={this.props.RowValue.BmlSumOfStockQuantity} ref={this.BmlSumOfStockQuantity} />
                                <input type="hidden" name="BmlSKULink" defaultValue={this.props.RowValue.BmlSKULink} ref={this.BmlSKULink} />

                                <input type="hidden" name="DgkIDCatalogConfig" defaultValue={this.props.RowValue.DgkIDCatalogConfig} ref={this.DgkIDCatalogConfig} />
                                <input type="hidden" name="DgkSKUName" defaultValue={this.props.RowValue.DgkSKUName} ref={this.DgkSKUName} />
                                <input type="hidden" name="DgkImgLink" defaultValue={this.props.RowValue.DgkImgLink} ref={this.DgkImgLink} />
                                <input type="hidden" name="DgkSupplierName" defaultValue={this.props.RowValue.DgkSupplierName} ref={this.DgkSupplierName} />
                                <input type="hidden" name="DgkBrand" defaultValue={this.props.RowValue.DgkBrand} ref={this.DgkBrand} />
                                <input type="hidden" name="DgkConfigSnapshot" defaultValue={this.props.RowValue.DgkConfigSnapshot} ref={this.DgkConfigSnapshot} />
                                <input type="hidden" name="DgkStock" defaultValue={this.props.RowValue.DgkStock} ref={this.DgkStock} />
                                <input type="hidden" name="DgkSKULink" defaultValue={this.props.RowValue.DgkSKULink} ref={this.DgkSKULink} />
                                <button type="submit" value="Submit" className="btn btn-primary mb-2" >See history</button>
                            </form>
                        </div>
                    </div>
                </td>
                <td key='3'><strong><div className="p"> {this.props.RowValue.BmlAvgPrice}</div></strong><strong>{this.props.RowValue.BmlAvgSpecialPrice}</strong>  </td>
                <td key='4'>
                    <div className="row">
                        <div className="col-xs-6">
                            <div className="img__wrap">
                                <img src={this.props.RowValue.DgkImgLink} alt={this.props.RowValue.DgkSKUName} className="img-thumbnail" height="180" width="180" />
                                <p className="img__description">
                                    <strong>{this.props.RowValue.DgkSKUName}</strong> <br /><br />
                                    Stock: <strong>{this.props.RowValue.DgkStock}</strong><br />
                                    Supplier: <strong>{this.props.RowValue.DgkSupplierName}</strong><br />
                                    Brand: <strong>{this.props.RowValue.DgkBrand}</strong><br />
                                    Last Update: <strong>{this.props.RowValue.DgkConfigSnapshot}</strong><br />
                                    Match Score: <strong>{this.props.RowValue.DgkScore}</strong><br />
                                </p>
                            </div>
                        </div>
                        <div className="col-xs-3">
                            <form className="row_button" action={this.props.RowValue.DgkSKULink} target="_blank">
                                <button type="submit" className="btn btn-primary mb-2" >See on Digikala</button>
                            </form>
                            <form className="row_button" action="/bmldgkskumanualmatching" method="post" target="_blank">
                                <input type="hidden" name="BmlIDCatalogConfig" defaultValue={this.props.RowValue.BmlIDCatalogConfig} ref={this.BmlIDCatalogConfig} />
                                <input type="hidden" name="DgkImgLink" defaultValue={this.props.RowValue.DgkImgLink} ref={this.DgkImgLink} />
                                <button type="submit" value="Submit" className="btn btn-primary mb-2" >Manual change</button>
                            </form>


                            <form>

                                
                                    <input
                                        type="checkbox"
                                        name={this.props.RowKey}

                                        defaultValue={this.state.GoodMatch}
                                        checked={this.props.RowValue.GoodMatch}
                                        onChange={this.handleInputChange} />
                                    Good match
                                

                            </form>

                        </div>

                    </div>
                </td>
                <td key='6'><strong><div className="p"> {this.props.RowValue.DgkAvgPrice}</div></strong><strong>{this.props.RowValue.DgkAvgSpecialPrice}</strong></td>

            </tr >
        )
    }
}


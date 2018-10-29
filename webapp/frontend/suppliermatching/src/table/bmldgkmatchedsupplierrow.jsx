import React, { Component } from 'react';

export default class BothSupplierRow extends Component {

    constructor(props) {
        super(props);

        this.input = React.createRef();
        this.handleInputChange = this.handleInputChange.bind(this);

    }



    handleInputChange() {

        this.props.SelectBmlSeller(this.props.ID, this.props.SellerName)

    }


    render() {
        return (

            <tr key={this.props.rowKey}>


                <td key='1'>
                    <div className="slim_row">
                        <div className="col-xs-6">
                            {this.props.BmlSellerName}
                        </div>
                    </div>
                </td>
                <td key='2'>
                    <div class="slim_row">
                        <div className="col-xs-6">
                            {this.props.DgkSellerName}

                        </div>
                    </div>
                </td>



            </tr>
        )
    }
}

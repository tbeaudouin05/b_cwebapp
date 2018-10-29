import React, { Component } from 'react';

export default class SellerRow extends Component {

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
                        <div className="row">
                            <div className="col-xs-6">
                                {this.props.SellerName}
                            </div>

                            <div className="col-xs-1">
                            <div class="custom-control custom-radio">
                                <input type="radio" class="custom-control-input" id="defaultUnchecked" name="bmlRadios" onChange={this.handleInputChange}  /> 

                            </div>
                        </div>
                            </div>
                        </td>
                        
                

            </tr>
        )
    }
}

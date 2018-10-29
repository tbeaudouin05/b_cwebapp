import React, { Component } from 'react';

export default class CategoryAssortmentRowDeep extends Component {

    constructor(props) {
        super(props);

        this.input = React.createRef();
    }

    render() {
        return (

            <tr key={this.props.rowKey}>
                    
                        <td key='1'>
                        <div className="row">
                            <div className="col-xs-6">
                                {this.props.CategoryTwoName}
                            </div>
                            </div>
                        </td>
                        <td key='2'>
                        <div className="row">
                            <div className="col-xs-6">
                                {this.props.ConfigCountCat2}
                            </div>
                            </div>
                        </td>
                

            </tr>
        )
    }
}

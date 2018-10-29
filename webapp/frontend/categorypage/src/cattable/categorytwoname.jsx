import React, { Component } from 'react';

export default class CategoryTwoName extends Component {

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
                           <div>
                                {this.props.Category2Name}
                               
                            </div>
                            </div>
                            </div>
                        
                        </td>

            </tr>
        )
    }
}

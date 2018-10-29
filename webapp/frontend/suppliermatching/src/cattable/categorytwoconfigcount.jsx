import React, { Component } from 'react';

export default class CategoryTwoConfigCount extends Component {

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
                                {this.props.ConfigCountCat2}
                               
                            </div>
                            </div>
                            </div>
                        
                        </td>

            </tr>
        )
    }
}

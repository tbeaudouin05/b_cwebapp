import React, { Component } from 'react';
import CategoryTwoName from './categorytwoname.jsx'
import CategoryTwoConfigCount from './categorytwoconfigcount.jsx'

export default class CategoryAssortmentRow extends Component {

    constructor(props) {
        super(props);

        this.input = React.createRef();
    }

    render() {
        return (

            <tr key={this.props.rowKey}>
                
                    <td key='1'>
                        <div className="row">
                            <div className="col-xs-7">
                            <ul>
                            <strong>{this.props.CategoryOneName}</strong>
                            <ul>  </ul>
                            {this.props.CategoryTwo.map(row => { 
                                return (
                                    <CategoryTwoName Category2Name={row.CategoryTwoName}  RowKey={row.RowKey} />
                                 )
                            })}
                            </ul>
                            </div>
                        </div>
                    </td>
                    <td key='2'>
                        <div className="row">
                            <div className="col-xs-5">
                            <ul>
                            <strong>{this.props.ConfigCountCat1}</strong>
                            <ul>  </ul>
                            {this.props.CategoryTwo.map(row => { 
                                return (
                                    <CategoryTwoConfigCount ConfigCountCat2={row.ConfigCountCat2} RowKey={row.RowKey} />
                                 )
                            })}
                            </ul>


                            </div>
                        </div>
                    </td>
                   
            </tr>
        )
    }
}

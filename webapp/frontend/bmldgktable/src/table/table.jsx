import React, { Component } from 'react';
import BmlDgkTableRow from './row.jsx'


export default class BmlDgkTable extends Component {

    render() {
        return (
            <table className="table">
                <thead>
                <tr>
                    <th scope="col">Bamilo SKU</th>
                    <th scope="col">Bamilo Price</th>
                    <th scope="col">Digikala SKU</th>
                    <th scope="col">Digikala Price</th>
                </tr>
                </thead>
                <tbody>
                {this.props.table.map(row => {
                    return (
                        <BmlDgkTableRow RowValue={row.RowValue} RowKey={row.RowKey} seeSKUHistory={this.props.seeSKUHistory} SetGoodMatch={this.props.SetGoodMatch}/>
                    )
                })}
                </tbody>

            </table>
        )
    }
}
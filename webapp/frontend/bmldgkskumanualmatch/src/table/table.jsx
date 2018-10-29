import React, { Component } from 'react';
import BmlDgkTableRow from './row.jsx'


export default class BmlDgkTable extends Component {

    render() {
        return (
            <table className="table">
                <thead>
                <tr>
                    <th scope="col">Digikala SKU</th>
                    <th scope="col">Digikala Price</th>
                </tr>
                </thead>
                <tbody>
                {this.props.table.map(row => {
                    return (
                        <BmlDgkTableRow DgkRowValue={row.DgkRowValue} rowKey={row.DgkRowValue} ApplyManualMatching={this.props.ApplyManualMatching} UpdateFrequency={this.props.UpdateFrequency} />
                    )
                })}
                </tbody>

            </table>
        )
    }
}
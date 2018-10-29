import React, { Component } from 'react';
import DgkSellerRow from './dgksupplierrow.jsx'


export default class DgkSupplierTable extends Component {

    render() {
        return (
            <table className="table">
                <thead>
                <tr>
                    <th scope="col">Supplier name</th>
                   
                </tr>
                </thead>
                <tbody>
                {this.props.SellerInfo.map(row => {
                    return (
                        <DgkSellerRow SellerName={row.SellerName} ConfigCount={row.ConfigCount}  ID={row.ID} SelectDgkSeller={this.props.SelectDgkSeller}/>
                    )
                })}
                </tbody>

            </table>
        )
    }
}
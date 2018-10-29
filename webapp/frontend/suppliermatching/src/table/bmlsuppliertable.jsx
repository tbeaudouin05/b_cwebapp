import React, { Component } from 'react';
import SellerRow from './bmlsupplierrow.jsx'


export default class BmlSupplierTable extends Component {

    render() {
        return (
            <div className='flex-vertical'>
                <table  className='table' responsive hover>
                    <thead>
                        <tr>
                            <th scope="col">Supplier name</th>

                        </tr>
                    </thead>
                    <tbody>
                        {this.props.SellerInfo.map(row => {
                            return (
                                <SellerRow SellerName={row.SellerName} ConfigCount={row.ConfigCount} ID={row.ID} SelectBmlSeller={this.props.SelectBmlSeller} />
                            )
                        })}
                    </tbody>

                </table>
                </div>
        )
    }
}
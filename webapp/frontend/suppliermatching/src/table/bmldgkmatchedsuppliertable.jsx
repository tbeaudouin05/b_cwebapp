import React, { Component } from 'react';
import BothSupplierRow from './bmldgkmatchedsupplierrow.jsx'


export default class BmlDgkSupplierTable extends Component {

    render() {
        return (
            <div className='flex-vertical'>
                <table  className='table table-bordered' responsive hover>
                    <thead>
                        <tr>
                            <th scope="col" style={{color:'Orange'}}>Bamilo</th>
                            <th scope="col" style={{color:'Red'}} >Digikala</th>

                        </tr>
                    </thead>
                    <tbody>
                        {this.props.SellerInfo.map(row => {
                            return (
                                <BothSupplierRow BmlSellerName={row.BmlSellerName} DgkSellerName={row.DgkSellerName} />
                            )
                        })}
                    </tbody>

                </table>
                </div>
        )
    }
}


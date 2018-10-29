import React, { Component } from 'react';
import CategoryAssortmentRow from './categoryassortmentrow.jsx'


export default class CategoryAssortmentTable extends Component {

    render() {
        return (
            <table className="table">
                <thead>
                <tr>
                    <th scope="col">Category</th>
                    <th scope="col">Config Count</th>
                </tr>
                </thead>
                <tbody>
                {this.props.CategoryAssortmentList.map(row => {
                    return (
                        <CategoryAssortmentRow CategoryOneName={row.CategoryOneName} ConfigCountCat1={row.ConfigCountCat1} CategoryTwo={row.CategoryTwo} RowKey={row.RowKey} />
                    )
                })}
                </tbody>

            </table>
        )
    }
}
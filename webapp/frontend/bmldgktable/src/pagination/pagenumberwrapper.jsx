import React, { Component } from 'react';
import PageNumber from './pagenumber.jsx'


export default class PageNumberWrapper extends Component {

    render() {
        return (
            //<nav aria-label="Page navigation">
            <ul className = 'pagination pagination-lg'>
                {this.props.pageNumberList.map(pageNumber => {
                    return (
                        <PageNumber pageNumber={pageNumber} changePageNumber ={this.props.changePageNumber} />
                    )
                })}
            </ul>
            //</nav>

        )
    }
}


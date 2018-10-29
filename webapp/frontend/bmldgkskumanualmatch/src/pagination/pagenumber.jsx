import React, { Component } from 'react';

export default class PageNumber extends Component {

    onClick(e){
        e.preventDefault();
        const {changePageNumber} = this.props;
        changePageNumber(this.props.pageNumber);
      }

  render() {
    return (
      <li className="page-item"><a className="page-link" onClick={this.onClick.bind(this)}>{this.props.pageNumber}</a></li>
    )
  }
}


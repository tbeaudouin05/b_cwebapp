import React, { Component } from 'react';

export default class SearchDgkSellerName extends Component {
  constructor(props) {
    super(props);
    this.state = { value: '' };

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(event) {
    this.setState({ value: event.target.value });
  }

  handleSubmit(event) {
    this.props.SearchDgkSellerName(this.state.value);
    event.preventDefault();
  }

  render() {
    return (

      <form className="form-inline" onSubmit={this.handleSubmit}>
           <p > </p>
           <p >  </p>

        <div className="form-group mx-sm-3 mb-2">
          <label htmlFor="search" className="sr-only">Search for supplier name..</label>
          <input type="text" className="form-control" id="search" placeholder="Search for supplier name.." value={this.state.value} onChange={this.handleChange} />
        </div>
        <p > </p>
        <button type="submit" className="btn btn-primary mb-2">Search</button>
        
      </form>


    );
  }
}
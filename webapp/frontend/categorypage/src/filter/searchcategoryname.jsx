import React, { Component } from 'react';

export default class SearchCategoryName extends Component {
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
    this.props.SearchCategoryName(this.state.value);
    event.preventDefault();
  }

  render() {
    return (

      <form className="form-inline" onSubmit={this.handleSubmit}>
        <div className="form-group mx-sm-3 mb-2">
          <label htmlFor="search" className="sr-only">Search for supplier name..</label>
          <input type="text" className="form-control" id="search" placeholder="Search for Category name.." value={this.state.value} onChange={this.handleChange} />
        </div>
        <button type="submit" className="btn btn-primary mb-2">Search</button>
      </form>


    );
  }
}
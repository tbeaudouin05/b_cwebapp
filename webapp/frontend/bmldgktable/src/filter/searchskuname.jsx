import React, { Component } from 'react';

export default class SearchSkuName extends Component {
    constructor(props) {
      super(props);
      this.state = {value: ''};
  
      this.handleChange = this.handleChange.bind(this);
      this.handleSubmit = this.handleSubmit.bind(this);
    }
  
    handleChange(event) {
      this.setState({value: event.target.value});
      console.log(event.target.value)
    }
  
    handleSubmit(event) {
      this.props.searchSkuName(this.state.value);
      console.log(this.state.value)
      event.preventDefault();
    }
  
    render() {
      return (

<form className="form-inline" onSubmit={this.handleSubmit}>
<div className=" form-group mx-sm-2 mb-1">
  <label htmlFor="search" className="sr-only">Search..</label>
  <input type="text" className="form-control" id="search" placeholder="Search.." value={this.state.value} onChange={this.handleChange}/>
</div>
<button type="submit" className="btn btn-primary mb-2">Search</button>
</form>

        
      );
    }
  }
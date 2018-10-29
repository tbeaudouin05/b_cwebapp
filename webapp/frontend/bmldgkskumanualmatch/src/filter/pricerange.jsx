import React, { Component } from 'react';

export default class FliterPriceRange extends Component {
  constructor(props) {
    super(props);
    this.state = {
      minValue: '-5',
      maxValue: '9999999999'
    };

    this.handleChangeMin = this.handleChangeMin.bind(this);
    this.handleChangeMax = this.handleChangeMax.bind(this);
    this.onClick = this.onClick.bind(this);
  }

  handleChangeMin(event) {
    this.setState({ minValue: event.target.value });
  }
  handleChangeMax(event) {
    this.setState({ maxValue: event.target.value });
  }

  onClick(event) {
    this.props.ApplyPriceFilter(this.state.minValue, this.state.maxValue);
    event.preventDefault();
  }

  render() {
    return (
      <div className='row'>
        <label className='center'>Price range:</label>

        <div className='row-gap-sm1'>
          <form className="form-inline" onSubmit={this.handleSubmit}>
            <div className="form-group mx-sm-3 mb-2">
              <label  className="sr-only">Minimum price</label>
              <input type="number" className="form-control " id="minprice" placeholder="Minimum price" value={this.state.value} onChange={this.handleChangeMin} />
            </div>
          </form>
        </div>

        <div className='row-gap-sm1'>
          <form className="form-inline" onSubmit={this.handleSubmit}>
            <div className="form-group mx-sm-3 mb-2">
              <label className="sr-only">Maximum price</label>
              <input type="number" className="form-control" id="maxprice" placeholder="Maximum price" value={this.state.value} onChange={this.handleChangeMax} />
            </div>
          </form>
        </div>
        <form className="row_button" onClick={this.onClick.bind(this)} target="_blank">
          <button type="submit" className="btn btn-primary mb-2 center">Filter</button>
        </form>

      </div>


    );
  }
}
import React, { Component } from 'react';

export default class OptionFilterCategory2 extends Component {
  render() {
    return (
      <option value={this.props.optionValue}>{this.props.optionText}</option>
    )
  }
}


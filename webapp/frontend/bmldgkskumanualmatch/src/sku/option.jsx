import React, { Component } from 'react';

export default class OptionFrequency extends Component {
    render() {
      return (
        <option value={this.props.OptionValue}>{this.props.OptionText}</option>
      )
    }
  }
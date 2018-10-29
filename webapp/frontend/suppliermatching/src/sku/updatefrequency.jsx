import React, { Component } from 'react';
import OptionFrequency from './option.jsx'


export default class UpdateFrequency extends Component {

    // cf. this: https://reactjs.org/docs/forms.html

    constructor(props) {
        super(props);
        this.state = { value: "1" };

        this.handleChange = this.handleChange.bind(this);
    }

    handleChange(event) {
        this.setState({ value: event.target.value });
        console.log(event.target.value)
        event.preventDefault();
    }

    onClick(e) {
        e.preventDefault();
        const { UpdateFrequency } = this.props;
        UpdateFrequency(this.state.value);
    }

    render() {
        return (
            <div className="row">
            <form> 
            <p > </p>
            <p ><strong> Change update frequency:</strong> </p>
                <select className="form-control" value={this.state.value} onChange={this.handleChange}>
                    {this.props.optionList.map(option => {
                        return (
                            <OptionFrequency OptionValue={option.OptionValue} OptionText={option.OptionText} />
                        )
                    })}
                </select>
                <form className="row_button" onClick={this.onClick.bind(this)} target="_blank">
                    <button type="submit" className="btn btn-primary mb-2">Submit</button>
                </form>
                <form className="row" UpdateStatus={this.props.UpdateStatus} >
                </form>
                <td key='4'> {this.props.UpdateStatus} </td>
            </form>
            </div>
        )
    }
}

